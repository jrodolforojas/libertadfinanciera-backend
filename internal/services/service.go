package services

import (
	"context"
	"errors"
	"sort"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/repositories"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services/scrapper"
)

type Service interface {
	GetDollarColonesChange(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetAllDollarColonesChangesResponse
	GetExchangeRatesByFilter(ctx context.Context, req GetDataByFilterRequest) *GetAllDollarColonesChangesResponse
	GetTodayExchangeRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayExchangeRateResponse
	GetBasicPassiveRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetBasicPassiveRatesResponse
	GetTodayBasicPassiveRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayBasicPassiveRateResponse
	GetMonetaryPolicyRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetMonetaryPolicyRatesResponse
	GetTodayMonetaryPolicyRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayMonetaryPolicyRateResponse
	GetPrimeRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetPrimeRatesResponse
	GetTodayPrimeRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayPrimeRateResponse
	GetCostaRicaInflationRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetCostaRicaInflationRatesResponse
	GetCostaRicaInflationRatesByFilter(ctx context.Context, req GetDataByFilterRequest) *GetCostaRicaInflationRatesResponse
	GetCostaRicaInflationRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayCostaRicaInflationRateResponse
	GetTreasuryRatesUSA(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetTreasuryRatesUSAResponse
	GetTodayTreasuryRateUSA(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayTreasuryRateUSAResponse
	GetUSAInflationRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetUSAInflationRatesResponse
	GetUSAInflationRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayUSAInflationRateResponse
}

type ServiceAPI struct {
	logger     log.Logger
	Scrapper   scrapper.Scrapper
	Repository repositories.Repository
}

func NewService(logger log.Logger, scrapperService scrapper.Scrapper, repo repositories.Repository) *ServiceAPI {
	return &ServiceAPI{
		logger:     logger,
		Scrapper:   scrapperService,
		Repository: repo,
	}
}

func (service *ServiceAPI) GetDollarColonesChange(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetAllDollarColonesChangesResponse {
	exchangeRates := []models.ExchangeRate{}

	dateFrom := req.DateFrom

	dateRanges := []models.DateRange{}
	for {
		if dateFrom.Month() == req.DateTo.Month() && dateFrom.Year() == req.DateTo.Year() {
			dateRanges = append(dateRanges, models.DateRange{
				DateFrom: dateFrom,
				DateTo:   req.DateTo,
			})
			break
		}
		dateTo := dateFrom.AddDate(0, 1, 0) // add 1 month to dateFrom
		dateRanges = append(dateRanges, models.DateRange{
			DateFrom: dateFrom,
			DateTo:   dateTo,
		})
		dateFrom = dateTo
	}

	errc := make(chan error, len(dateRanges))
	for _, dateRange := range dateRanges {
		go func(dateFrom time.Time, dateTo time.Time) {
			result, err := service.Scrapper.GetDollarColonesChangeByDates(dateFrom, dateTo, 0)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping exchange rate by dates",
					"date_from", dateFrom, "date_to", dateTo)
				errc <- err
				return
			}
			exchangeRates = append(exchangeRates, result...)
			errc <- nil
		}(dateRange.DateFrom, dateRange.DateTo)
	}

	for i := 0; i < len(dateRanges); i++ {
		if err := <-errc; err != nil {
			return &GetAllDollarColonesChangesResponse{
				ExchangesRates: nil,
				Err:            err,
			}
		}
	}

	sort.Slice(exchangeRates, func(i, j int) bool {
		return exchangeRates[i].Date.After(exchangeRates[j].Date)
	})

	return &GetAllDollarColonesChangesResponse{
		ExchangesRates: exchangeRates,
		Err:            nil,
	}
}

func (service *ServiceAPI) GetExchangeRatesByFilter(ctx context.Context, req GetDataByFilterRequest) *GetAllDollarColonesChangesResponse {
	filtersArray := []int64{}

	if req.Periodicity == "quarterly" {
		filtersArray = []int64{31, 91, 182, 274, 366}
	}
	if req.Periodicity == "biannual" {
		filtersArray = []int64{182, 366}
	}
	if req.Periodicity == "annual" || req.Periodicity == "quinquennium" {
		filtersArray = []int64{366}
	}

	exchangeRates := []models.ExchangeRate{}

	minimumDate := time.Date(1983, 0, 31, 0, 0, 0, 0, time.UTC)
	bridgeDate := time.Date(2003, 0, 1, 0, 0, 0, 0, time.UTC)
	today := time.Now()

	errc := make(chan error, len(filtersArray))

	for _, filter := range filtersArray {
		go func(filter int64) {
			result, err := service.Scrapper.GetDollarColonesChangeByDates(minimumDate, bridgeDate, filter)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping inflation rate by filter",
					"date_from", minimumDate, "date_to", bridgeDate, "filter", filter, "error", err)
				errc <- err
				return
			}
			exchangeRates = append(exchangeRates, result...)
			errc <- nil

		}(filter)
	}

	errcBridgeToday := make(chan error, len(filtersArray))
	for _, filter := range filtersArray {
		go func(filter int64) {
			result, err := service.Scrapper.GetDollarColonesChangeByDates(bridgeDate, today, filter)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping inflation rate by filter",
					"date_from", bridgeDate, "date_to", today, "filter", filter, "error", err)
				errcBridgeToday <- err
				return
			}
			exchangeRates = append(exchangeRates, result...)
			errcBridgeToday <- nil

		}(filter)
	}

	for i := 0; i < len(filtersArray); i++ {
		if err := <-errc; err != nil {
			return &GetAllDollarColonesChangesResponse{
				ExchangesRates: nil,
				Err:            err,
			}
		}
	}

	for i := 0; i < len(filtersArray); i++ {
		if err := <-errcBridgeToday; err != nil {
			return &GetAllDollarColonesChangesResponse{
				ExchangesRates: nil,
				Err:            err,
			}
		}
	}

	if req.Periodicity == "quinquennium" {
		// filter exchanges rates by quinquennium from 1976 to today
		exchangeRatesQuinquennium := []models.ExchangeRate{}
		for i := 0; i < len(exchangeRates); i++ {
			if i%5 == 0 {
				exchangeRatesQuinquennium = append(exchangeRatesQuinquennium, exchangeRates[i])
			}
		}

		exchangeRates = exchangeRatesQuinquennium
	}

	sort.Slice(exchangeRates, func(i, j int) bool {
		return exchangeRates[i].Date.After(exchangeRates[j].Date)
	})

	return &GetAllDollarColonesChangesResponse{
		ExchangesRates: exchangeRates,
		Err:            nil,
	}
}
func (service *ServiceAPI) GetTodayExchangeRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayExchangeRateResponse {
	date := time.Now()
	for i := 0; i < MAXIMUM_TRIES; i++ {
		rate, err := service.Scrapper.GetExchangeRateByDate(date)
		if rate.BuyPrice != 0 && rate.SalePrice != 0 {
			return &GetTodayExchangeRateResponse{
				ExchangesRate: rate,
				Err:           err,
			}
		}
		date = date.AddDate(0, 0, -1)
	}
	_ = level.Error(service.logger).Log("msg", "error scrapping USA treasury rate by date", "date", date)
	return &GetTodayExchangeRateResponse{
		ExchangesRate: nil,
		Err:           errors.New("no results found"),
	}
}

