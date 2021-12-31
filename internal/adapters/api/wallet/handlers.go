package wallet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/wallet"
	"github.com/skwol/wallet/pkg/logging"
)

const (
	walletURL  = "/api/v1/wallets/{record_id}"
	walletsURL = "/api/v1/wallets"
)

type handler struct {
	walletService wallet.Service
}

func NewHandler(service wallet.Service) (adapters.Handler, error) {
	return &handler{walletService: service}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(walletURL, h.getWallet).Methods("GET")
	router.HandleFunc(walletURL, h.updateWallet).Methods("PATCH")
	router.HandleFunc(walletsURL, h.createWallet).Methods("POST")
	router.HandleFunc(walletsURL, h.getAllWallets).Methods("GET")
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
	response, err := json.Marshal(walletDTO)
	if err != nil {
		logger.Errorf("error marshaling wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling wallet: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(response)
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
	var request wallet.UpdateWalletDTO
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Errorf("error unmarshaling request: %s", err.Error())
		http.Error(w, fmt.Sprintf("error unmarshaling request: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	walletDTO, err := h.walletService.Update(r.Context(), id, &request)
	if err != nil {
		logger.Errorf("error updating wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error updating wallet: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}
	response, err := json.Marshal(walletDTO)
	if err != nil {
		logger.Errorf("error marshaling wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling wallet: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

func (h *handler) createWallet(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("error reading request body: %s", err.Error())
		http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	var request wallet.CreateWalletDTO
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Errorf("error unmarshaling request: %s", err.Error())
		http.Error(w, fmt.Sprintf("error unmarshaling request: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	walletDTO, err := h.walletService.Create(r.Context(), &request)
	if err != nil {
		logger.Errorf("error creating wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error updating wallet: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	response, err := json.Marshal(walletDTO)
	if err != nil {
		logger.Errorf("error marshaling wallet: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling wallet: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(response)
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

	wallets, err := h.walletService.GetAll(r.Context(), limit, offset)
	if err != nil {
		logger.Errorf("error returned from service: %s", err.Error())
		http.Error(w, fmt.Sprintf("error returned from service: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	type Data struct {
		Wallets []*wallet.WalletDTO `json:"wallets"`
	}
	response, err := json.Marshal(Data{Wallets: wallets})
	if err != nil {
		logger.Errorf("error marshaling wallets: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling wallets: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
