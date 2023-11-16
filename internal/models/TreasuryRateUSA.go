package models

import "time"

type TreasuryRateUSAHTML struct {
	Date  string
	Value string
}

type TreasuryRateUSA struct {
	Value float64   `json:"value"`
	Date  time.Time `json:"date"`
}
