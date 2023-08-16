package services

import "github.com/jrodolforojas/libertadfinanciera-backend/internal/models"

type GetAllDollarColonesChangesResponse struct {
	ExchangesRates []models.ExchangeRate `json:"data"`
	Err            error                 `json:"error,omitempty"`
}

func (r GetAllDollarColonesChangesResponse) error() error { return r.Err }

type GetTodayExchangeRateResponse struct {
	ExchangesRate *models.ExchangeRate `json:"data"`
	Err           error                `json:"error,omitempty"`
}

func (r GetTodayExchangeRateResponse) error() error { return r.Err }

type GetBasicPassiveRatesResponse struct {
	BasicPassiveRates []models.BasicPassiveRate `json:"data"`
	Err               error                     `json:"error,omitempty"`
}

func (r GetBasicPassiveRatesResponse) error() error { return r.Err }

type GetTodayBasicPassiveRateResponse struct {
	BasicPassiveRate *models.BasicPassiveRate `json:"data"`
	Err              error                    `json:"error,omitempty"`
}

func (r GetTodayBasicPassiveRateResponse) error() error { return r.Err }