const MAXIMUM_BASIC_PASSIVE_RATE_YEAR = 12

func (service *ServiceAPI) addYears(dateFrom time.Time, years int) time.Time {
	return dateFrom.AddDate(years-1, 0, 0)
}
func (service *ServiceAPI) GetBasicPassiveRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetBasicPassiveRatesResponse {
	yearDifference := req.DateTo.Year() - req.DateFrom.Year()
	var basicPassiveRates []models.BasicPassiveRate

	dateFrom := req.DateFrom
	for {
		if yearDifference >= MAXIMUM_BASIC_PASSIVE_RATE_YEAR {
			newDateTo := service.addYears(dateFrom, MAXIMUM_BASIC_PASSIVE_RATE_YEAR)
			result, err := service.Scrapper.GetBasicPassiveRateByDates(dateFrom, newDateTo)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping basic passive rates by dates",
					"date_from", dateFrom, "date_to", newDateTo)
				break
			}
			basicPassiveRates = append(basicPassiveRates, result...)
			dateFrom = newDateTo.AddDate(0, 0, 1)
			yearDifference = req.DateTo.Year() - dateFrom.Year()

		} else {
			result, err := service.Scrapper.GetBasicPassiveRateByDates(dateFrom, req.DateTo)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping basic passive rates by dates",
					"date_from", dateFrom, "date_to", req.DateTo)
				break
			}
			basicPassiveRates = append(basicPassiveRates, result...)
			break
		}
	}

	sort.Slice(basicPassiveRates, func(i, j int) bool {
		return basicPassiveRates[i].Date.After(basicPassiveRates[j].Date)
	})
	return &GetBasicPassiveRatesResponse{
		BasicPassiveRates: basicPassiveRates,
		Err:               nil,
	}
}

