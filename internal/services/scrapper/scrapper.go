package scrapper

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/log"
	"github.com/gocolly/colly/v2"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/configuration"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
	"github.com/jrodolforojas/libertadfinanciera-backend/internal/utils"
)

type Scrapper interface {
	GetDollarColonesChangeByDates(dateFrom time.Time, dateTo time.Time) ([]models.ExchangeRate, error)
	GetExchangeRateByDate(date time.Time) (*models.ExchangeRate, error)
	GetBasicPassiveRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.BasicPassiveRate, error)
	GetBasicPassiveDateByDate(date time.Time) (*models.BasicPassiveRate, error)
	GetMonetaryPolicyRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.MonetaryPolicyRate, error)
	GetMonetaryPolicyRateByDate(date time.Time) (*models.MonetaryPolicyRate, error)
	GetPrimeRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.PrimeRate, error)
	GetPrimeRateByDate(date time.Time) (*models.PrimeRate, error)
	GetCostaRicaInflationRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.CostaRicaInflationRate, error)
	GetCostaRicaInflationRateByDate(date time.Time) (*models.CostaRicaInflationRate, error)
	GetTreasuryRateUSAByDates(dateFrom time.Time, dateTo time.Time) ([]models.TreasuryRateUSA, error)
	GetTreasuryRateUSAByDate(date time.Time) (*models.TreasuryRateUSA, error)
	GetUSAInflationRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.USAInflationRate, error)
	GetUSAInflationRateByDate(date time.Time) (*models.USAInflationRate, error)
}

type BCCRScrapper struct {
	logger log.Logger
	urls   configuration.ScrapperConfig
}

func NewBCCRScrapper(logger log.Logger, urls configuration.ScrapperConfig) *BCCRScrapper {
	return &BCCRScrapper{
		logger: logger,
		urls:   urls,
	}
}

func (scrapper *BCCRScrapper) getScrappingUrl(url string, dateFrom time.Time, dateTo time.Time) string {
	return fmt.Sprintf(url, dateFrom.Format(utils.DATE_FORMAT), dateTo.Format(utils.DATE_FORMAT))
}

func (scrapper *BCCRScrapper) GetDollarColonesChangeByDates(dateFrom time.Time, dateTo time.Time) ([]models.ExchangeRate, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.ExchangeRateUrl, dateFrom, dateTo)
	collyCollector := colly.NewCollector()

	exchangesRates := []models.ExchangeRate{}

	collyCollector.OnHTML("#theTable400", func(h *colly.HTMLElement) {
		dates := h.ChildTexts("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		buys := h.ChildTexts("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td")
		sales := h.ChildTexts("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(3) > table > tbody > tr > td > table > tbody > tr > td")

		var exchangesRatesHTML []models.ExchangeRateHTML
		for index, date := range dates {
			buy := buys[index]
			sale := sales[index]
			if buy != "" && sale != "" {
				exchangeRateHTML := models.ExchangeRateHTML{
					SalePrice: sale,
					BuyPrice:  buy,
					Date:      date,
				}
				exchangesRatesHTML = append(exchangesRatesHTML, exchangeRateHTML)
			}
		}

		for _, exchangeRateHTML := range exchangesRatesHTML {
			toExchangeRate, err := exchangeRateHTML.ToExchangeRate()
			if err != nil {
				_ = level.Debug(scrapper.logger).Log("message", "error converting from html to exchange rate", "result", exchangesRatesHTML, "error", err)
				return
			}
			exchangesRates = append(exchangesRates, toExchangeRate)
		}
	})

	collyCollector.Visit(url)

	return exchangesRates, nil
}

func (scrapper *BCCRScrapper) GetExchangeRateByDate(date time.Time) (*models.ExchangeRate, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.ExchangeRateUrl, date, date)
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
	url := scrapper.getScrappingUrl(scrapper.urls.BasicPassiveRateUrl, dateFrom, dateTo)

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

