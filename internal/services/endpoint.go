package services

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	GetAllDolarColonesChanges  endpoint.Endpoint
	GetExchangeRatesByFilter   endpoint.Endpoint
	GetTodayExchangeRate       endpoint.Endpoint
	GetBasicPassiveRates       endpoint.Endpoint
	GetTodayBasicPassiveRate   endpoint.Endpoint
	GetMonetaryPolicyRates     endpoint.Endpoint
	GetTodayMonetaryPolicyRate endpoint.Endpoint
	GetPrimeRates              endpoint.Endpoint
	GetPrimeRate               endpoint.Endpoint
	GetCostaRicaInflationRates endpoint.Endpoint
	GetCostaRicaInflationRate  endpoint.Endpoint
	GetTreasuryRatesUSA        endpoint.Endpoint
	GetTreasuryRateUSA         endpoint.Endpoint
	GetUSAInflationRates       endpoint.Endpoint
	GetUSAInflationRate        endpoint.Endpoint
}

func MakeEndpoints(s *ServiceAPI) Endpoints {
	return Endpoints{
		GetAllDolarColonesChanges:  makeGetAllDolarColonesChangesEndpoint(s),
		GetExchangeRatesByFilter:   makeGetExchangeRatesByFilterEndpoint(s),
		GetTodayExchangeRate:       makeGetTodayExchangeRateEndpoint(s),
		GetBasicPassiveRates:       makeGetBasicPassiveRatesEndpoint(s),
		GetTodayBasicPassiveRate:   makeGetTodayBasicPassiveRateEndpoint(s),
		GetMonetaryPolicyRates:     makeGetMonetaryPolicyRatesEndpoint(s),
		GetTodayMonetaryPolicyRate: makeGetTodayMonetaryPolicyRateEndpoint(s),
		GetPrimeRates:              makeGetPrimeRatesEndpoint(s),
		GetPrimeRate:               makeGetPrimeRateEndpoint(s),
		GetCostaRicaInflationRates: makeGetCostaRicaInflationRatesEndpoint(s),
		GetCostaRicaInflationRate:  makeGetCostaRicaInflationRateEndpoint(s),
		GetTreasuryRatesUSA:        makeGetTreasuryRatesUSAEndpoint(s),
		GetTreasuryRateUSA:         makeGetTreasuryRateUSAEndpoint(s),
		GetUSAInflationRates:       makeGetUSAInflationRatesEndpoint(s),
		GetUSAInflationRate:        makeGetUSAInflationRateEndpoint(s),
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

func makeGetExchangeRatesByFilterEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetDataByFilterRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetDataByFilterRequest")
		}

		result := s.GetExchangeRatesByFilter(ctx, req)
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

func makeGetBasicPassiveRatesEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetAllDollarColonesChangesRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetAllDollarColonesChangesRequest")
		}

		result := s.GetBasicPassiveRates(ctx, req)
		return result, nil
	}
}

func makeGetTodayBasicPassiveRateEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetTodayExchangeRateRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetTodayExchangeRateRequest")
		}
		result := s.GetTodayBasicPassiveRate(ctx, req)

		return result, nil
	}
}

func makeGetMonetaryPolicyRatesEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetAllDollarColonesChangesRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetAllDollarColonesChangesRequest")
		}

		result := s.GetMonetaryPolicyRates(ctx, req)
		return result, nil
	}
}

func makeGetTodayMonetaryPolicyRateEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetTodayExchangeRateRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetTodayExchangeRateRequest")
		}
		result := s.GetTodayMonetaryPolicyRate(ctx, req)

		return result, nil
	}
}

func makeGetPrimeRatesEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetAllDollarColonesChangesRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetAllDollarColonesChangesRequest")
		}

		result := s.GetPrimeRates(ctx, req)
		return result, nil
	}
}

func makeGetPrimeRateEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetTodayExchangeRateRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetTodayExchangeRateRequest")
		}
		result := s.GetTodayPrimeRate(ctx, req)

		return result, nil
	}
}

func makeGetCostaRicaInflationRatesEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetAllDollarColonesChangesRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetAllDollarColonesChangesRequest")
		}

		result := s.GetCostaRicaInflationRates(ctx, req)
		return result, nil
	}
}

func makeGetCostaRicaInflationRateEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetTodayExchangeRateRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetTodayExchangeRateRequest")
		}
		result := s.GetCostaRicaInflationRate(ctx, req)

		return result, nil
	}
}

func makeGetTreasuryRatesUSAEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetAllDollarColonesChangesRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetAllDollarColonesChangesRequest")
		}
		result := s.GetTreasuryRatesUSA(ctx, req)

		return result, nil
	}
}

func makeGetTreasuryRateUSAEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetTodayExchangeRateRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetTodayExchangeRateRequest")
		}
		result := s.GetTodayTreasuryRateUSA(ctx, req)

		return result, nil
	}
}

func makeGetUSAInflationRatesEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetAllDollarColonesChangesRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetAllDollarColonesChangesRequest")
		}
		result := s.GetUSAInflationRates(ctx, req)

		return result, nil
	}
}

func makeGetUSAInflationRateEndpoint(s *ServiceAPI) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(GetTodayExchangeRateRequest)
		if !ok {
			return nil, errors.New("unable to cast the request to a GetTodayExchangeRateRequest")
		}
		result := s.GetUSAInflationRate(ctx, req)

		return result, nil
	}
}
