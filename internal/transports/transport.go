package transports

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
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

type responseWriterWrapper struct {
	http.ResponseWriter
	wroteHeader bool
}

func (w *responseWriterWrapper) WriteHeader(statusCode int) {
	if !w.wroteHeader {
		w.ResponseWriter.WriteHeader(statusCode)
		w.wroteHeader = true
	}
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
	ww := &responseWriterWrapper{ResponseWriter: w}
	if !ww.wroteHeader {
		ww.Header().Set("Content-Type", "application/json; charset=utf-8")
		ww.WriteHeader(http.StatusInternalServerError) // Only call once
	}

	json.NewEncoder(ww).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case utils.ErrNotFound:
		return http.StatusNotFound
	case utils.ErrDateInvalidFormat, utils.ErrInvalidDateRange:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func MakeHTTPHandler(ctx context.Context, s *services.ServiceAPI, logger log.Logger) http.Handler {
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

	router.Methods(http.MethodGet).Path("/exchange_rates/filter").Handler(httptransport.NewServer(
		endpoints.GetExchangeRatesByFilter,
		decodeGetDataByFilterRequest,
		encodeResponse,
	))

	router.Methods(http.MethodGet).Path("/country_interes_rates/cr").Handler(httptransport.NewServer(
		endpoints.GetBasicPassiveRates,
		decodeGetAllDolarColonesChangesRequest,
		encodeResponse,
	))

	router.Methods(http.MethodGet).Path("/country_interes_rates/cr/today").Handler(httptransport.NewServer(
		endpoints.GetTodayBasicPassiveRate,
		decodeTodayExchangeRateRequest,
		encodeResponse,
	))
	router.Methods(http.MethodGet).Path("/country_interes_rates/usa").Handler(httptransport.NewServer(
		endpoints.GetTreasuryRatesUSA,
		decodeGetAllDolarColonesChangesRequest,
		encodeResponse,
	))
	router.Methods(http.MethodGet).Path("/country_interes_rates/usa/today").Handler(httptransport.NewServer(
		endpoints.GetTreasuryRateUSA,
		decodeTodayExchangeRateRequest,
		encodeResponse,
	))
	router.Methods(http.MethodGet).Path("/monetary_policy_rates").Handler(httptransport.NewServer(
		endpoints.GetMonetaryPolicyRates,
		decodeGetAllDolarColonesChangesRequest,
		encodeResponse,
	))
	router.Methods(http.MethodGet).Path("/monetary_policy_rates/today").Handler(httptransport.NewServer(
		endpoints.GetTodayMonetaryPolicyRate,
		decodeTodayExchangeRateRequest,
		encodeResponse,
	))
	router.Methods(http.MethodGet).Path("/prime_rates").Handler(httptransport.NewServer(
		endpoints.GetPrimeRates,
		decodeGetAllDolarColonesChangesRequest,
		encodeResponse,
	))
	router.Methods(http.MethodGet).Path("/prime_rates/today").Handler(httptransport.NewServer(
		endpoints.GetPrimeRate,
		decodeTodayExchangeRateRequest,
		encodeResponse,
	))

	router.Methods(http.MethodGet).Path("/inflation_rates/cr").Handler(httptransport.NewServer(
		endpoints.GetCostaRicaInflationRates,
		decodeGetAllDolarColonesChangesRequest,
		encodeResponse,
	))

	router.Methods(http.MethodGet).Path("/inflation_rates/cr/filter").Handler(httptransport.NewServer(
		endpoints.GetCostaRicaInflationRatesByFilter,
		decodeGetDataByFilterRequest,
		encodeResponse,
	))

	router.Methods(http.MethodGet).Path("/inflation_rates/cr/today").Handler(httptransport.NewServer(
		endpoints.GetCostaRicaInflationRate,
		decodeTodayExchangeRateRequest,
		encodeResponse,
	))
	router.Methods(http.MethodGet).Path("/inflation_rates/usa").Handler(httptransport.NewServer(
		endpoints.GetUSAInflationRates,
		decodeGetAllDolarColonesChangesRequest,
		encodeResponse,
	))
	router.Methods(http.MethodGet).Path("/inflation_rates/usa/today").Handler(httptransport.NewServer(
		endpoints.GetUSAInflationRate,
		decodeTodayExchangeRateRequest,
		encodeResponse,
	))
	return router
}

func decodeGetDataByFilterRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	periocity := r.FormValue("periocity")

	if periocity == "" {
		return services.GetDataByFilterRequest{
			Periodicity: "monthly",
		}, nil
	}

	if periocity != "quarterly" && periocity != "biannual" && periocity != "annual" &&
		periocity != "quinquennium" && periocity != "monthly" {
		return nil, utils.ErrPeriodicity
	}

	return services.GetDataByFilterRequest{
		Periodicity: periocity,
	}, nil
}

func decodeGetAllDolarColonesChangesRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	dateFromParam := r.FormValue("date_from")
	dateToParam := r.FormValue("date_to")

	if dateFromParam != "" && dateToParam != "" {
		dateFrom, error := utils.ConvertStringDate(dateFromParam)
		if error != nil {
			return nil, utils.ErrDateInvalidFormat
		}

		dateTo, error := utils.ConvertStringDate(dateToParam)
		if error != nil {
			return nil, utils.ErrDateInvalidFormat
		}

		if !utils.IsDatesValid(dateFrom, dateTo) {
			return nil, utils.ErrInvalidDateRange
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
