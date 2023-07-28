package services

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
)

type Endpoints struct {
	GetAllDolarColonesChanges endpoint.Endpoint
	GetTodayExchangeRate      endpoint.Endpoint
}

func MakeEndpoints(s *Service) Endpoints {
	return Endpoints{
		GetAllDolarColonesChanges: makeGetAllDolarColonesChangesEndpoint(s),
		GetTodayExchangeRate:      makeGetTodayExchangeRateEndpoint(s),
	}
}

func makeGetAllDolarColonesChangesEndpoint(s *Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		exchangesRates, err := s.GetDolarColonesChange(ctx)

		return GetAllDolarColonesChangesResponse{
			ExchangesRates: exchangesRates,
			Err:            err,
		}, err
	}
}

type GetAllDolarColonesChangesRequest struct {
}

type GetAllDolarColonesChangesResponse struct {
	ExchangesRates []models.ExchangeRate `json:"data"`
	Err            error                 `json:"error,omitempty"`
}

func (r GetAllDolarColonesChangesResponse) error() error { return r.Err }

type GetTodayExchangeRateRequest struct {
}

type GetTodayExchangeRateResponse struct {
	ExchangesRate models.ExchangeRate `json:"data"`
	Err           error               `json:"error,omitempty"`
}

func (r GetTodayExchangeRateResponse) error() error { return r.Err }

func makeGetTodayExchangeRateEndpoint(s *Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		exchangeRate, err := s.GetTodayExchangeRate(ctx)

		return GetTodayExchangeRateResponse{
			ExchangesRate: exchangeRate,
			Err:           err,
		}, err
	}
}
