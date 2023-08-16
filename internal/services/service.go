package services

import (
	"context"
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
	GetTodayExchangeRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayExchangeRateResponse
	GetBasicPassiveRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetBasicPassiveRatesResponse
	GetTodayBasicPassiveRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayBasicPassiveRateResponse
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
	dateFrom := req.DateFrom

	var allExchangeRates []models.ExchangeRate
	for dateFrom.Before(req.DateTo) {
		nextMonthDate := dateFrom.AddDate(0, 1, 0) // add 1 month to dateFrom
		if nextMonthDate.After(req.DateTo) {
			nextMonthDate = req.DateTo
		}
		exchangeRates, err := service.Scrapper.GetDollarColonesChangeByDates(dateFrom, nextMonthDate)
		if err != nil {
			_ = level.Error(service.logger).Log("msg", "error scrapping", "dateFrom", dateFrom, "dateTo", nextMonthDate, "error", err)
			return &GetAllDollarColonesChangesResponse{
				ExchangesRates: nil,
				Err:            err,
			}
		}
		allExchangeRates = append(allExchangeRates, exchangeRates...)
		dateFrom = nextMonthDate
	}

	sort.Slice(allExchangeRates, func(i, j int) bool {
		return allExchangeRates[i].Date.After(allExchangeRates[j].Date)
	})

	return &GetAllDollarColonesChangesResponse{
		ExchangesRates: allExchangeRates,
		Err:            nil,
	}
}

func (service *ServiceAPI) GetTodayExchangeRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayExchangeRateResponse {
	date := time.Now()
	todayExchangeRate, err := service.Scrapper.GetExchangeRateByDate(date)
	if err != nil {
		_ = level.Error(service.logger).Log("msg", "error scrapping", "dateFrom", date, "dateTo", date, "error", err)
		return &GetTodayExchangeRateResponse{
			ExchangesRate: nil,
			Err:           err,
		}
	}

	return &GetTodayExchangeRateResponse{
		ExchangesRate: todayExchangeRate,
		Err:           nil,
	}
}

const MAXIMUM_BASIC_PASSIVE_RATE_YEAR = 12

func (service *ServiceAPI) add12years(dateFrom time.Time) time.Time {
	return dateFrom.AddDate(MAXIMUM_BASIC_PASSIVE_RATE_YEAR-1, 0, 0)
}
func (service *ServiceAPI) GetBasicPassiveRates(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetBasicPassiveRatesResponse {
	yearDifference := req.DateTo.Year() - req.DateFrom.Year()
	var basicPassiveRates []models.BasicPassiveRate

	dateFrom := req.DateFrom
	for {
		if yearDifference >= MAXIMUM_BASIC_PASSIVE_RATE_YEAR {
			newDateTo := service.add12years(dateFrom)
			result, err := service.Scrapper.GetBasicPassiveRateByDates(dateFrom, newDateTo)
			if err != nil {
				_ = level.Error(service.logger).Log("msg", "error scrapping basic passive rates by dates",
					"date_from", dateFrom, "date_to", newDateTo)
				break
			}
			basicPassiveRates = append(basicPassiveRates, result...)
			dateFrom = newDateTo.AddDate(0, 0, 1)
			yearDifference = req.DateTo.Year() - dateFrom.Year()

		} else {
			result, err := service.Scrapper.GetBasicPassiveRateByDates(dateFrom, req.DateTo)
			if err != nil {
				_ = level.Error(service.logger).Log("msg", "error scrapping basic passive rates by dates",
					"date_from", dateFrom, "date_to", req.DateTo)
				break
			}
			basicPassiveRates = append(basicPassiveRates, result...)
			break
		}
	}

	sort.Slice(basicPassiveRates, func(i, j int) bool {
		return basicPassiveRates[i].Date.Before(basicPassiveRates[j].Date)
	})
	return &GetBasicPassiveRatesResponse{
		BasicPassiveRates: basicPassiveRates,
		Err:               nil,
	}
}

func (service *ServiceAPI) GetTodayBasicPassiveRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayBasicPassiveRateResponse {
	date := time.Now()
	todayBasicPassiveRate, err := service.Scrapper.GetBasicPassiveDateByDate(date)
	if err != nil {
		_ = level.Error(service.logger).Log("msg", "error scrapping basic passive rate by date", "date", date,
			"error", err)
		return &GetTodayBasicPassiveRateResponse{
			BasicPassiveRate: nil,
			Err:              err,
		}
	}

	return &GetTodayBasicPassiveRateResponse{
		BasicPassiveRate: todayBasicPassiveRate,
		Err:              nil,
	}
}
