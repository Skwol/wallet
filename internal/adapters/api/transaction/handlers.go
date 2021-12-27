package transaction

import (
	"net/http"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/transaction"
)

const (
	getTransactionURL  = "/api/v1/transactions/{record_id}"
	getTransactionsURL = "/api/v1/transactions"
)

type handler struct {
	transactionService transaction.Service
}

func NewHandler(service transaction.Service) adapters.Handler {
	return &handler{}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(getTransactionURL, h.getTransaction)
	router.HandleFunc(getTransactionsURL, h.getTransactions)
}

func (h *handler) getTransaction(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) getTransactions(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("transactions"))
	w.WriteHeader(http.StatusOK)
}
