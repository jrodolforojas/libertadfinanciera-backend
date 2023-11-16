package models

import "time"

type USAInflationRateHTML struct {
	Date  string
	Value string
}

type USAInflationRate struct {
	Value float64   `json:"value"`
	Date  time.Time `json:"date"`
}
