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

type GetMonetaryPolicyRatesResponse struct {
	MonetaryPolicyRates []models.MonetaryPolicyRate `json:"data"`
	Err                 error                       `json:"error,omitempty"`
}

func (r GetMonetaryPolicyRatesResponse) error() error { return r.Err }

type GetTodayMonetaryPolicyRateResponse struct {
	MonetaryPolicyRate *models.MonetaryPolicyRate `json:"data"`
	Err                error                      `json:"error,omitempty"`
}

func (r GetTodayMonetaryPolicyRateResponse) error() error { return r.Err }

type GetPrimeRatesResponse struct {
	PrimeRates []models.PrimeRate `json:"data"`
	Err        error              `json:"error,omitempty"`
}

func (r GetPrimeRatesResponse) error() error { return r.Err }

type GetTodayPrimeRateResponse struct {
	PrimeRate *models.PrimeRate `json:"data"`
	Err       error             `json:"error,omitempty"`
}

func (r GetTodayPrimeRateResponse) error() error { return r.Err }
