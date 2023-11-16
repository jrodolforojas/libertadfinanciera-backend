package models

import "time"

type CostaRicaInflationRateHTML struct {
	Date  string
	Value string
}

type CostaRicaInflationRate struct {
	Value float64   `json:"value"`
	Date  time.Time `json:"date"`
}
