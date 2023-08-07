package services

import (
	"strconv"
	"strings"
	"time"

	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
)

type Month struct {
	Name   string
	Number int
}

var months = map[string]int{
	"Ene": 1,
	"Feb": 2,
	"Mar": 3,
	"Abr": 4,
	"May": 5,
	"Jun": 6,
	"Jul": 7,
	"Ago": 8,
	"Set": 9,
	"Oct": 10,
	"Nov": 11,
	"Dic": 12,
}

type ExchangeRateHTML struct {
	SalePrice string `json:"sale"`
	BuyPrice  string `json:"buy"`
	Date      string `json:"date"`
}

func (exchangeRateHTML ExchangeRateHTML) ToExchangeRate() (models.ExchangeRate, error) {
	salePrice := strings.ReplaceAll(exchangeRateHTML.SalePrice, ",", ".")
	sale, err := strconv.ParseFloat(salePrice, 64)
	if err != nil {
		return models.ExchangeRate{}, err
	}

	buyPrice := strings.ReplaceAll(exchangeRateHTML.BuyPrice, ",", ".")
	buy, err := strconv.ParseFloat(buyPrice, 64)
	if err != nil {
		return models.ExchangeRate{}, err
	}

	dateArray := strings.Split(exchangeRateHTML.Date, " ") // 0: Day, 1: Month, 2: Year
	year, err := strconv.Atoi(dateArray[2])
	if err != nil {
		return models.ExchangeRate{}, err
	}

	day, err := strconv.Atoi(dateArray[0])
	if err != nil {
		return models.ExchangeRate{}, err
	}

	month := months[dateArray[1]]

	date := time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)

	return models.ExchangeRate{
		SalePrice: sale,
		BuyPrice:  buy,
		Date:      date,
		CreatedAt: time.Now().UTC(),
	}, nil
}
