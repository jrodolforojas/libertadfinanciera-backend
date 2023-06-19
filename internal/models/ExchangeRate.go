package models

type ExchangeRate struct {
	SalePrice float64 `json:"sale"`
	BuyPrice  float64 `json:"buy"`
	Date      string  `json:"date"`
}
