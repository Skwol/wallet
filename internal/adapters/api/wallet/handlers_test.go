package wallet

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	dbwallet "github.com/skwol/wallet/internal/adapters/db/wallet"
	"github.com/skwol/wallet/internal/domain/wallet"
	domainwallet "github.com/skwol/wallet/internal/domain/wallet"
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

		storage, err := dbwallet.NewStorage(dbClient)
		if err != nil {
			t.Fatalf("error creating wallet storage %s", err.Error())
		}
		service, err := domainwallet.NewService(storage)
		if err != nil {
			t.Fatalf("error creating wallet service %s", err.Error())
		}
		handlerInterface, err := NewHandler(service)
		if err != nil {
			t.Fatalf("error creating wallet handler %s", err.Error())
		}
		walletHandler := handlerInterface.(*handler)

		walletHandler.Register(router)
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

func TestGetWallets(t *testing.T) {
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

	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO wallet (id, name, balance) VALUES ($1, $2, $3);", 3, "test_wallet_three", 300); err != nil {
		t.Fatalf("error creating wallet three: %s", err.Error())
	}

	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO wallet (id, name, balance) VALUES ($1, $2, $3);", 4, "test_wallet_four", 400); err != nil {
		t.Fatalf("error creating wallet four: %s", err.Error())
	}

	tranOneDate := time.Date(2020, 10, 11, 10, 0, 0, 0, time.UTC)
	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO transaction (id, sender_id, receiver_id, amount, date, tran_type) values ($1, $2, $2, $3, $4, 'deposit');", 1, 1, 100, tranOneDate); err != nil {
		t.Fatalf("error creating transaction one: %s", err.Error())
	}

	tranTwoDate := time.Date(2021, 10, 11, 10, 0, 0, 0, time.UTC)
	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO transaction (id, sender_id, receiver_id, amount, date, tran_type) values ($1, $2, $3, $4, $5, 'transfer');", 2, 2, 1, 100, tranTwoDate); err != nil {
		t.Fatalf("error creating transaction one: %s", err.Error())
	}

	tranThreeDate := time.Date(2021, 10, 12, 10, 0, 0, 0, time.UTC)
	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO transaction (id, sender_id, receiver_id, amount, date, tran_type) values ($1, $2, $2, $3, $4, 'withdraw');", 3, 1, 100, tranThreeDate); err != nil {
		t.Fatalf("error creating transaction one: %s", err.Error())
	}

	ts := httptest.NewServer(router)
	defer ts.Close()

	type args struct {
		limit, offset int
		endpoint      string
	}
	tests := []struct {
		name           string
		args           args
		want           []domainwallet.WalletDTO
		wantStatusCode int
		singleValue    bool
	}{
		{
			name:           "no wallet",
			args:           args{endpoint: fmt.Sprintf("%s/api/v1/wallets/5", ts.URL)},
			wantStatusCode: http.StatusNotFound,
			singleValue:    true,
		},
		{
			name: "wallet without transactions",
			args: args{endpoint: fmt.Sprintf("%s/api/v1/wallets/1", ts.URL)},
			want: []domainwallet.WalletDTO{
				{ID: 1, Name: "test_wallet_one", Balance: 100},
			},
			wantStatusCode: http.StatusOK,
			singleValue:    true,
		},
		{
			name: "all wallets",
			args: args{endpoint: fmt.Sprintf("%s/api/v1/wallets?limit=10&offset=0", ts.URL)},
			want: []domainwallet.WalletDTO{
				{ID: 1, Name: "test_wallet_one", Balance: 100},
				{ID: 2, Name: "test_wallet_two", Balance: 200},
				{ID: 3, Name: "test_wallet_three", Balance: 300},
				{ID: 4, Name: "test_wallet_four", Balance: 400},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "all wallets limited",
			args: args{endpoint: fmt.Sprintf("%s/api/v1/wallets?limit=1&offset=1", ts.URL)},
			want: []domainwallet.WalletDTO{
				{ID: 2, Name: "test_wallet_two", Balance: 200},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "all wallets no results",
			args:           args{endpoint: fmt.Sprintf("%s/api/v1/wallets?limit=10&offset=4", ts.URL)},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "wallet with all transactions",
			args: args{endpoint: fmt.Sprintf("%s/api/v1/wallets-with-transactions/1?limit=10&offset=0", ts.URL)},
			want: []domainwallet.WalletDTO{
				{ID: 1, Name: "test_wallet_one", Balance: 100, Transactions: []domainwallet.TransactionDTO{
					{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranOneDate, Type: domainwallet.TranTypeDeposit},
					{ID: 2, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranTwoDate, Type: domainwallet.TranTypeTransfer},
					{ID: 3, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranThreeDate, Type: domainwallet.TranTypeWithdraw},
				}},
			},
			wantStatusCode: http.StatusOK,
			singleValue:    true,
		},
		{
			name: "wallet with all transactions limited with offset",
			args: args{endpoint: fmt.Sprintf("%s/api/v1/wallets-with-transactions/1?limit=1&offset=1", ts.URL)},
			want: []domainwallet.WalletDTO{
				{ID: 1, Name: "test_wallet_one", Balance: 100, Transactions: []domainwallet.TransactionDTO{
					{ID: 2, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranTwoDate, Type: domainwallet.TranTypeTransfer},
				}},
			},
			wantStatusCode: http.StatusOK,
			singleValue:    true,
		},
		{
			name: "wallet with all transactions limited with offset out of values",
			args: args{endpoint: fmt.Sprintf("%s/api/v1/wallets-with-transactions/1?limit=1&offset=10", ts.URL)},
			want: []domainwallet.WalletDTO{
				{ID: 1, Name: "test_wallet_one", Balance: 100},
			},
			wantStatusCode: http.StatusOK,
			singleValue:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(newReq(t, http.MethodGet, tt.args.endpoint, nil))
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

			if len(tt.want) == 0 {
				if len(result) > 0 {
					t.Fatalf("test %s: empty response expected, got '%s'", tt.name, string(result))
				}
				return
			}

			var got []domainwallet.WalletDTO
			if tt.singleValue {
				var response domainwallet.WalletDTO
				if err := json.Unmarshal(result, &response); err != nil {
					t.Fatalf("test %s: error unmarshaling response: %s", tt.name, err.Error())
				}
				got = []domainwallet.WalletDTO{response}
			} else {
				if err := json.Unmarshal(result, &got); err != nil {
					t.Fatalf("test %s: error unmarshaling response: %s", tt.name, err.Error())
				}
			}
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("test %s: wrong wallets returned, expected: %+v, got: %+v", tt.name, tt.want, got)
			}
		})
	}
}