func (scrapper *BCCRScrapper) GetBasicPassiveDateByDate(date time.Time) (*models.BasicPassiveRate, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.BasicPassiveRateUrl, date, date)

	collyCollector := colly.NewCollector()

	basicPassiveRate := models.BasicPassiveRate{}

	collyCollector.OnHTML("#Table17", func(h *colly.HTMLElement) {
		valueHTML := h.ChildText("#Table17 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(2) > p")
		dateHTML := h.ChildText("#Table17 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(1) > p")
		yearHTML := h.ChildText("#Table17 > tbody > tr:nth-child(1) > td:nth-child(2) > span > table > tbody > tr > td.celda17 > p")

		basicPassiveRateHTML := models.BasicPassiveRateHTML{
			Value: valueHTML,
			Date:  dateHTML + " " + yearHTML,
		}

		toBasicPassiveRate, err := basicPassiveRateHTML.ToBasicPassiveRate()
		if err != nil {
			_ = level.Debug(scrapper.logger).Log("msg", "error converting from BasicPassiveRateHTML to BasicPassiveRate models", "error", err)
			return
		}
		basicPassiveRate = toBasicPassiveRate
		basicPassiveRate.Date = date
	})

	collyCollector.Visit(url)

	return &basicPassiveRate, nil
}

func (scrapper *BCCRScrapper) GetMonetaryPolicyRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.MonetaryPolicyRate, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.MonetaryPolicyRateUrl, dateFrom, dateTo)

	yearFrom := dateFrom.Year()
	yearTo := dateTo.Year()

	yearDifference := (yearTo - yearFrom) + 2

	collyCollector := colly.NewCollector()

	monetaryPolicyRates := []models.MonetaryPolicyRate{}

	collyCollector.OnHTML("#Table779 > tbody", func(h *colly.HTMLElement) {
		column := h.ChildTexts("#Table779 > tbody > tr > td > span > table > tbody > tr > td")

		yearsHeader := column[1:yearDifference]

		result := [][]string{}
		for i := yearDifference; i < len(column); i += (len(yearsHeader) + 1) {
			row := column[i : i+len(yearsHeader)+1]
			result = append(result, row)
		}

		var monetaryPolicyRatesHTML []models.MonetaryPolicyRateHTML
		for _, row := range result {
			values := row[1:] // <-- Get rates without first element (the date)
			for i := 0; i < len(values); i++ {
				date := row[0] + " " + yearsHeader[i]
				value := values[i]
				if value != "" {
					monetaryPolicyRate := models.MonetaryPolicyRateHTML{
						Value: value,
						Date:  date,
					}
					monetaryPolicyRatesHTML = append(monetaryPolicyRatesHTML, monetaryPolicyRate)
				}
			}
		}
		for _, monetaryPolicyRateHTML := range monetaryPolicyRatesHTML {
			toMonetaryPolicyRate, err := monetaryPolicyRateHTML.ToMonetaryPolicyRate()
			if err != nil {
				_ = level.Debug(scrapper.logger).Log("message", "error converting from html to monetary policy rate", "result", result, "error", err)
				return
			}

			monetaryPolicyRates = append(monetaryPolicyRates, toMonetaryPolicyRate)
		}
	})

	collyCollector.Visit(url)

	return monetaryPolicyRates, nil
}

func (scrapper *BCCRScrapper) GetMonetaryPolicyRateByDate(date time.Time) (*models.MonetaryPolicyRate, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.MonetaryPolicyRateUrl, date, date)

	collyCollector := colly.NewCollector()

	monetaryPolicyRate := models.MonetaryPolicyRate{}

	collyCollector.OnHTML("#Table779", func(h *colly.HTMLElement) {
		valueHTML := h.ChildText("#Table779 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(2) > p")
		dateHTML := h.ChildText("#Table779 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(1) > p")
		yearHTML := h.ChildText("#Table779 > tbody > tr:nth-child(1) > td:nth-child(2) > span > table > tbody > tr > td.celda779 > p")

		monetaryPolicyRateHTML := models.MonetaryPolicyRateHTML{
			Value: valueHTML,
			Date:  dateHTML + " " + yearHTML,
		}

		toMonetaryPolicyRate, err := monetaryPolicyRateHTML.ToMonetaryPolicyRate()
		if err != nil {
			_ = level.Debug(scrapper.logger).Log("msg", "error converting from MonetaryPolicyRateHTML to MonetaryPolicyRate models", "error", err)
			return
		}
		monetaryPolicyRate = toMonetaryPolicyRate
		monetaryPolicyRate.Date = date
	})

	collyCollector.Visit(url)

	return &monetaryPolicyRate, nil
}

