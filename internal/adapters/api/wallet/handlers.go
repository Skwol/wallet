package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/skwol/wallet/pkg/logging"

	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/wallet"
)

const (
	walletURL                 = "/api/v1/wallets/{record_id}"
	walletWithTransactionsURL = "/api/v1/wallets-with-transactions/{record_id}"
	walletsURL                = "/api/v1/wallets"
)

type handler struct {
	walletService wallet.Service
	logger        logging.Logger
}

func NewHandler(service wallet.Service, logger logging.Logger) (adapters.Handler, error) {
	return &handler{walletService: service, logger: logger}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(walletsURL, h.getAllWallets).Methods(http.MethodGet)
	router.HandleFunc(walletURL, h.getWallet).Methods(http.MethodGet)
	router.HandleFunc(walletWithTransactionsURL, h.getWalletWithTransactions).Methods(http.MethodGet)

	router.HandleFunc(walletURL, h.updateWallet).Methods(http.MethodPatch)

	router.HandleFunc(walletsURL, h.createWallet).Methods(http.MethodPost)
}

func (h *handler) getWallet(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	id, err := strconv.ParseInt(mux.Vars(r)["record_id"], 10, 64)
	if err != nil {
		logger.Errorf("error parsing id: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing id: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}
	walletDTO, err := h.walletService.GetByID(r.Context(), id)
	if err != nil {
		logger.Errorf("error returned from service: %s", err.Error())
		http.Error(w, fmt.Sprintf("error returned from service: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if walletDTO.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	response, err := json.Marshal(newWallet(walletDTO))
	if err != nil {
		logger.Errorf("error marshaling wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling wallet: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(response); err != nil {
		logger.Errorf("error writing response: %s", err.Error())
		http.Error(w, fmt.Sprintf("error writing response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *handler) getWalletWithTransactions(w http.ResponseWriter, r *http.Request) {
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

	id, err := strconv.ParseInt(mux.Vars(r)["record_id"], 10, 64)
	if err != nil {
		logger.Errorf("error parsing id: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing id: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}
	walletDTO, err := h.walletService.GetByIDWithTransactions(r.Context(), id, limit, offset)
	if err != nil {
		logger.Errorf("error returned from service: %s", err.Error())
		http.Error(w, fmt.Sprintf("error returned from service: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(newWallet(walletDTO))
	if err != nil {
		logger.Errorf("error marshaling wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling wallet: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(response); err != nil {
		logger.Errorf("error writing response: %s", err.Error())
		http.Error(w, fmt.Sprintf("error writing response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *handler) updateWallet(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	id, err := strconv.ParseInt(mux.Vars(r)["record_id"], 10, 64)
	if err != nil {
		logger.Errorf("error parsing id: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing id: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("error reading request body: %s", err.Error())
		http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}
	var request Wallet
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Errorf("error unmarshaling request: %s", err.Error())
		http.Error(w, fmt.Sprintf("error unmarshaling request: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	updateRequest := request.toUpdateRequest()
	walletDTO, err := h.walletService.Update(r.Context(), id, &updateRequest)
	if err != nil {
		logger.Errorf("error updating wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error updating wallet: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}
	response, err := json.Marshal(newWallet(walletDTO))
	if err != nil {
		logger.Errorf("error marshaling wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling wallet: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(response); err != nil {
		logger.Errorf("error writing response: %s", err.Error())
		http.Error(w, fmt.Sprintf("error writing response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *handler) createWallet(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("error reading request body: %s", err.Error())
		http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	var request Wallet
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Errorf("error unmarshaling request: %s", err.Error())
		http.Error(w, fmt.Sprintf("error unmarshaling request: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	createRequest := request.toCreateRequest()
	walletDTO, err := h.walletService.Create(r.Context(), &createRequest)
	if err != nil {
		logger.Errorf("error creating wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error creating wallet: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	response, err := json.Marshal(newWallet(walletDTO))
	if err != nil {
		logger.Errorf("error marshaling wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling wallet: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(response); err != nil {
		logger.Errorf("error writing response: %s", err.Error())
		http.Error(w, fmt.Sprintf("error writing response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}

func (h *handler) getAllWallets(w http.ResponseWriter, r *http.Request) {
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

	walletDTOs, err := h.walletService.GetAll(r.Context(), limit, offset)
	if err != nil {
		logger.Errorf("error returned from service: %s", err.Error())
		http.Error(w, fmt.Sprintf("error returned from service: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if len(walletDTOs) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	wallets := make([]Wallet, 0, len(walletDTOs))
	for _, dto := range walletDTOs {
		wallets = append(wallets, newWallet(dto))
	}
	response, err := json.Marshal(wallets)
	if err != nil {
		logger.Errorf("error marshaling wallets: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling wallets: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(response); err != nil {
		logger.Errorf("error writing response: %s", err.Error())
		http.Error(w, fmt.Sprintf("error writing response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
