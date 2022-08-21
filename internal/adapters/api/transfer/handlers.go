package transfer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/skwol/wallet/pkg/logging"

	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/transfer"
)

const transferURL = "/api/v1/transfers"

type handler struct {
	transferService transfer.Service
	logger          logging.Logger
}

func NewHandler(service transfer.Service, logger logging.Logger) (adapters.Handler, error) {
	return &handler{transferService: service, logger: logger}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(transferURL, h.createTransfer).Methods(http.MethodPost)
}

func (h *handler) createTransfer(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Errorf("error reading request body: %s", err.Error())
		http.Error(w, fmt.Sprintf("error reading request body: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	var request CreateTransferRequest
	if err := json.Unmarshal(body, &request); err != nil {
		h.logger.Errorf("error unmarshaling request: %s", err.Error())
		http.Error(w, fmt.Sprintf("error unmarshaling request: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	createRequest := request.toCreateRequest()
	transferDTO, err := h.transferService.Create(r.Context(), &createRequest)
	if err != nil {
		h.logger.Errorf("error creating transfer: %s", err.Error())
		http.Error(w, fmt.Sprintf("error creating wallet: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	response, err := json.Marshal(newTransfer(transferDTO))
	if err != nil {
		h.logger.Errorf("error marshaling transfer: %s", err.Error())
		http.Error(w, fmt.Sprintf("error marshaling transfer: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(response); err != nil {
		h.logger.Errorf("error writing response: %s", err.Error())
		http.Error(w, fmt.Sprintf("error writing response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
}