func (scrapper *BCCRScrapper) GetPrimeRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.PrimeRate, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.PrimeRateUrl, dateFrom, dateTo)

	yearFrom := dateFrom.Year()
	yearTo := dateTo.Year()

	yearDifference := (yearTo - yearFrom) + 2

	collyCollector := colly.NewCollector()

	primeRates := []models.PrimeRate{}

	collyCollector.OnHTML("#Table60 > tbody", func(h *colly.HTMLElement) {
		column := h.ChildTexts("#Table60 > tbody > tr > td > span > table > tbody > tr > td")

		yearsHeader := column[1:yearDifference]

		result := [][]string{}
		for i := yearDifference; i < len(column); i += (len(yearsHeader) + 1) {
			row := column[i : i+len(yearsHeader)+1]
			result = append(result, row)
		}

		var primeRatesHTML []models.PrimeRateHTML
		for _, row := range result {
			values := row[1:] // <-- Get rates without first element (the date)
			for i := 0; i < len(values); i++ {
				date := row[0] + " " + yearsHeader[i]
				value := values[i]
				if value != "" {
					primeRateHTML := models.PrimeRateHTML{
						Value: value,
						Date:  date,
					}
					primeRatesHTML = append(primeRatesHTML, primeRateHTML)
				}
			}
		}
		for _, primeRateHTML := range primeRatesHTML {
			toPrimeRate, err := primeRateHTML.ToPrimeRate()
			if err != nil {
				_ = level.Debug(scrapper.logger).Log("message", "error converting from html to monetary policy rate", "result", result, "error", err)
				return
			}

			primeRates = append(primeRates, toPrimeRate)
		}
	})

	collyCollector.Visit(url)

	return primeRates, nil
}

func (scrapper *BCCRScrapper) GetPrimeRateByDate(date time.Time) (*models.PrimeRate, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.PrimeRateUrl, date, date)

	collyCollector := colly.NewCollector()

	primeRate := models.PrimeRate{}

	collyCollector.OnHTML("#Table60", func(h *colly.HTMLElement) {
		valueHTML := h.ChildText("#Table60 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(2) > p")
		dateHTML := h.ChildText("#Table60 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(1) > p")
		yearHTML := h.ChildText("#Table60 > tbody > tr:nth-child(1) > td:nth-child(2) > span > table > tbody > tr > td.celda60 > p")

		primeRateHTML := models.PrimeRateHTML{
			Value: valueHTML,
			Date:  dateHTML + " " + yearHTML,
		}

		toPrimeRate, err := primeRateHTML.ToPrimeRate()
		if err != nil {
			_ = level.Debug(scrapper.logger).Log("msg", "error converting from MonetaryPolicyRateHTML to MonetaryPolicyRate models", "error", err)
			return
		}
		primeRate = toPrimeRate
		primeRate.Date = date
	})

	collyCollector.Visit(url)

	return &primeRate, nil
}

