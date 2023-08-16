package models

import "time"

type BasicPassiveRateHTML struct {
	Value string
	Date  string
}

type BasicPassiveRate struct {
	Value float64   `json:"value"`
	Date  time.Time `json:"date"`
}