func TestUpdateWallet(t *testing.T) {
	setup(t)
	ctx := context.Background()

	if _, err := dbClient.Conn.QueryContext(ctx, "truncate wallet cascade;"); err != nil {
		t.Fatalf("error truncating wallet: %s", err.Error())
	}
	if _, err := dbClient.Conn.QueryContext(ctx, "truncate transaction;"); err != nil {
		t.Fatalf("error truncating transaction: %s", err.Error())
	}

	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO wallet (id, name, balance) VALUES ($1, $2, $3);", 1, "wallet_one", 100); err != nil {
		t.Fatalf("error creating wallet one: %s", err.Error())
	}

	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO wallet (id, name, balance) VALUES ($1, $2, $3);", 2, "wallet_two", 200); err != nil {
		t.Fatalf("error creating wallet two: %s", err.Error())
	}

	ts := httptest.NewServer(router)
	defer ts.Close()

	type args struct {
		request wallet.UpdateWalletDTO
		enpoint string
	}
	tests := []struct {
		name             string
		args             args
		want             domainwallet.WalletDTO
		wantTransactions []domainwallet.TransactionDTO
		wantStatusCode   int
	}{
		{
			name:             "update wallet, withdraw 100 to become 0",
			args:             args{request: domainwallet.UpdateWalletDTO{Name: "wallet_one", Balance: 0}, enpoint: "/api/v1/wallets/1"},
			want:             domainwallet.WalletDTO{ID: 1, Name: "wallet_one", Balance: 0},
			wantTransactions: []domainwallet.TransactionDTO{{SenderID: 1, ReceiverID: 1, Amount: 100, Type: wallet.TranTypeWithdraw}},
			wantStatusCode:   http.StatusOK,
		},
		{
			name:           "update wallet, withdraw 300 to become negative",
			args:           args{request: domainwallet.UpdateWalletDTO{Name: "wallet_two", Balance: -100}, enpoint: "/api/v1/wallets/2"},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:             "update wallet, deposit 100 to become 300",
			args:             args{request: domainwallet.UpdateWalletDTO{Name: "wallet_two", Balance: 300}, enpoint: "/api/v1/wallets/2"},
			want:             domainwallet.WalletDTO{ID: 2, Name: "wallet_two", Balance: 300},
			wantTransactions: []domainwallet.TransactionDTO{{SenderID: 2, ReceiverID: 2, Amount: 100, Type: domainwallet.TranTypeDeposit}},
			wantStatusCode:   http.StatusOK,
		},
		{
			name:           "update non existing wallet",
			args:           args{request: domainwallet.UpdateWalletDTO{Name: "wallet_three", Balance: 300}, enpoint: "/api/v1/wallets/3"},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(newReq(t, http.MethodPatch, ts.URL+tt.args.enpoint, tt.args.request))
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

			if tt.want.Name == "" {
				var got domainwallet.WalletDTO
				if err := json.Unmarshal(result, &got); err == nil {
					t.Fatalf("test %s: should not receive correct response from server", tt.name)
				}
				return
			}

			var got domainwallet.WalletDTO
			if err := json.Unmarshal(result, &got); err != nil {
				t.Fatalf("test %s: error unmarshaling response: %s", tt.name, err.Error())
			}
			tt.want.ID = got.ID
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("test %s: wrong wallet returned, expected: %+v, got: %+v", tt.name, tt.want, got)
			}

			// test wallet in db
			row := dbClient.Conn.QueryRowContext(ctx, `SELECT id, name, balance FROM wallet WHERE id = $1;`, got.ID)
			var gotInDB domainwallet.WalletDTO
			switch err := row.Scan(&gotInDB.ID, &gotInDB.Name, &gotInDB.Balance); err {
			case sql.ErrNoRows:
				t.Fatalf("test %s: wallet was not created", tt.name)
			default:
				if !reflect.DeepEqual(tt.want, gotInDB) {
					t.Fatalf("test %s: wrong wallet in db, expected: %+v, got: %+v", tt.name, tt.want, gotInDB)
				}
			}

			// test transactions in db
			var transactionsInDB []domainwallet.TransactionDTO

			rows, err := dbClient.Conn.QueryContext(ctx, "SELECT sender_id, receiver_id, amount, tran_type FROM transaction WHERE sender_id = $1 OR receiver_id = $1 ORDER BY ID ASC", got.ID)
			if err != nil {
				t.Fatalf("test %s: error getting transactions from db: %s", tt.name, err.Error())
			}
			var tran domainwallet.TransactionDTO
			for rows.Next() {
				if len(tt.wantTransactions) == 0 {
					t.Fatalf("test %s: expeted 0 transactions, got some in db", tt.name)
				}
				if err := rows.Scan(&tran.SenderID, &tran.ReceiverID, &tran.Amount, &tran.Type); err != nil {
					t.Fatalf("test %s: error scanning transaction from db: %s", tt.name, err.Error())
				}
				transactionsInDB = append(transactionsInDB, tran)
			}
			if !reflect.DeepEqual(tt.wantTransactions, transactionsInDB) {
				t.Fatalf("test %s: wrong transactions in db, expected: %+v, got: %+v", tt.name, tt.wantTransactions, transactionsInDB)
			}
		})
	}
}