func (service *ServiceAPI) GetTodayBasicPassiveRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayBasicPassiveRateResponse {
	date := time.Now()
	for i := 0; i < MAXIMUM_TRIES; i++ {
		rate, err := service.Scrapper.GetBasicPassiveDateByDate(date)
		if rate.Value != 0 {
			return &GetTodayBasicPassiveRateResponse{
				BasicPassiveRate: rate,
				Err:              err,
			}
		}
		date = date.AddDate(0, 0, -1)
	}
	_ = level.Error(service.logger).Log("msg", "error scrapping USA treasury rate by date", "date", date)
	return &GetTodayBasicPassiveRateResponse{
		BasicPassiveRate: nil,
		Err:              errors.New("no results found"),
	}
}

func (service *ServiceAPI) GetMonetaryPolicyRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetMonetaryPolicyRatesResponse {
	yearDifference := req.DateTo.Year() - req.DateFrom.Year()
	var monetaryPolicyRates []models.MonetaryPolicyRate

	const MAXIMUM_MONETARY_POLICY_RATE_YEAR = 5
	dateFrom := req.DateFrom
	for {
		if yearDifference >= MAXIMUM_MONETARY_POLICY_RATE_YEAR {
			newDateTo := service.addYears(dateFrom, MAXIMUM_MONETARY_POLICY_RATE_YEAR)
			result, err := service.Scrapper.GetMonetaryPolicyRateByDates(dateFrom, newDateTo)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping basic passive rates by dates",
					"date_from", dateFrom, "date_to", newDateTo)
				break
			}
			monetaryPolicyRates = append(monetaryPolicyRates, result...)
			dateFrom = newDateTo.AddDate(0, 0, 1)
			yearDifference = req.DateTo.Year() - dateFrom.Year()

		} else {
			result, err := service.Scrapper.GetMonetaryPolicyRateByDates(dateFrom, req.DateTo)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping basic passive rates by dates",
					"date_from", dateFrom, "date_to", req.DateTo)
				break
			}
			monetaryPolicyRates = append(monetaryPolicyRates, result...)
			break
		}
	}

	sort.Slice(monetaryPolicyRates, func(i, j int) bool {
		return monetaryPolicyRates[i].Date.After(monetaryPolicyRates[j].Date)
	})
	return &GetMonetaryPolicyRatesResponse{
		MonetaryPolicyRates: monetaryPolicyRates,
		Err:                 nil,
	}
}

func (service *ServiceAPI) GetTodayMonetaryPolicyRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayMonetaryPolicyRateResponse {
	date := time.Now()
	for i := 0; i < MAXIMUM_TRIES; i++ {
		rate, err := service.Scrapper.GetMonetaryPolicyRateByDate(date)
		if rate.Value != 0 {
			return &GetTodayMonetaryPolicyRateResponse{
				MonetaryPolicyRate: rate,
				Err:                err,
			}
		}
		date = date.AddDate(0, 0, -1)
	}
	_ = level.Error(service.logger).Log("msg", "error scrapping USA treasury rate by date", "date", date)
	return &GetTodayMonetaryPolicyRateResponse{
		MonetaryPolicyRate: nil,
		Err:                errors.New("no results found"),
	}
}

