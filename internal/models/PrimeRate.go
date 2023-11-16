package models

import "time"

type PrimeRateHTML struct {
	Date  string
	Value string
}

type PrimeRate struct {
	Value float64   `json:"value"`
	Date  time.Time `json:"date"`
}