func TestCreateWallet(t *testing.T) {
	setup(t)
	ctx := context.Background()

	if _, err := dbClient.Conn.QueryContext(ctx, "truncate wallet cascade;"); err != nil {
		t.Fatalf("error truncating wallet: %s", err.Error())
	}
	if _, err := dbClient.Conn.QueryContext(ctx, "truncate transaction;"); err != nil {
		t.Fatalf("error truncating transaction: %s", err.Error())
	}

	ts := httptest.NewServer(router)
	defer ts.Close()

	type args struct {
		request wallet.CreateWalletDTO
	}
	tests := []struct {
		name             string
		args             args
		want             domainwallet.WalletDTO
		wantTransactions []domainwallet.TransactionDTO
		wantStatusCode   int
	}{
		{
			name:           "create wallet with 0 balance",
			args:           args{domainwallet.CreateWalletDTO{Name: "wallet_one"}},
			want:           domainwallet.WalletDTO{Name: "wallet_one"},
			wantStatusCode: http.StatusCreated,
		},
		{
			name:             "create wallet with 100 balance",
			args:             args{domainwallet.CreateWalletDTO{Name: "wallet_two", Balance: 100}},
			want:             domainwallet.WalletDTO{Name: "wallet_two", Balance: 100},
			wantTransactions: []domainwallet.TransactionDTO{{Amount: 100, Type: domainwallet.TranTypeDeposit}},
			wantStatusCode:   http.StatusCreated,
		},
		{
			name:           "create wallet negative balance",
			args:           args{domainwallet.CreateWalletDTO{Name: "wallet_three", Balance: -10}},
			wantStatusCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(newReq(t, http.MethodPost, ts.URL+"/api/v1/wallets", tt.args.request))
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

			if tt.want.Name == "" {
				var got domainwallet.WalletDTO
				if err := json.Unmarshal(result, &got); err == nil {
					t.Fatalf("test %s: should not receive correct response from server", tt.name)
				}
				return
			}

			var got domainwallet.WalletDTO
			if err := json.Unmarshal(result, &got); err != nil {
				t.Fatalf("test %s: error unmarshaling response: %s", tt.name, err.Error())
			}
			tt.want.ID = got.ID
			if !reflect.DeepEqual(tt.want, got) {
				t.Fatalf("test %s: wrong wallet returned, expected: %+v, got: %+v", tt.name, tt.want, got)
			}

			// test wallet in db
			row := dbClient.Conn.QueryRowContext(ctx, `SELECT id, name, balance FROM wallet WHERE id = $1;`, got.ID)
			var gotInDB domainwallet.WalletDTO
			switch err := row.Scan(&gotInDB.ID, &gotInDB.Name, &gotInDB.Balance); err {
			case sql.ErrNoRows:
				t.Fatalf("test %s: wallet was not created", tt.name)
			default:
				if !reflect.DeepEqual(tt.want, gotInDB) {
					t.Fatalf("test %s: wrong wallet in db, expected: %+v, got: %+v", tt.name, tt.want, gotInDB)
				}
			}

			// test transactions in db
			var transactionsInDB []domainwallet.TransactionDTO

			rows, err := dbClient.Conn.QueryContext(ctx, "SELECT amount, tran_type FROM transaction WHERE sender_id = $1 OR receiver_id = $1 ORDER BY ID ASC", got.ID)
			if err != nil {
				t.Fatalf("test %s: error getting transactions from db: %s", tt.name, err.Error())
			}
			var tran domainwallet.TransactionDTO
			for rows.Next() {
				if len(tt.wantTransactions) == 0 {
					t.Fatalf("test %s: expeted 0 transactions, got some in db", tt.name)
				}
				if err := rows.Scan(&tran.Amount, &tran.Type); err != nil {
					t.Fatalf("test %s: error scanning transaction from db: %s", tt.name, err.Error())
				}
				transactionsInDB = append(transactionsInDB, tran)
			}
			if !reflect.DeepEqual(tt.wantTransactions, transactionsInDB) {
				t.Fatalf("test %s: wrong transactions in db, expected: %+v, got: %+v", tt.name, tt.wantTransactions, transactionsInDB)
			}
		})
	}
}
