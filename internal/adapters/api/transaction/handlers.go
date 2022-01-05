package transaction

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	router.HandleFunc(transactionURL, h.getTransaction).Methods(http.MethodGet)
	router.HandleFunc(transactionsURL, h.getAllTransactions).Methods(http.MethodGet)

	router.HandleFunc(transactionsURL, h.getFilteredTransactions).Methods(http.MethodPost)
}

func (h *handler) getTransaction(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	id, err := strconv.ParseInt(mux.Vars(r)["record_id"], 10, 64)
	if err != nil {
		logger.Errorf("error parsing id: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing id: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}
	transactionDTO, err := h.transactionService.GetByID(r.Context(), id)
	if err != nil {
		logger.Errorf("error returned from service: %s", err.Error())
		http.Error(w, fmt.Sprintf("error returned from service: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if transactionDTO.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response, err := json.Marshal(transactionDTO)
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
	if len(transactions) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response, err := json.Marshal(transactions)
	if err != nil {
		logger.Errorf("error marshaling transactions: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling transactions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *handler) getFilteredTransactions(w http.ResponseWriter, r *http.Request) {
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("error reading request body: %s", err.Error())
		http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	var request transaction.FilterTransactionsDTO
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Errorf("error unmarshaling request: %s", err.Error())
		http.Error(w, fmt.Sprintf("error unmarshaling request: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	transactions, err := h.transactionService.GetFiltered(r.Context(), &request, limit, offset)
	if err != nil {
		logger.Errorf("error returned from service: %s", err.Error())
		http.Error(w, fmt.Sprintf("error returned from service: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if len(transactions) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response, err := json.Marshal(transactions)
	if err != nil {
		logger.Errorf("error marshaling transactions: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling transactions: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
