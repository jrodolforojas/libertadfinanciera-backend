package services

import "time"

type GetAllDollarColonesChangesRequest struct {
	DateFrom time.Time `json:"date_from"`
	DateTo   time.Time `json:"date_to"`
}

type GetTodayExchangeRateRequest struct {
}

type GetDataByFilterRequest struct {
	Periodicity string `json:"periodicity"`
}
