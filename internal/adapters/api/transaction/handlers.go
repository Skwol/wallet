package transaction

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/transaction"
	"github.com/skwol/wallet/pkg/logging"
)

const (
	transactionURL  = "/api/v1/transactions/{record_id}"
	transactionsURL = "/api/v1/transactions"
)

type handler struct {
	transactionService transaction.Service
}

func NewHandler(service transaction.Service) (adapters.Handler, error) {
	return &handler{transactionService: service}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(transactionURL, h.getTransaction)
	router.HandleFunc(transactionsURL, h.getAllTransactions)
}

func (h *handler) getTransaction(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	id, err := strconv.ParseInt(mux.Vars(r)["record_id"], 10, 64)
	if err != nil {
		logger.Errorf("error parsing id: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing id: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}
	walletDTO, err := h.transactionService.GetByID(r.Context(), id)
	if err != nil {
		logger.Errorf("error returned from service: %s", err.Error())
		http.Error(w, fmt.Sprintf("error returned from service: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(walletDTO)
	if err != nil {
		logger.Errorf("error marshaling transaction: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling transaction: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(response)
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
	type Data struct {
		Transactions []*transaction.TransactionDTO `json:"transactions"`
	}
	response, err := json.Marshal(Data{Transactions: transactions})
	if err != nil {
		logger.Errorf("error marshaling transactions: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling transactions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
