package services

import (
	"context"

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
	exchangesRates, err := service.Scrapper.GetCurrentTableDollarColonesChange()

	return &GetAllDollarColonesChangesResponse{
		ExchangesRates: exchangesRates,
		Err:            err,
	}
}

func (service *ServiceAPI) GetTodayExchangeRate(ctx context.Context, req GetTodayExchangeRateRequest) *GetTodayExchangeRateResponse {
	todayExchangeRate, err := service.Scrapper.GetLatestExchangeRate()
	if err != nil {
		return &GetTodayExchangeRateResponse{
			ExchangesRate: nil,
			Err:           err,
		}
	}

	// insert into database
	result, err := service.Repository.SaveLatestExchangeRate(*todayExchangeRate)
	if err != nil {
		return &GetTodayExchangeRateResponse{
			ExchangesRate: nil,
			Err:           err,
		}
	}

	return &GetTodayExchangeRateResponse{
		ExchangesRate: result,
		Err:           nil,
	}
}
