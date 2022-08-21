package common

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/time/rate"

	"github.com/skwol/wallet/pkg/logging"

	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/common"
)

const (
	generateFakeDataURL = "/api/v1/generate_fake_data"
)

// ~ 1 per 25 minutes
var limiter = rate.NewLimiter(0.0007, 1)

type handler struct {
	service common.Service
	logger  logging.Logger
}

func NewHandler(service common.Service, logger logging.Logger) (adapters.Handler, error) {
	return &handler{service: service, logger: logger}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(generateFakeDataURL, h.generateFakeData).Methods("POST")
}

func (h *handler) generateFakeData(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	ctx := r.Context()

	numberOfRecordsToCreate, err := strconv.Atoi(r.FormValue("records"))
	if err != nil {
		logger.Errorf("error parsing records query param: %s", err.Error())
		http.Error(w, fmt.Sprintf("error parsing records query param: %s", err.Error()), http.StatusUnprocessableEntity)
		return
	}

	if !limiter.Allow() {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	go func() {
		if err := h.service.GenerateFakeData(ctx, numberOfRecordsToCreate); err != nil {
			logger.Errorf("error during generating data: %s", err.Error())
			return
		}
	}()
	w.WriteHeader(http.StatusCreated)
}
