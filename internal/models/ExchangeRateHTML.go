package models

import (
	"strconv"
	"strings"
	"time"
)

var months = map[string]int{
	"Ene": int(time.January),
	"Feb": int(time.February),
	"Mar": int(time.March),
	"Abr": int(time.April),
	"May": int(time.May),
	"Jun": int(time.June),
	"Jul": int(time.July),
	"Ago": int(time.August),
	"Set": int(time.September),
	"Oct": int(time.October),
	"Nov": int(time.November),
	"Dic": int(time.December),
}

type ExchangeRateHTML struct {
	SalePrice string `json:"sale"`
	BuyPrice  string `json:"buy"`
	Date      string `json:"date"`
}

func (exchangeRateHTML ExchangeRateHTML) ToExchangeRate() (ExchangeRate, error) {
	salePrice := strings.ReplaceAll(exchangeRateHTML.SalePrice, ",", ".")
	sale, err := strconv.ParseFloat(salePrice, 64)
	if err != nil {
		return ExchangeRate{}, err
	}

	buyPrice := strings.ReplaceAll(exchangeRateHTML.BuyPrice, ",", ".")
	buy, err := strconv.ParseFloat(buyPrice, 64)
	if err != nil {
		return ExchangeRate{}, err
	}

	dateArray := strings.Split(exchangeRateHTML.Date, " ") // 0: Day, 1: Month, 2: Year
	year, err := strconv.Atoi(dateArray[2])
	if err != nil {
		return ExchangeRate{}, err
	}

	day, err := strconv.Atoi(dateArray[0])
	if err != nil {
		return ExchangeRate{}, err
	}

	month := months[dateArray[1]]

	date := time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)

	return ExchangeRate{
		SalePrice: sale,
		BuyPrice:  buy,
		Date:      date,
	}, nil
}

func (basicPassiveRateHTML BasicPassiveRateHTML) ToBasicPassiveRate() (BasicPassiveRate, error) {
	valueHTML := strings.ReplaceAll(basicPassiveRateHTML.Value, ",", ".")
	value, err := strconv.ParseFloat(valueHTML, 64)
	if err != nil && valueHTML != "" {
		return BasicPassiveRate{}, err
	}

	dateArray := strings.Split(basicPassiveRateHTML.Date, " ") // 0: Day, 1: Month, 2: Year
	year, err := strconv.Atoi(dateArray[2])
	if err != nil {
		return BasicPassiveRate{}, err
	}

	day, err := strconv.Atoi(dateArray[0])
	if err != nil {
		return BasicPassiveRate{}, err
	}

	month := months[dateArray[1]]

	date := time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)

	return BasicPassiveRate{
		Value: value,
		Date:  date,
	}, nil
}

func (monetaryPolicyRateHTML MonetaryPolicyRateHTML) ToMonetaryPolicyRate() (MonetaryPolicyRate, error) {
	valueHTML := strings.ReplaceAll(monetaryPolicyRateHTML.Value, ",", ".")
	value, err := strconv.ParseFloat(valueHTML, 64)
	if err != nil && valueHTML != "" {
		return MonetaryPolicyRate{}, err
	}

	dateArray := strings.Split(monetaryPolicyRateHTML.Date, " ") // 0: Day, 1: Month, 2: Year
	year, err := strconv.Atoi(dateArray[2])
	if err != nil {
		return MonetaryPolicyRate{}, err
	}

	day, err := strconv.Atoi(dateArray[0])
	if err != nil {
		return MonetaryPolicyRate{}, err
	}

	month := months[dateArray[1]]

	date := time.Date(year, time.Month(month), day, 12, 0, 0, 0, time.UTC)

	return MonetaryPolicyRate{
		Value: value,
		Date:  date,
	}, nil
}
