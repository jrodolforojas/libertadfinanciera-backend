package scrapper

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/gocolly/colly/v2"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/utils"
)

type Scrapper interface {
	GetDollarColonesChangeByDates(dateFrom time.Time, dateTo time.Time) ([]models.ExchangeRate, error)
	GetExchangeRateByDate(date time.Time) (*models.ExchangeRate, error)
	GetBasicPassiveRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.BasicPassiveRate, error)
}

type BCCRScrapper struct {
	logger              log.Logger
	url                 string
	basicPassiveRateUrl string
}

func NewBCCRScrapper(logger log.Logger, url string, basicPassiveRateUrl string) *BCCRScrapper {
	return &BCCRScrapper{
		logger:              logger,
		url:                 url,
		basicPassiveRateUrl: basicPassiveRateUrl,
	}
}

func (scrapper *BCCRScrapper) getScrappingUrl(url string, dateFrom time.Time, dateTo time.Time) string {
	return fmt.Sprintf(url, dateFrom.Format(utils.DATE_FORMAT), dateTo.Format(utils.DATE_FORMAT))
}

func (scrapper *BCCRScrapper) GetDollarColonesChangeByDates(dateFrom time.Time, dateTo time.Time) ([]models.ExchangeRate, error) {
	url := scrapper.getScrappingUrl(scrapper.url, dateFrom, dateTo)

	collyCollector := colly.NewCollector()

	exchangesRates := []models.ExchangeRate{}

	collyCollector.OnHTML("#theTable400", func(h *colly.HTMLElement) {
		dates := h.ChildTexts("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		buys := h.ChildTexts("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td")
		sales := h.ChildTexts("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(3) > table > tbody > tr > td > table > tbody > tr > td")

		if len(dates) != len(buys) || len(dates) != len(sales) {
			_ = level.Debug(scrapper.logger).Log("message", "The number of dates, buys and sales are not the same")
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
				_ = level.Debug(scrapper.logger).Log("message", "error converting from html to exchange rate", "result", result, "error", err)
				return
			}
			exchangesRates = append(exchangesRates, toExchangeRate)
		}
	})

	collyCollector.Visit(url)

	return exchangesRates, nil
}

func (scrapper *BCCRScrapper) GetExchangeRateByDate(date time.Time) (*models.ExchangeRate, error) {
	url := scrapper.getScrappingUrl(scrapper.url, date, date)
	collyCollector := colly.NewCollector()

	todayExchangeRate := models.ExchangeRate{}

	collyCollector.OnHTML("#theTable400", func(h *colly.HTMLElement) {
		date := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		buyHTML := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td")
		saleHTML := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(3) > table > tbody > tr > td > table > tbody > tr > td")

		result := models.ExchangeRateHTML{
			SalePrice: saleHTML,
			BuyPrice:  buyHTML,
			Date:      date,
		}

		toExchangeRate, err := result.ToExchangeRate()
		if err != nil {
			_ = level.Debug(scrapper.logger).Log("message", "error converting from html to exchange rate", "result", result, "error", err)
			return
		}

		todayExchangeRate = toExchangeRate
	})

	collyCollector.Visit(url)

	return &todayExchangeRate, nil
}

func (scrapper *BCCRScrapper) GetBasicPassiveRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.BasicPassiveRate, error) {
	url := scrapper.getScrappingUrl(scrapper.basicPassiveRateUrl, dateFrom, dateTo)

	yearFrom := dateFrom.Year()
	yearTo := dateTo.Year()

	yearDifference := (yearTo - yearFrom) + 2

	collyCollector := colly.NewCollector()

	basicPassiveRates := []models.BasicPassiveRate{}

	collyCollector.OnHTML("#Table17 > tbody", func(h *colly.HTMLElement) {
		column := h.ChildTexts("#Table17 > tbody > tr > td > span > table > tbody > tr > td")

		yearsHeader := column[1:yearDifference]

		result := [][]string{}
		for i := yearDifference; i < len(column); i += (len(yearsHeader) + 1) {
			row := column[i : i+len(yearsHeader)+1]
			result = append(result, row)
		}

		var basicPassiveRatesHTML []models.BasicPassiveRateHTML
		for _, row := range result {
			values := row[1:] // <-- Get rates without first element (the date)
			for i := 0; i < len(values); i++ {
				date := row[0] + " " + yearsHeader[i]
				value := values[i]
				if value != "" {
					basicPassiveRate := models.BasicPassiveRateHTML{
						Value: value,
						Date:  date,
					}
					basicPassiveRatesHTML = append(basicPassiveRatesHTML, basicPassiveRate)
				}
			}
		}

		for _, basicPassiveRateHTML := range basicPassiveRatesHTML {
			toBasicPassiveRate, err := basicPassiveRateHTML.ToBasicPassiveRate()
			if err != nil {
				_ = level.Debug(scrapper.logger).Log("message", "error converting from html to basic passive rate", "result", result, "error", err)
				return
			}

			basicPassiveRates = append(basicPassiveRates, toBasicPassiveRate)
		}
	})

	collyCollector.Visit(url)

	return basicPassiveRates, nil
}
