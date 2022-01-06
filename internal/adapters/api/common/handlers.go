package common

import (
	"net/http"

	"github.com/gorilla/mux"
	adapters "github.com/skwol/wallet/internal/adapters/api"
	"github.com/skwol/wallet/internal/domain/common"
	"github.com/skwol/wallet/pkg/logging"
	"golang.org/x/time/rate"
)

const (
	generateFakeDataURL = "/api/v1/generate_fake_data"
)

// ~ 1 per 25 minutes
var limiter = rate.NewLimiter(0.0007, 1)

type handler struct {
	service common.Service
}

func NewHandler(service common.Service) (adapters.Handler, error) {
	return &handler{service: service}, nil
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(generateFakeDataURL, limit(h.generateFakeData)).Methods("POST")
}

func limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}
		next(w, r)
	}
}

func (h *handler) generateFakeData(w http.ResponseWriter, r *http.Request) {
	logger := logging.GetLogger()
	ctx := r.Context()
	go func() {
		if err := h.service.GenerateFakeData(ctx); err != nil {
			logger.Errorf("error during generating data: %s", err.Error())
			return
		}
	}()
	w.WriteHeader(http.StatusCreated)
}