func (service *ServiceAPI) GetPrimeRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetPrimeRatesResponse {
	yearDifference := req.DateTo.Year() - req.DateFrom.Year()
	var primeRates []models.PrimeRate

	const MAXIMUM_PRIME_RATE_YEAR = 9
	dateFrom := req.DateFrom
	for {
		if yearDifference >= MAXIMUM_PRIME_RATE_YEAR {
			newDateTo := service.addYears(dateFrom, MAXIMUM_PRIME_RATE_YEAR)
			result, err := service.Scrapper.GetPrimeRateByDates(dateFrom, newDateTo)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping basic passive rates by dates",
					"date_from", dateFrom, "date_to", newDateTo)
				break
			}
			primeRates = append(primeRates, result...)
			dateFrom = newDateTo.AddDate(0, 0, 1)
			yearDifference = req.DateTo.Year() - dateFrom.Year()

		} else {
			result, err := service.Scrapper.GetPrimeRateByDates(dateFrom, req.DateTo)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping basic passive rates by dates",
					"date_from", dateFrom, "date_to", req.DateTo)
				break
			}
			primeRates = append(primeRates, result...)
			break
		}
	}

	sort.Slice(primeRates, func(i, j int) bool {
		return primeRates[i].Date.After(primeRates[j].Date)
	})
	return &GetPrimeRatesResponse{
		PrimeRates: primeRates,
		Err:        nil,
	}
}

const MAXIMUM_TRIES = 5

func (service *ServiceAPI) GetTodayPrimeRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayPrimeRateResponse {
	date := time.Now()
	for i := 0; i < MAXIMUM_TRIES; i++ {
		primeRate, err := service.Scrapper.GetPrimeRateByDate(date)
		if primeRate.Value != 0 {
			return &GetTodayPrimeRateResponse{
				PrimeRate: primeRate,
				Err:       err,
			}
		}
		date = date.AddDate(0, 0, -1)
	}
	_ = level.Error(service.logger).Log("msg", "error scrapping USA treasury rate by date", "date", date)
	return &GetTodayPrimeRateResponse{
		PrimeRate: nil,
		Err:       errors.New("no results found"),
	}
}

const MAXIMUM_INFLATION_RATE_YEAR = 2

func (service *ServiceAPI) GetCostaRicaInflationRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetCostaRicaInflationRatesResponse {
	inflationRates := []models.CostaRicaInflationRate{}

	dateFrom := req.DateFrom
	dateRanges := []models.DateRange{}

	yearDifference := req.DateTo.Year() - req.DateFrom.Year()

	for {
		if yearDifference <= MAXIMUM_INFLATION_RATE_YEAR {
			dateRanges = append(dateRanges, models.DateRange{
				DateFrom: dateFrom,
				DateTo:   req.DateTo,
			})
			break
		} else {
			// add 1 year to dateFrom
			dateTo := dateFrom.AddDate(MAXIMUM_INFLATION_RATE_YEAR, 0, 0)
			dateRanges = append(dateRanges, models.DateRange{
				DateFrom: dateFrom,
				DateTo:   dateTo,
			})
			dateFrom = dateTo
			yearDifference = req.DateTo.Year() - dateFrom.Year()
		}
	}

	errc := make(chan error, len(dateRanges))
	for _, dateRange := range dateRanges {
		go func(dateFrom time.Time, dateTo time.Time) {
			result, err := service.Scrapper.GetCostaRicaInflationRateByDates(dateFrom, dateTo, 0)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping inflation rate by dates",
					"date_from", dateFrom, "date_to", dateTo)
				errc <- err
				return
			}
			inflationRates = append(inflationRates, result...)
			errc <- nil
		}(dateRange.DateFrom, dateRange.DateTo)
	}

	for i := 0; i < len(dateRanges); i++ {
		if err := <-errc; err != nil {
			return &GetCostaRicaInflationRatesResponse{
				InflationRates: nil,
				Err:            err,
			}
		}
	}

	sort.Slice(inflationRates, func(i, j int) bool {
		return inflationRates[i].Date.After(inflationRates[j].Date)
	})
	return &GetCostaRicaInflationRatesResponse{
		InflationRates: inflationRates,
		Err:            nil,
	}
}

