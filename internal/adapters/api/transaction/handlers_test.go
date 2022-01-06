package transaction

import (
	"bytes"
	"context"
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
	dbtransaction "github.com/skwol/wallet/internal/adapters/db/transaction"
	domaintransaction "github.com/skwol/wallet/internal/domain/transaction"
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

		storage, err := dbtransaction.NewStorage(dbClient)
		if err != nil {
			t.Fatalf("error creating transaction storage %s", err.Error())
		}
		service, err := domaintransaction.NewService(storage)
		if err != nil {
			t.Fatalf("error creating transaction service %s", err.Error())
		}
		handlerInterface, err := NewHandler(service)
		if err != nil {
			t.Fatalf("error creating transaction handler %s", err.Error())
		}
		transactionHandler := handlerInterface.(*handler)

		transactionHandler.Register(router)
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

func prepareAllTransactionsInDB(ctx context.Context, t *testing.T) []time.Time {
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

	tranOneDate := time.Date(2020, 10, 11, 10, 0, 0, 0, time.UTC)
	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO transaction (id, sender_id, receiver_id, amount, date, tran_type) values ($1, $2, $2, $3, $4, 'deposit');", 1, 1, 100, tranOneDate); err != nil {
		t.Fatalf("error creating transaction one: %s", err.Error())
	}

	tranTwoDate := time.Date(2021, 10, 11, 10, 0, 0, 0, time.UTC)
	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO transaction (id, sender_id, receiver_id, amount, date, tran_type) values ($1, $2, $2, $3, $4, 'deposit');", 2, 2, 200, tranTwoDate); err != nil {
		t.Fatalf("error creating transaction one: %s", err.Error())
	}

	tranThreeDate := time.Date(2021, 10, 12, 10, 0, 0, 0, time.UTC)
	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO transaction (id, sender_id, receiver_id, amount, date, tran_type) values ($1, $2, $3, $4, $5, 'transfer');", 3, 2, 1, 100, tranThreeDate); err != nil {
		t.Fatalf("error creating transaction one: %s", err.Error())
	}

	tranFourDate := time.Date(2021, 10, 13, 10, 0, 0, 0, time.UTC)
	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO transaction (id, sender_id, receiver_id, amount, date, tran_type) values ($1, $2, $2, $3, $4, 'withdraw');", 4, 2, 100, tranFourDate); err != nil {
		t.Fatalf("error creating transaction one: %s", err.Error())
	}
	return []time.Time{tranOneDate, tranTwoDate, tranThreeDate, tranFourDate}
}

func TestGetTransaction(t *testing.T) {
	setup(t)
	ctx := context.Background()

	if _, err := dbClient.Conn.QueryContext(ctx, "truncate wallet cascade;"); err != nil {
		t.Fatalf("error truncating wallet: %s", err.Error())
	}

	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO wallet (id, name, balance) VALUES ($1, $2, $3);", 1, "test_wallet_one", 100); err != nil {
		t.Fatalf("error creating wallet one: %s", err.Error())
	}

	tranDate := time.Date(2020, 10, 11, 10, 0, 0, 0, time.UTC)
	if _, err := dbClient.Conn.QueryContext(ctx, "INSERT INTO transaction (id, sender_id, receiver_id, amount, date, tran_type) values ($1, $2, $2, $3, $4, 'deposit');", 1, 1, 100, tranDate); err != nil {
		t.Fatalf("error creating transaction one: %s", err.Error())
	}

	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.DefaultClient.Do(newReq(t, http.MethodGet, ts.URL+"/api/v1/transactions/2", nil))
	if err != nil {
		t.Fatalf("error getting response: %s", err.Error())
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.StatusCode)
	}

	resp, err = http.DefaultClient.Do(newReq(t, http.MethodGet, ts.URL+"/api/v1/transactions/1", nil))
	if err != nil {
		t.Fatalf("error getting response: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading request: %s", err.Error())
	}

	var response Transaction
	if err := json.Unmarshal(result, &response); err != nil {
		t.Fatalf("error unmarshaling response: %s", err.Error())
	}
	expectedResponse := Transaction{
		ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDate, Type: string(domaintransaction.TranTypeDeposit),
	}
	if !reflect.DeepEqual(expectedResponse, response) {
		t.Fatalf("wrong transaction returned, expected: %+v, got: %+v", expectedResponse, response)
	}
}

