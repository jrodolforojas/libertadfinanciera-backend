package models

import "time"

type MonetaryPolicyRateHTML struct {
	Date  string
	Value string
}

type MonetaryPolicyRate struct {
	Value float64   `json:"value"`
	Date  time.Time `json:"date"`
}
