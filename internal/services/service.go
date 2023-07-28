package services

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (service *Service) GetDolarColonesChange(ctx context.Context) ([]models.ExchangeRate, error) {
	url := "https://gee.bccr.fi.cr/indicadoreseconomicos/Cuadros/frmVerCatCuadro.aspx?idioma=1&CodCuadro=%20400"

	collyCollector := colly.NewCollector()

	exchangesRates := []models.ExchangeRate{}

	collyCollector.OnHTML("#theTable400", func(h *colly.HTMLElement) {
		dates := h.ChildTexts("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		buys := h.ChildTexts("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td")
		sales := h.ChildTexts("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(3) > table > tbody > tr > td > table > tbody > tr > td")

		if len(dates) != len(buys) || len(dates) != len(sales) {
			fmt.Println("Error: The number of dates, buys and sales are not the same")
			return
		}

		for i := 0; i < len(dates); i++ {
			salePrice := strings.ReplaceAll(sales[i], ",", ".")
			sale, err := strconv.ParseFloat(salePrice, 64)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			buyPrice := strings.ReplaceAll(buys[i], ",", ".")
			buy, err := strconv.ParseFloat(buyPrice, 64)
			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			exchangeRate := models.ExchangeRate{
				Date:      dates[i],
				SalePrice: sale,
				BuyPrice:  buy,
			}
			exchangesRates = append(exchangesRates, exchangeRate)
		}
	})

	collyCollector.Visit(url)

	if len(exchangesRates) > 0 {
		return exchangesRates, nil
	}

	return []models.ExchangeRate{}, nil
}

func (service *Service) GetTodayExchangeRate(ctx context.Context) (models.ExchangeRate, error) {
	url := "https://gee.bccr.fi.cr/indicadoreseconomicos/Cuadros/frmVerCatCuadro.aspx?idioma=1&CodCuadro=%20400"

	collyCollector := colly.NewCollector()

	todayExchangeRate := models.ExchangeRate{}

	collyCollector.OnHTML("#theTable400", func(h *colly.HTMLElement) {
		date := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr:nth-child(30) > td")
		buyHTML := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr:nth-child(30) > td")
		saleHTML := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(3) > table > tbody > tr > td > table > tbody > tr:nth-child(30) > td")

		salePrice := strings.ReplaceAll(saleHTML, ",", ".")
		sale, err := strconv.ParseFloat(salePrice, 64)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		buyPrice := strings.ReplaceAll(buyHTML, ",", ".")
		buy, err := strconv.ParseFloat(buyPrice, 64)
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		todayExchangeRate = models.ExchangeRate{
			Date:      date,
			SalePrice: sale,
			BuyPrice:  buy,
		}
	})

	collyCollector.Visit(url)

	return todayExchangeRate, nil
}
