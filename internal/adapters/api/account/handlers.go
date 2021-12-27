package account

import (
	"net/http"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/account"
)

const (
	getAccountURL  = "/api/v1/accounts/{record_id}"
	getAccountsURL = "/api/v1/accounts"
)

type handler struct {
	accountService account.Service
}

func NewHandler(service account.Service) adapters.Handler {
	return &handler{}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(getAccountURL, h.getAccount)
	router.HandleFunc(getAccountsURL, h.getAccounts)
}

func (h *handler) getAccount(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) getAccounts(w http.ResponseWriter, r *http.Request) {
	// h.accountService.GetAll(r.Context(), 0, 0)
	w.Write([]byte("accounts"))
	w.WriteHeader(http.StatusOK)
}
