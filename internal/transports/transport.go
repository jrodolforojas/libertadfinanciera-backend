package transports

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/configuration"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/middleware"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/utils"
)

type errorer interface {
	error() error
}

func MakeHTTPHandler(ctx context.Context, s *services.ServiceAPI) http.Handler {
	router := mux.NewRouter()
	endpoints := services.MakeEndpoints(s)

	subRouter := router.PathPrefix("/api").Subrouter()
	router = router.PathPrefix("/").Subrouter()

	corsMethods := []string{http.MethodOptions, http.MethodGet}
	config, _ := configuration.Read()
	router.Use(middleware.CORSPolicies(corsMethods, config.Address.AllowedOrigins))
	subRouter.Use(middleware.CORSPolicies(corsMethods, config.Address.AllowedOrigins))

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
	dateFromParam := r.FormValue("date_from")
	dateToParam := r.FormValue("date_to")

	if dateFromParam != "" && dateToParam != "" {
		dateFrom, error := utils.ConvertStringDate(dateFromParam)
		if error != nil {
			return nil, errDateInvalidFormat
		}

		dateTo, error := utils.ConvertStringDate(dateToParam)
		if error != nil {
			return nil, errDateInvalidFormat
		}

		if !utils.IsDatesValid(dateFrom, dateTo) {
			return nil, errInvalidDateRange
		}

		return services.GetAllDollarColonesChangesRequest{
			DateFrom: dateFrom,
			DateTo:   dateTo,
		}, nil
	}

	dateFrom, dateTo := utils.GetDateFromDateToFromToday(utils.DEFAULT_DAYS_TO_GO_BACK)
	return services.GetAllDollarColonesChangesRequest{
		DateFrom: dateFrom,
		DateTo:   dateTo,
	}, nil
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
	errDateInvalidFormat = errors.New("invalid date format. Should be in format: YYYY/MM/DD")
	errInvalidDateRange  = errors.New("invalid date range. Should be between 1983 and today")
	errNotFound          = errors.New("not found")
)

func codeFrom(err error) int {
	switch err {
	case errNotFound:
		return http.StatusNotFound
	case errDateInvalidFormat, errInvalidDateRange:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