func (scrapper *BCCRScrapper) GetCostaRicaInflationRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.CostaRicaInflationRate, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.InflationCostaRicaUrl, dateFrom, dateTo)
	collyCollector := colly.NewCollector()

	inflationRates := []models.CostaRicaInflationRate{}

	collyCollector.OnHTML("#theTable2732 > tbody", func(h *colly.HTMLElement) {
		columns := h.ChildTexts("#theTable2732 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		interanualVariation := h.ChildTexts("#theTable2732 > tbody > tr:nth-child(2) > td:nth-child(4) > table > tbody > tr > td > table > tbody > tr > td")

		inflationRatesHTML := []models.CostaRicaInflationRateHTML{}
		for index, value := range interanualVariation {
			dateHTML := columns[index]
			if value != "" {
				inflationRateHTML := models.CostaRicaInflationRateHTML{
					Value: value,
					Date:  dateHTML,
				}
				inflationRatesHTML = append(inflationRatesHTML, inflationRateHTML)
			}
		}

		for _, inflationRateHTML := range inflationRatesHTML {
			inflationRate, err := toCostaRicaInflationRate(inflationRateHTML)
			if err != nil {
				_ = level.Error(scrapper.logger).Log("msg", "error converting from CostaRicaInflationRateHTML to CostaRicaInflationRate models", "error", err)
				return
			}
			inflationRates = append(inflationRates, inflationRate)
		}
	})

	collyCollector.Visit(url)

	return inflationRates, nil
}

