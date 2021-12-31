package transfer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/transfer"
	"github.com/skwol/wallet/pkg/logging"
)

const transferURL = "/api/v1/transfers"

type handler struct {
	transferService transfer.Service
}

func NewHandler(service transfer.Service) (adapters.Handler, error) {
	return &handler{transferService: service}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(transferURL, h.createTransfer).Methods("POST")
}

func (h *handler) createTransfer(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("error reading request body: %s", err.Error())
		http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	var request transfer.CreateTransferDTO
	if err := json.Unmarshal(body, &request); err != nil {
		logger.Errorf("error unmarshaling request: %s", err.Error())
		http.Error(w, fmt.Sprintf("error unmarshaling request: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	transferDTO, err := h.transferService.Create(r.Context(), &request)
	if err != nil {
		logger.Errorf("error creating transfer: %s", err.Error())
		http.Error(w, fmt.Sprintf("error creating wallet: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	response, err := json.Marshal(transferDTO)
	if err != nil {
		logger.Errorf("error marshaling transfer: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling transfer: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
