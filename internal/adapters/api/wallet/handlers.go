package wallet

import (
	"net/http"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/wallet"
)

const (
	getWalletURL  = "/api/v1/wallets/{record_id}"
	getWalletsURL = "/api/v1/wallets"
)

type handler struct {
	transactionService wallet.Service
}

func NewHandler(service wallet.Service) (adapters.Handler, error) {
	return &handler{}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(getWalletURL, h.getWallet)
	router.HandleFunc(getWalletsURL, h.getWallets)
}

func (h *handler) getWallet(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) getWallets(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("wallets"))
	w.WriteHeader(http.StatusOK)
}
