package scrapper

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
)

type Scrapper interface {
	GetCurrentTableDollarColonesChange() ([]models.ExchangeRate, error)
	GetDollarColonesChangeByDates(dateFrom time.Time, dateTo time.Time) ([]models.ExchangeRate, error)
	GetLatestExchangeRate() (*models.ExchangeRate, error)
}

type BCCRScrapper struct {
	url string
}

func NewBCCRScrapper(url string) *BCCRScrapper {
	return &BCCRScrapper{
		url: url,
	}
}

func (scrapper *BCCRScrapper) GetCurrentTableDollarColonesChange() ([]models.ExchangeRate, error) {
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
			result := models.ExchangeRateHTML{
				SalePrice: sales[i],
				BuyPrice:  buys[i],
				Date:      dates[i],
			}

			toExchangeRate, err := result.ToExchangeRate()
			if err != nil {
				fmt.Println("error converting from html to exchange rate: ", err)
				return
			}
			exchangesRates = append(exchangesRates, toExchangeRate)
		}
	})

	collyCollector.Visit(scrapper.url)

	return exchangesRates, nil
}

func (scrapper *BCCRScrapper) GetDollarColonesChangeByDates(dateFrom time.Time, dateTo time.Time) ([]models.ExchangeRate, error) {
	return nil, nil
}

func (scrapper *BCCRScrapper) GetLatestExchangeRate() (*models.ExchangeRate, error) {
	collyCollector := colly.NewCollector()

	todayExchangeRate := models.ExchangeRate{}

	collyCollector.OnHTML("#theTable400", func(h *colly.HTMLElement) {
		date := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr:nth-child(30) > td")
		buyHTML := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr:nth-child(30) > td")
		saleHTML := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(3) > table > tbody > tr > td > table > tbody > tr:nth-child(30) > td")

		result := models.ExchangeRateHTML{
			SalePrice: saleHTML,
			BuyPrice:  buyHTML,
			Date:      date,
		}

		toExchangeRate, err := result.ToExchangeRate()
		if err != nil {
			fmt.Println("error converting from html to exchange rate: ", err)
			return
		}

		todayExchangeRate = toExchangeRate
	})

	collyCollector.Visit(scrapper.url)

	return &todayExchangeRate, nil
}