func (service *ServiceAPI) GetCostaRicaInflationRatesByFilter(ctx context.Context, req GetDataByFilterRequest) *GetCostaRicaInflationRatesResponse {
	filtersArray := []int64{}

	if req.Periodicity == "quarterly" {
		filtersArray = []int64{1, 3, 6, 9, 12}
	}
	if req.Periodicity == "biannual" {
		filtersArray = []int64{6, 12}
	}
	if req.Periodicity == "annual" || req.Periodicity == "quinquennium" {
		filtersArray = []int64{12}
	}

	inflationRates := []models.CostaRicaInflationRate{}

	minimumDate := time.Date(1976, 0, 31, 0, 0, 0, 0, time.UTC)
	bridgeDate := time.Date(1995, 0, 1, 0, 0, 0, 0, time.UTC)
	today := time.Now()

	errc := make(chan error, len(filtersArray))

	for _, filter := range filtersArray {
		go func(filter int64) {
			result, err := service.Scrapper.GetCostaRicaInflationRateByDates(minimumDate, bridgeDate, filter)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping inflation rate by filter",
					"date_from", minimumDate, "date_to", bridgeDate, "filter", filter, "error", err)
				errc <- err
				return
			}
			inflationRates = append(inflationRates, result...)
			errc <- nil

		}(filter)
	}

	errcBridgeToday := make(chan error, len(filtersArray))
	for _, filter := range filtersArray {
		go func(filter int64) {
			result, err := service.Scrapper.GetCostaRicaInflationRateByDates(bridgeDate, today, filter)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping inflation rate by filter",
					"date_from", bridgeDate, "date_to", today, "filter", filter, "error", err)
				errcBridgeToday <- err
				return
			}
			inflationRates = append(inflationRates, result...)
			errcBridgeToday <- nil

		}(filter)
	}

	for i := 0; i < len(filtersArray); i++ {
		if err := <-errc; err != nil {
			return &GetCostaRicaInflationRatesResponse{
				InflationRates: nil,
				Err:            err,
			}
		}
	}

	for i := 0; i < len(filtersArray); i++ {
		if err := <-errcBridgeToday; err != nil {
			return &GetCostaRicaInflationRatesResponse{
				InflationRates: nil,
				Err:            err,
			}
		}
	}

	if req.Periodicity == "quinquennium" {
		// filter exchanges rates by quinquennium from 1976 to today
		inflationRatesQuinquennium := []models.CostaRicaInflationRate{}
		for i := 0; i < len(inflationRates); i++ {
			if i%5 == 0 {
				inflationRatesQuinquennium = append(inflationRatesQuinquennium, inflationRates[i])
			}
		}

		inflationRates = inflationRatesQuinquennium
	}

	sort.Slice(inflationRates, func(i, j int) bool {
		return inflationRates[i].Date.After(inflationRates[j].Date)
	})

	return &GetCostaRicaInflationRatesResponse{
		InflationRates: inflationRates,
		Err:            nil,
	}
}

func (service *ServiceAPI) GetCostaRicaInflationRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayCostaRicaInflationRateResponse {
	date := time.Now()
	for i := 0; i < MAXIMUM_TRIES; i++ {
		inflationRate, err := service.Scrapper.GetCostaRicaInflationRateByDate(date)
		if inflationRate.Value != 0 {
			return &GetTodayCostaRicaInflationRateResponse{
				InflationRate: inflationRate,
				Err:           err,
			}
		}
		date = date.AddDate(0, -1, 0)
	}
	_ = level.Error(service.logger).Log("msg", "error scrapping Costa Rica inflation rate by date", "date", date)
	return &GetTodayCostaRicaInflationRateResponse{
		InflationRate: nil,
		Err:           errors.New("no results found"),
	}
}

