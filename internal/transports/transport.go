package transports

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/middleware"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services"
)

type errorer interface {
	error() error
}

func MakeHTTPHandler(ctx context.Context, s *services.Service) http.Handler {
	router := mux.NewRouter()
	endpoints := services.MakeEndpoints(s)

	subRouter := router.PathPrefix("/api").Subrouter()
	router = router.PathPrefix("/").Subrouter()

	corsMethods := []string{http.MethodOptions, http.MethodGet}
	router.Use(middleware.CORSPolicies(corsMethods))
	subRouter.Use(middleware.CORSPolicies(corsMethods))

	router.Methods(http.MethodGet).Path("/exchange_rates").Handler(httptransport.NewServer(
		endpoints.GetAllDolarColonesChanges,
		decodeGetAllDolarColonesChangesRequest,
		encodeResponse,
	))

	router.Methods(http.MethodGet).Path("/exchange_rates/today").Handler(httptransport.NewServer(
		endpoints.GetTodayExchangeRate,
		decodeTodayExchangeRateRequest,
		encodeResponse,
	))

	return router
}

func decodeGetAllDolarColonesChangesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req services.GetAllDolarColonesChangesRequest
	return req, nil
}

func decodeTodayExchangeRateRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req services.GetTodayExchangeRateRequest
	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		log.Println("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

var (
	// ErrInconsistentIDs inconsistent IDs error returns int error.
	errInconsistentIDs = errors.New("inconsistent IDs")
	// ErrAlreadyExists already exists error returns int error.
	errAlreadyExists = errors.New("already exists")
	// ErrNotFound not found error returns int error.
	errNotFound = errors.New("not found")
)

func codeFrom(err error) int {
	switch err {
	case errNotFound:
		return http.StatusNotFound
	case errAlreadyExists, errInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