func (scrapper *BCCRScrapper) GetCostaRicaInflationRateByDate(date time.Time) (*models.CostaRicaInflationRate, error) {
	dateFrom := time.Date(date.Year(), date.Month()-1, 1, 0, 0, 0, 0, time.UTC)
	dateTo := time.Date(date.Year(), date.Month()-1, 31, 0, 0, 0, 0, time.UTC)
	url := scrapper.getScrappingUrl(scrapper.urls.InflationCostaRicaUrl, dateFrom, dateTo)
	collyCollector := colly.NewCollector()

	inflationRate := models.CostaRicaInflationRate{}

	collyCollector.OnHTML("#theTable2732 > tbody", func(h *colly.HTMLElement) {
		dateHTML := h.ChildText("#theTable2732 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		valueHTML := h.ChildText("#theTable2732 > tbody > tr:nth-child(2) > td:nth-child(4) > table > tbody > tr > td > table > tbody > tr > td")

		inflationRateHTML := models.CostaRicaInflationRateHTML{
			Value: valueHTML,
			Date:  dateHTML,
		}

		costaRicaInflationRate, err := toCostaRicaInflationRate(inflationRateHTML)
		if err != nil {
			_ = level.Error(scrapper.logger).Log("msg", "error converting from CostaRicaInflationRateHTML to CostaRicaInflationRate models", "error", err)
			return
		}

		inflationRate = costaRicaInflationRate
		inflationRate.Date = dateTo

	})

	collyCollector.Visit(url)

	return &inflationRate, nil
}

func (scrapper *BCCRScrapper) GetUSAInflationRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.USAInflationRate, error) {
	url := scrapper.urls.InflationUSAUrl

	inflationRates := []models.USAInflationRate{}
	response, err := http.Get(url)
	if err != nil {
		_ = level.Error(scrapper.logger).Log("msg", "error getting response from url", "url", url, "error", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		_ = level.Error(scrapper.logger).Log("msg", "error reading response body", "error", err)
	}

	responseString := string(responseData)
	tableRows := strings.Split(responseString, "~")

	for _, row := range tableRows {
		fields := strings.Split(row, ",,")
		value := ""
		if len(fields) >= 3 {
			value = strings.ReplaceAll(fields[3], " ", "")
		}

		if value != "" {
			inflationRateHTML := models.USAInflationRateHTML{
				Value: value,
				Date:  fields[0] + " " + fields[1],
			}

			inflationRate, err := toUSAInflationRate(inflationRateHTML)
			if err != nil {
				_ = level.Error(scrapper.logger).Log("msg", "error converting from USAInflationRateHTML to USAInflationRate models", "error", err)
				return nil, err
			}

			// check if date is between dateFrom and dateTo
			if (inflationRate.Date.After(dateFrom) || inflationRate.Date.Equal(dateFrom)) && (inflationRate.Date.Before(dateTo) || inflationRate.Date.Equal(dateTo)) {
				inflationRates = append(inflationRates, inflationRate)
			}
		}
	}

	return inflationRates, nil
}

func (scrapper *BCCRScrapper) GetUSAInflationRateByDate(date time.Time) (*models.USAInflationRate, error) {
	url := scrapper.urls.InflationUSAUrl

	inflationRate := models.USAInflationRate{}
	response, err := http.Get(url)
	if err != nil {
		_ = level.Error(scrapper.logger).Log("msg", "error getting response from url", "url", url, "error", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		_ = level.Error(scrapper.logger).Log("msg", "error reading response body", "error", err)
	}

	responseString := string(responseData)
	tableRows := strings.Split(responseString, "~")

	row := tableRows[0]

	fields := strings.Split(row, ",,")
	value := ""
	if len(fields) >= 3 {
		value = strings.ReplaceAll(fields[3], " ", "")
	}

	if value != "" {
		inflationRateHTML := models.USAInflationRateHTML{
			Value: value,
			Date:  fields[0] + " " + fields[1],
		}

		todayInflationRate, err := toUSAInflationRate(inflationRateHTML)
		if err != nil {
			_ = level.Error(scrapper.logger).Log("msg", "error converting from USAInflationRateHTML to USAInflationRate models", "error", err)
			return nil, err
		}

		inflationRate = todayInflationRate
	}

	return &inflationRate, nil
}

func (scrapper *BCCRScrapper) GetTreasuryRateUSAByDates(dateFrom time.Time, dateTo time.Time) ([]models.TreasuryRateUSA, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.TreasuryRateUSAUrl, dateFrom, dateTo)
	collyCollector := colly.NewCollector()

	treasuryRates := []models.TreasuryRateUSA{}

	collyCollector.OnHTML("#theTable677 > tbody", func(h *colly.HTMLElement) {
		columns := h.ChildTexts("#theTable677 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		rates := h.ChildTexts("#col_135401 > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr > td")

		treasuryRatesHTML := []models.TreasuryRateUSAHTML{}
		for index, value := range rates {
			dateHTML := columns[index]
			if value != "" {
				treasuryRateHTML := models.TreasuryRateUSAHTML{
					Value: value,
					Date:  dateHTML,
				}
				treasuryRatesHTML = append(treasuryRatesHTML, treasuryRateHTML)
			}
		}

		for _, treasuryRateHTML := range treasuryRatesHTML {
			treasuryRate, err := toTreasuryRateUSA(treasuryRateHTML)
			if err != nil {
				_ = level.Error(scrapper.logger).Log("msg", "error converting from TreasuryRateUSAHTML to TreasuryRateUSA models", "error", err)
				return
			}
			treasuryRates = append(treasuryRates, treasuryRate)
		}
	})

	collyCollector.Visit(url)

	return treasuryRates, nil
}

func (scrapper *BCCRScrapper) GetTreasuryRateUSAByDate(date time.Time) (*models.TreasuryRateUSA, error) {
	yesterday := date.AddDate(0, 0, -1)
	url := scrapper.getScrappingUrl(scrapper.urls.TreasuryRateUSAUrl, yesterday, yesterday)

	collyCollector := colly.NewCollector()

	treasuryRate := models.TreasuryRateUSA{}

	collyCollector.OnHTML("#theTable677 > tbody", func(h *colly.HTMLElement) {
		valueHTML := h.ChildText("#col_135401 > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr > td")
		dateHTML := h.ChildText("#theTable677 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")

		treasuryRateHTML := models.TreasuryRateUSAHTML{
			Value: valueHTML,
			Date:  dateHTML,
		}

		rate, err := toTreasuryRateUSA(treasuryRateHTML)
		if err != nil {
			_ = level.Debug(scrapper.logger).Log("msg", "error converting from MonetaryPolicyRateHTML to MonetaryPolicyRate models", "error", err)
			return
		}
		treasuryRate = rate
		treasuryRate.Date = yesterday
	})

	collyCollector.Visit(url)

	return &treasuryRate, nil
}
