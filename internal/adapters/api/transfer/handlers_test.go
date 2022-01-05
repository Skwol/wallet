package transfer

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	dbtransfer "github.com/skwol/wallet/internal/adapters/db/transfer"
	domaintransfer "github.com/skwol/wallet/internal/domain/transfer"
	"github.com/skwol/wallet/pkg/client/pgdb"
	"github.com/skwol/wallet/pkg/logging"
	"github.com/skwol/wallet/pkg/testdb"
)

var (
	once     sync.Once
	router   *mux.Router
	dbClient *pgdb.PGDB
)

func setup(t *testing.T) {
	once.Do(func() {
		logging.Init()

		router = mux.NewRouter()

		var err error
		dbClient, err = testdb.DBClient()
		if err != nil {
			t.Fatalf("error creating db client: %s", err.Error())
		}
		if dbClient == nil {
			t.Fatal("missing db client")
		}

		storage, err := dbtransfer.NewStorage(dbClient)
		if err != nil {
			t.Fatalf("error creating transfer storage %s", err.Error())
		}
		service, err := domaintransfer.NewService(storage)
		if err != nil {
			t.Fatalf("error creating transfer service %s", err.Error())
		}
		handlerInterface, err := NewHandler(service)
		if err != nil {
			t.Fatalf("error creating transfer handler %s", err.Error())
		}
		transferHandler := handlerInterface.(*handler)

		transferHandler.Register(router)
	})
}

func newReq(t *testing.T, method, url string, body interface{}) *http.Request {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		t.Fatalf("test %s: error encoding request: %s", body, err.Error())
	}
	r, err := http.NewRequest(method, url, &buf)
	if err != nil {
		t.Fatal(err)
	}
	return r
}

func TestCreateTransfer(t *testing.T) {
	setup(t)
	ctx := context.Background()

	if _, err := dbClient.Conn.QueryContext(ctx, "truncate wallet cascade;"); err != nil {
		t.Fatalf("error truncating wallet: %s", err.Error())
	}
	if _, err := dbClient.Conn.QueryContext(ctx, "truncate transaction;"); err != nil {
		t.Fatalf("error truncating transaction: %s", err.Error())
	}

	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO wallet (id, name, balance) VALUES ($1, $2, $3);", 1, "test_wallet_one", 100); err != nil {
		t.Fatalf("error creating wallet one: %s", err.Error())
	}

	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO wallet (id, name, balance) VALUES ($1, $2, $3);", 2, "test_wallet_two", 200); err != nil {
		t.Fatalf("error creating wallet two: %s", err.Error())
	}

	ts := httptest.NewServer(router)
	defer ts.Close()

	type args struct {
		request domaintransfer.CreateTransferDTO
	}
	tests := []struct {
		name               string
		args               args
		want               domaintransfer.TransferDTO
		wantTransaction    domaintransfer.TransferDTO
		wantWalletBalances map[int]float64
		wantStatusCode     int
	}{
		{
			name:           "transfer when sender == receiver",
			args:           args{domaintransfer.CreateTransferDTO{Amount: 100, Sender: domaintransfer.WalletDTO{ID: 1}, Receiver: domaintransfer.WalletDTO{ID: 1}}},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:           "transfer when sender does not have enough",
			args:           args{domaintransfer.CreateTransferDTO{Amount: 1000, Sender: domaintransfer.WalletDTO{ID: 1}, Receiver: domaintransfer.WalletDTO{ID: 2}}},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "transfer OK",
			args:               args{domaintransfer.CreateTransferDTO{Amount: 100, Sender: domaintransfer.WalletDTO{ID: 1}, Receiver: domaintransfer.WalletDTO{ID: 2}}},
			want:               domaintransfer.TransferDTO{CreateTransferDTO: domaintransfer.CreateTransferDTO{Amount: 100, Sender: domaintransfer.WalletDTO{ID: 1, Balance: 0}, Receiver: domaintransfer.WalletDTO{ID: 2, Balance: 300}}},
			wantTransaction:    domaintransfer.TransferDTO{CreateTransferDTO: domaintransfer.CreateTransferDTO{Amount: 100, Sender: domaintransfer.WalletDTO{ID: 1}, Receiver: domaintransfer.WalletDTO{ID: 2}}},
			wantWalletBalances: map[int]float64{1: 0, 2: 300},
			wantStatusCode:     http.StatusCreated,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(newReq(t, http.MethodPost, ts.URL+"/api/v1/transfers", tt.args.request))
			if err != nil {
				t.Fatalf("test %s: error getting response: %s", tt.name, err.Error())
			}
			if resp == nil {
				t.Fatalf("test %s: missing response", tt.name)
			}
			if resp.StatusCode != tt.wantStatusCode {
				t.Fatalf("test %s: expected status %d, got %d", tt.name, tt.wantStatusCode, resp.StatusCode)
			}
			result, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("test %s: error reading request: %s", tt.name, err.Error())
			}

			if tt.want.CreateTransferDTO.Amount == 0 {
				var got domaintransfer.TransferDTO
				if err := json.Unmarshal(result, &got); err == nil {
					t.Fatalf("test %s: should not receive correct response from server", tt.name)
				}
				return
			}

			var got domaintransfer.TransferDTO
			if err := json.Unmarshal(result, &got); err != nil {
				t.Fatalf("test %s: error unmarshaling response: %s", tt.name, err.Error())
			}
			tt.want.ID = got.ID
			tt.want.Timestamp = got.Timestamp
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("test %s: wrong transfer returned, expected: %+v, got: %+v", tt.name, tt.want, got)
			}

			// test wallets in db
			var balance float64
			for walletID, expectedBalance := range tt.wantWalletBalances {
				row := dbClient.Conn.QueryRowContext(ctx, `SELECT balance FROM wallet WHERE id = $1;`, walletID)
				switch err := row.Scan(&balance); err {
				case sql.ErrNoRows:
					t.Fatalf("test %s: missing wallet %d in db", tt.name, walletID)
				default:
					if balance != expectedBalance {
						t.Fatalf("test %s: wrong balance of wallet %d in db, expected %f, got %f", tt.name, walletID, expectedBalance, balance)
					}
				}
			}

			// test transactions in db
			var transactionInDB domaintransfer.TransferDTO

			row := dbClient.Conn.QueryRowContext(ctx, "SELECT sender_id, receiver_id, amount FROM transaction WHERE id = $1 and tran_type = 'transfer'", got.ID)
			switch err := row.Scan(&transactionInDB.CreateTransferDTO.Sender.ID, &transactionInDB.CreateTransferDTO.Receiver.ID, &transactionInDB.CreateTransferDTO.Amount); err {
			case sql.ErrNoRows:
				t.Fatalf("test %s: missing transaction %d in db", tt.name, got.ID)
			default:
				if !reflect.DeepEqual(tt.wantTransaction, transactionInDB) {
					t.Fatalf("test %s: wrong transaction in db, expected: %+v, got: %+v", tt.name, tt.wantTransaction, transactionInDB)
				}
			}
		})
	}
}