func TestGetTransactions(t *testing.T) {
	setup(t)
	ctx := context.Background()
	tranDates := prepareAllTransactionsInDB(ctx, t)

	ts := httptest.NewServer(router)
	defer ts.Close()

	type args struct {
		limit, offset int
	}
	tests := []struct {
		name           string
		args           args
		want           []Transaction
		wantStatusCode int
	}{
		{
			name:           "no transactions",
			args:           args{limit: 100, offset: 4},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name: "all transactions",
			args: args{limit: 10, offset: 0},
			want: []Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDates[0], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
				{ID: 4, SenderID: 2, ReceiverID: 2, Amount: 100, Timestamp: tranDates[3], Type: string(domaintransaction.TranTypeWithdraw)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name: "all transactions with offset and limit",
			args: args{limit: 2, offset: 1},
			want: []Transaction{
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
			},
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(newReq(t, http.MethodGet, fmt.Sprintf("%s/api/v1/transactions?limit=%d&offset=%d", ts.URL, tt.args.limit, tt.args.offset), nil))
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
			var response []Transaction
			if err := json.Unmarshal(result, &response); err != nil {
				t.Fatalf("test %s: error unmarshaling response: %s", tt.name, err.Error())
			}
			if !reflect.DeepEqual(tt.want, response) {
				t.Fatalf("test %s: wrong transactions returned, expected: %+v, got: %+v", tt.name, tt.want, response)
			}
		})
	}

	type argsFilters struct {
		limit, offset int
	}
	testsFilters := []struct {
		name           string
		args           argsFilters
		request        Filter
		want           []Transaction
		wantStatusCode int
	}{
		{
			name:           "filtered no transactions",
			args:           argsFilters{limit: 100, offset: 4},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:    "filtered all transactions",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{},
			want: []Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDates[0], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
				{ID: 4, SenderID: 2, ReceiverID: 2, Amount: 100, Timestamp: tranDates[3], Type: string(domaintransaction.TranTypeWithdraw)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered all transactions with offset and limit",
			args:    argsFilters{limit: 2, offset: 1},
			request: Filter{},
			want: []Transaction{
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered transactions by single sender id",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{SenderIDs: []int64{1}},
			want: []Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDates[0], Type: string(domaintransaction.TranTypeDeposit)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered transactions by several sender id",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{SenderIDs: []int64{1, 2, 3}},
			want: []Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDates[0], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
				{ID: 4, SenderID: 2, ReceiverID: 2, Amount: 100, Timestamp: tranDates[3], Type: string(domaintransaction.TranTypeWithdraw)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered transactions by amount more then 100",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{Amount: FloatRangeFilter{From: 100}},
			want: []Transaction{
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered transactions by amount less then 200",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{Amount: FloatRangeFilter{To: 200}},
			want: []Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDates[0], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
				{ID: 4, SenderID: 2, ReceiverID: 2, Amount: 100, Timestamp: tranDates[3], Type: string(domaintransaction.TranTypeWithdraw)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered transactions by amount more less then or equal to 200 and more then or equal to 100",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{Amount: FloatRangeFilter{From: 100, To: 200}},
			want: []Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDates[0], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
				{ID: 4, SenderID: 2, ReceiverID: 2, Amount: 100, Timestamp: tranDates[3], Type: string(domaintransaction.TranTypeWithdraw)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "filtered transactions by amount more less then 200 more then 100",
			args:           argsFilters{limit: 10, offset: 0},
			request:        Filter{Amount: FloatRangeFilter{From: 101, To: 199}},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:           "filtered transactions by amount more less then 100",
			args:           argsFilters{limit: 10, offset: 0},
			request:        Filter{Amount: FloatRangeFilter{To: 100}},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:    "filtered transactions by timestamp after first",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{Timestamp: DateRangeFilter{From: tranDates[0]}},
			want: []Transaction{
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
				{ID: 4, SenderID: 2, ReceiverID: 2, Amount: 100, Timestamp: tranDates[3], Type: string(domaintransaction.TranTypeWithdraw)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered transactions by timestamp before last",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{Timestamp: DateRangeFilter{To: tranDates[3]}},
			want: []Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDates[0], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered transactions by timestamp before last after first (inclusive)",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{Timestamp: DateRangeFilter{From: tranDates[0], To: tranDates[3]}},
			want: []Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDates[0], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
				{ID: 4, SenderID: 2, ReceiverID: 2, Amount: 100, Timestamp: tranDates[3], Type: string(domaintransaction.TranTypeWithdraw)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered transactions by timestamp before last after first (inclusive)",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{Timestamp: DateRangeFilter{From: tranDates[0].Add(time.Minute), To: tranDates[3].Add(-time.Minute)}},
			want: []Transaction{
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "filtered transactions by timestamp before first",
			args:           argsFilters{limit: 10, offset: 0},
			request:        Filter{Timestamp: DateRangeFilter{To: tranDates[0]}},
			wantStatusCode: http.StatusNotFound,
		},
		{
			name:    "filtered transactions by types deposit and transfer",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{Types: []string{string(domaintransaction.TranTypeDeposit), string(domaintransaction.TranTypeTransfer)}},
			want: []Transaction{
				{ID: 1, SenderID: 1, ReceiverID: 1, Amount: 100, Timestamp: tranDates[0], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 2, SenderID: 2, ReceiverID: 2, Amount: 200, Timestamp: tranDates[1], Type: string(domaintransaction.TranTypeDeposit)},
				{ID: 3, SenderID: 2, ReceiverID: 1, Amount: 100, Timestamp: tranDates[2], Type: string(domaintransaction.TranTypeTransfer)},
			},
			wantStatusCode: http.StatusOK,
		},
		{
			name:    "filtered transactions by types withdraw",
			args:    argsFilters{limit: 10, offset: 0},
			request: Filter{Types: []string{string(domaintransaction.TranTypeWithdraw)}},
			want: []Transaction{
				{ID: 4, SenderID: 2, ReceiverID: 2, Amount: 100, Timestamp: tranDates[3], Type: string(domaintransaction.TranTypeWithdraw)},
			},
			wantStatusCode: http.StatusOK,
		},
	}
	for _, tt := range testsFilters {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := http.DefaultClient.Do(newReq(t, http.MethodPost, fmt.Sprintf("%s/api/v1/transactions?limit=%d&offset=%d", ts.URL, tt.args.limit, tt.args.offset), tt.request))
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
			var response []Transaction
			if err := json.Unmarshal(result, &response); err != nil {
				t.Fatalf("test %s: error unmarshaling response: %s", tt.name, err.Error())
			}
			if !reflect.DeepEqual(tt.want, response) {
				t.Fatalf("test %s: wrong transactions returned, expected: %+v, got: %+v", tt.name, tt.want, response)
			}
		})
	}
}
