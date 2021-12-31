package common

import (
	"net/http"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/common"
	"github.com/skwol/wallet/pkg/logging"
)

const (
	generateFakeDataURL = "/api/v1/generate_fake_data"
)

type handler struct {
	service common.Service
}

func NewHandler(service common.Service) (adapters.Handler, error) {
	return &handler{service: service}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(generateFakeDataURL, h.generateFakeData).Methods("POST")
}

func (h *handler) generateFakeData(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	err := h.service.GenerateFakeData(r.Context())
	if err != nil {
		logger.Errorf("error during generating data: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
