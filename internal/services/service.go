package services

import (
	"context"
	"time"

	"github.com/jrodolforojas/libertadfinanciera-backend/internal/repositories"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/services/scrapper"
)

const MINIMUM_DAYS_TO_GO_BACK = 30

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
	dateFrom, dateTo := getDatesFromToday(MINIMUM_DAYS_TO_GO_BACK)
	exchangesRates, err := service.Scrapper.GetDollarColonesChangeByDates(dateFrom, dateTo)

	for _, exchangeRate := range exchangesRates {
		_, err := service.Repository.SaveExchangeRate(exchangeRate)
		if err != nil {
			return &GetAllDollarColonesChangesResponse{
				ExchangesRates: nil,
				Err:            err,
			}
		}
	}
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

	result, err := service.Repository.SaveExchangeRate(*todayExchangeRate)
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

// Get date from and date to from today's date and the number of days to go back
func getDatesFromToday(days int) (time.Time, time.Time) {
	dateTo := time.Now()
	dateFrom := dateTo.AddDate(0, 0, -days)

	return dateFrom, dateTo
}
