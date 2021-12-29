package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/wallet"
	"github.com/skwol/wallet/pkg/logging"
)

const (
	getWalletURL  = "/api/v1/wallets/{record_id}"
	getWalletsURL = "/api/v1/wallets"
)

type handler struct {
	walletService wallet.Service
}

func NewHandler(service wallet.Service) (adapters.Handler, error) {
	return &handler{walletService: service}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(getWalletURL, h.getWallet)
	router.HandleFunc(getWalletsURL, h.getAllWallets)
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
