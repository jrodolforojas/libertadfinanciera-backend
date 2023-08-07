package services

import (
	"context"
	"time"

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
	exchangesRates, err := service.Scrapper.GetDollarColonesChangeByDates(req.DateFrom, req.DateTo)

	return &GetAllDollarColonesChangesResponse{
		ExchangesRates: exchangesRates,
		Err:            err,
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