func (service *ServiceAPI) GetTreasuryRatesUSA(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetTreasuryRatesUSAResponse {
	treasuryRates := []models.TreasuryRateUSA{}
	dateFrom := req.DateFrom
	dateRanges := []models.DateRange{}
	for {
		if dateFrom.Month() == req.DateTo.Month() && dateFrom.Year() == req.DateTo.Year() {
			dateRanges = append(dateRanges, models.DateRange{
				DateFrom: dateFrom,
				DateTo:   req.DateTo,
			})
			break
		}
		dateTo := dateFrom.AddDate(0, 1, 0) // add 1 month to dateFrom
		dateRanges = append(dateRanges, models.DateRange{
			DateFrom: dateFrom,
			DateTo:   dateTo,
		})
		dateFrom = dateTo
	}

	errc := make(chan error, len(dateRanges))
	for _, dateRange := range dateRanges {
		go func(dateFrom time.Time, dateTo time.Time) {
			result, err := service.Scrapper.GetTreasuryRateUSAByDates(dateFrom, dateTo)
			if err != nil {
				_ = level.Debug(service.logger).Log("msg", "error scrapping USA treasury rate by dates", "error", err)
				errc <- err
				return
			}
			treasuryRates = append(treasuryRates, result...)
			errc <- nil
		}(dateRange.DateFrom, dateRange.DateTo)
	}

	for i := 0; i < len(dateRanges); i++ {
		if err := <-errc; err != nil {
			return &GetTreasuryRatesUSAResponse{
				TreasuryRatesUSA: nil,
				Err:              err,
			}
		}
	}

	sort.Slice(treasuryRates, func(i, j int) bool {
		return treasuryRates[i].Date.After(treasuryRates[j].Date)
	})

	return &GetTreasuryRatesUSAResponse{
		TreasuryRatesUSA: treasuryRates,
		Err:              nil,
	}
}

func (service *ServiceAPI) GetTodayTreasuryRateUSA(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayTreasuryRateUSAResponse {
	date := time.Now()
	const treasuryRateMaxiumTries = 30
	for i := 0; i < treasuryRateMaxiumTries; i++ {
		todayTreasuryRateUSA, err := service.Scrapper.GetTreasuryRateUSAByDate(date)
		if todayTreasuryRateUSA.Value != 0 {
			return &GetTodayTreasuryRateUSAResponse{
				TreasuryRateUSA: todayTreasuryRateUSA,
				Err:             err,
			}
		}
		date = date.AddDate(0, 0, -1)
	}
	_ = level.Error(service.logger).Log("msg", "error scrapping USA treasury rate by date", "date", date)
	return &GetTodayTreasuryRateUSAResponse{
		TreasuryRateUSA: nil,
		Err:             errors.New("no results found"),
	}
}

func (service *ServiceAPI) GetUSAInflationRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetUSAInflationRatesResponse {
	result, err := service.Scrapper.GetUSAInflationRateByDates(req.DateFrom, req.DateTo)
	if err != nil {
		_ = level.Debug(service.logger).Log("msg", "error scrapping USA inflation rate by dates",
			"date_from", req.DateFrom, "date_to", req.DateTo)
		return &GetUSAInflationRatesResponse{
			InflationRates: nil,
			Err:            err,
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Date.After(result[j].Date)
	})

	return &GetUSAInflationRatesResponse{
		InflationRates: result,
		Err:            nil,
	}
}

func (service *ServiceAPI) GetUSAInflationRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayUSAInflationRateResponse {
	date := time.Now()
	inflationRate, err := service.Scrapper.GetUSAInflationRateByDate(date)
	return &GetTodayUSAInflationRateResponse{
		InflationRate: inflationRate,
		Err:           err,
	}
}
