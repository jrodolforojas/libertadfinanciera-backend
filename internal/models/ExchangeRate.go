package models

import "time"

type ExchangeRate struct {
	SalePrice float64   `json:"sale"`
	BuyPrice  float64   `json:"buy"`
	Date      time.Time `json:"date"`
}
