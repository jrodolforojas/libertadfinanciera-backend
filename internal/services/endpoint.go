package services

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
)

type Endpoints struct {
	GetAllDolarColonesChanges endpoint.Endpoint
}

func MakeEndpoints(s *Service) Endpoints {
	return Endpoints{
		GetAllDolarColonesChanges: makeGetAllDolarColonesChangesEndpoint(s),
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
