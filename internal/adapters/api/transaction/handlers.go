package transaction

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/transaction"
	"github.com/skwol/wallet/pkg/logging"
)

const (
	getTransactionURL  = "/api/v1/transactions/{record_id}"
	getTransactionsURL = "/api/v1/transactions"
)

type handler struct {
	transactionService transaction.Service
}

func NewHandler(service transaction.Service) (adapters.Handler, error) {
	return &handler{transactionService: service}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(getTransactionURL, h.getTransaction)
	router.HandleFunc(getTransactionsURL, h.getAllTransactions)
}

func (h *handler) getTransaction(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) getAllTransactions(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	limit, err := strconv.Atoi(r.FormValue("limit"))
	if err != nil {
		logger.Errorf("error parsing limit query param: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing limit query param: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	offset, err := strconv.Atoi(r.FormValue("offset"))
	if err != nil {
		logger.Errorf("error parsing offset query param: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing offset query param: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	transactions, err := h.transactionService.GetAll(r.Context(), limit, offset)
	if err != nil {
		logger.Errorf("error returned from service: %s", err.Error())
		http.Error(w, fmt.Sprintf("error returned from service: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	tmpl, err := template.ParseFiles("templates/transactions/transactions.html")
	if err != nil {
		logger.Errorf("error parsing transactions template: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing transactions template: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if tmpl == nil {
		logger.Error("missing transactions template")
		http.Error(w, "missing transactions template", http.StatusInternalServerError)
		return
	}
	type Data struct {
		Transactions []*transaction.TransactionDTO
	}
	tmpl.Execute(w, Data{Transactions: transactions})
}
