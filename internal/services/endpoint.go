package services

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetAllDolarColonesChanges endpoint.Endpoint
	GetTodayExchangeRate      endpoint.Endpoint
}

func MakeEndpoints(s *ServiceAPI) Endpoints {
	return Endpoints{
		GetAllDolarColonesChanges: makeGetAllDolarColonesChangesEndpoint(s),
		GetTodayExchangeRate:      makeGetTodayExchangeRateEndpoint(s),
	}
}

func makeGetAllDolarColonesChangesEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetAllDollarColonesChangesRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetAllDollarColonesChangesRequest")
		}

		result := s.GetDollarColonesChange(ctx, req)
		return result, nil
	}
}

func makeGetTodayExchangeRateEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetTodayExchangeRateRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetTodayExchangeRateRequest")
		}
		result := s.GetTodayExchangeRate(ctx, req)

		return result, nil
	}
}
