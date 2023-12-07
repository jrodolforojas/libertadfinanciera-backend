package transports

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	httptransport "github.com/go-kit/kit/transport/http"
	kitHTTP "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/configuration"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/middleware"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/utils"
)

type errorer interface {
	error() error
}

// MakeJSONEncoder creates a new json enooder
func MakeJSONEncoder(logger log.Logger) kitHTTP.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(errorer); ok && e.error() != nil {
			encodeError(ctx, e.error(), w, logger)
			return nil
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		return json.NewEncoder(w).Encode(response)
	}
}

func encodeError(_ context.Context, err error, w http.ResponseWriter, logger log.Logger) {
	if err == nil {
		_ = level.Error(logger).Log("msg", "error encoding error", "error", err)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
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
		MakeJSONEncoder(logger),
	))

	router.Methods(http.MethodGet).Path("/exchange_rates/today").Handler(httptransport.NewServer(
		endpoints.GetTodayExchangeRate,
		decodeTodayExchangeRateRequest,
		MakeJSONEncoder(logger),
	))

	router.Methods(http.MethodGet).Path("/exchange_rates/filter").Handler(httptransport.NewServer(
		endpoints.GetExchangeRatesByFilter,
		decodeGetDataByFilterRequest,
		MakeJSONEncoder(logger),
	))

	router.Methods(http.MethodGet).Path("/country_interes_rates/cr").Handler(httptransport.NewServer(
		endpoints.GetBasicPassiveRates,
		decodeGetAllDolarColonesChangesRequest,
		MakeJSONEncoder(logger),
	))

	router.Methods(http.MethodGet).Path("/country_interes_rates/cr/today").Handler(httptransport.NewServer(
		endpoints.GetTodayBasicPassiveRate,
		decodeTodayExchangeRateRequest,
		MakeJSONEncoder(logger),
	))
	router.Methods(http.MethodGet).Path("/country_interes_rates/usa").Handler(httptransport.NewServer(
		endpoints.GetTreasuryRatesUSA,
		decodeGetAllDolarColonesChangesRequest,
		MakeJSONEncoder(logger),
	))
	router.Methods(http.MethodGet).Path("/country_interes_rates/usa/today").Handler(httptransport.NewServer(
		endpoints.GetTreasuryRateUSA,
		decodeTodayExchangeRateRequest,
		MakeJSONEncoder(logger),
	))
	router.Methods(http.MethodGet).Path("/monetary_policy_rates").Handler(httptransport.NewServer(
		endpoints.GetMonetaryPolicyRates,
		decodeGetAllDolarColonesChangesRequest,
		MakeJSONEncoder(logger),
	))
	router.Methods(http.MethodGet).Path("/monetary_policy_rates/today").Handler(httptransport.NewServer(
		endpoints.GetTodayMonetaryPolicyRate,
		decodeTodayExchangeRateRequest,
		MakeJSONEncoder(logger),
	))
	router.Methods(http.MethodGet).Path("/prime_rates").Handler(httptransport.NewServer(
		endpoints.GetPrimeRates,
		decodeGetAllDolarColonesChangesRequest,
		MakeJSONEncoder(logger),
	))
	router.Methods(http.MethodGet).Path("/prime_rates/today").Handler(httptransport.NewServer(
		endpoints.GetPrimeRate,
		decodeTodayExchangeRateRequest,
		MakeJSONEncoder(logger),
	))

	router.Methods(http.MethodGet).Path("/inflation_rates/cr").Handler(httptransport.NewServer(
		endpoints.GetCostaRicaInflationRates,
		decodeGetAllDolarColonesChangesRequest,
		MakeJSONEncoder(logger),
	))

	router.Methods(http.MethodGet).Path("/inflation_rates/cr/filter").Handler(httptransport.NewServer(
		endpoints.GetCostaRicaInflationRatesByFilter,
		decodeGetDataByFilterRequest,
		MakeJSONEncoder(logger),
	))

	router.Methods(http.MethodGet).Path("/inflation_rates/cr/today").Handler(httptransport.NewServer(
		endpoints.GetCostaRicaInflationRate,
		decodeTodayExchangeRateRequest,
		MakeJSONEncoder(logger),
	))
	router.Methods(http.MethodGet).Path("/inflation_rates/usa").Handler(httptransport.NewServer(
		endpoints.GetUSAInflationRates,
		decodeGetAllDolarColonesChangesRequest,
		MakeJSONEncoder(logger),
	))
	router.Methods(http.MethodGet).Path("/inflation_rates/usa/today").Handler(httptransport.NewServer(
		endpoints.GetUSAInflationRate,
		decodeTodayExchangeRateRequest,
		MakeJSONEncoder(logger),
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
