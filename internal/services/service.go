package services

import (
	"context"
	"time"

	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/repositories"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services/scrapper"
)

type Service interface {
	GetDollarColonesChange(ctx context.Context, req GetAllDollarColonesChangesRequest) *GetAllDollarColonesChangesResponse
	GetTodayExchangeRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayExchangeRateResponse
}

type ServiceAPI struct {
	Scrapper   scrapper.Scrapper
	Repository repositories.Repository
}

func NewService(scrapperService scrapper.Scrapper, repo repositories.Repository) *ServiceAPI {
	return &ServiceAPI{
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
			return &GetAllDollarColonesChangesResponse{
				ExchangesRates: nil,
				Err:            err,
			}
		}
		allExchangeRates = append(allExchangeRates, exchangeRates...)
		dateFrom = nextMonthDate
	}

	return &GetAllDollarColonesChangesResponse{
		ExchangesRates: allExchangeRates,
		Err:            nil,
	}
}

func (service *ServiceAPI) GetTodayExchangeRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayExchangeRateResponse {
	todayExchangeRate, err := service.Scrapper.GetExchangeRateByDate(time.Now())
	if err != nil {
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
