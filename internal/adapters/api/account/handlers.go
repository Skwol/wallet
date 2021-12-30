package account

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/account"
	"github.com/skwol/wallet/pkg/logging"
)

const (
	accountURL  = "/api/v1/accounts/{record_id}"
	accountsURL = "/api/v1/accounts"
)

type handler struct {
	accountService account.Service
}

func NewHandler(service account.Service) (adapters.Handler, error) {
	return &handler{
		accountService: service,
	}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(accountURL, h.getAccount)
	router.HandleFunc(accountsURL, h.getAllAccounts)
}

func (h *handler) getAccount(w http.ResponseWriter, r *http.Request) {
}

func (h *handler) getAllAccounts(w http.ResponseWriter, r *http.Request) {
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

	accounts, err := h.accountService.GetAll(r.Context(), limit, offset)
	if err != nil {
		logger.Errorf("error returned from service: %s", err.Error())
		http.Error(w, fmt.Sprintf("error returned from service: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	type Data struct {
		Accounts []*account.AccountDTO `json:"accounts"`
	}
	response, err := json.Marshal(Data{Accounts: accounts})
	if err != nil {
		logger.Errorf("error marshaling accounts: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling accounts: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}
