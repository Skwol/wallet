package wallet

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

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
	tmpl, err := template.ParseFiles("templates/wallets/wallet.html")
	if err != nil {
		logger.Errorf("error parsing wallet template: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing wallet template: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if tmpl == nil {
		logger.Error("missing wallet template")
		http.Error(w, "missing wallet template", http.StatusInternalServerError)
		return
	}
	type Data struct {
		Wallet *wallet.WalletDTO
	}

	tmpl.Execute(w, Data{Wallet: walletDTO})
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
	tmpl, err := template.ParseFiles("templates/wallets/wallets.html")
	if err != nil {
		logger.Errorf("error parsing wallets template: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing wallets template: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if tmpl == nil {
		logger.Error("missing wallets template")
		http.Error(w, "missing wallets template", http.StatusInternalServerError)
		return
	}
	type Data struct {
		Wallets []*wallet.WalletDTO
	}
	tmpl.Execute(w, Data{Wallets: wallets})
}
