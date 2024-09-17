package scrapper

import (
	"fmt"
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
	GetDollarColonesChangeByDates(dateFrom time.Time, dateTo time.Time, filtro int64) ([]models.ExchangeRate, error)
	GetExchangeRateByDate(date time.Time) (*models.ExchangeRate, error)
	GetBasicPassiveRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.BasicPassiveRate, error)
	GetBasicPassiveDateByDate(date time.Time) (*models.BasicPassiveRate, error)
	GetMonetaryPolicyRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.MonetaryPolicyRate, error)
	GetMonetaryPolicyRateByDate(date time.Time) (*models.MonetaryPolicyRate, error)
	GetPrimeRateByDates(dateFrom time.Time, dateTo time.Time) ([]models.PrimeRate, error)
	GetPrimeRateByDate(date time.Time) (*models.PrimeRate, error)
	GetCostaRicaInflationRateByDates(dateFrom time.Time, dateTo time.Time, filter int64) ([]models.CostaRicaInflationRate, error)
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

func (scrapper *BCCRScrapper) getScrappingUrlWithFilter(url string, dateFrom time.Time, dateTo time.Time, filter int64) string {
	return fmt.Sprintf(url, dateFrom.Format(utils.DATE_FORMAT), dateTo.Format(utils.DATE_FORMAT), fmt.Sprintf("%d", filter))
}
func (scrapper *BCCRScrapper) getScrappingUrl(url string, dateFrom time.Time, dateTo time.Time) string {
	return fmt.Sprintf(url, dateFrom.Format(utils.DATE_FORMAT), dateTo.Format(utils.DATE_FORMAT))
}

func (scrapper *BCCRScrapper) GetDollarColonesChangeByDates(dateFrom time.Time, dateTo time.Time, filtro int64) ([]models.ExchangeRate, error) {
	url := scrapper.getScrappingUrlWithFilter(scrapper.urls.ExchangeRateUrl, dateFrom, dateTo, filtro)
	fmt.Println(url)
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
	url := scrapper.getScrappingUrlWithFilter(scrapper.urls.ExchangeRateUrl, date, date, 0)
	fmt.Println(url)
	collyCollector := colly.NewCollector()

	todayExchangeRate := models.ExchangeRate{}

	collyCollector.OnHTML("#theTable400", func(h *colly.HTMLElement) {
		date := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		buyHTML := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(2) > table > tbody > tr > td > table > tbody > tr > td")
		saleHTML := h.ChildText("#theTable400 > tbody > tr:nth-child(2) > td:nth-child(3) > table > tbody > tr > td > table > tbody > tr > td")

		if buyHTML == "" || saleHTML == "" || date == "" {
			_ = level.Debug(scrapper.logger).Log("message", "error getting exchange rate from html", "url", url, "date", date)
			return
		}

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
	fmt.Println(url)
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
				if value != "" && date != "" {
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
	fmt.Println(url)
	collyCollector := colly.NewCollector()

	basicPassiveRate := models.BasicPassiveRate{}

	collyCollector.OnHTML("#Table17", func(h *colly.HTMLElement) {
		valueHTML := h.ChildText("#Table17 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(2) > p")
		dateHTML := h.ChildText("#Table17 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(1) > p")
		yearHTML := h.ChildText("#Table17 > tbody > tr:nth-child(1) > td:nth-child(2) > span > table > tbody > tr > td.celda17 > p")

		if valueHTML == "" || dateHTML == "" || yearHTML == "" {
			_ = level.Error(scrapper.logger).Log("msg", "error getting basic passive rate from html", "url", url, "date", date)
			return
		}

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
	fmt.Println(url)
	_ = level.Debug(scrapper.logger).Log("message", "url", url)
	yearFrom := dateFrom.Year()
	yearTo := dateTo.Year()

	yearDifference := (yearTo - yearFrom) + 2

	collyCollector := colly.NewCollector()

	monetaryPolicyRates := []models.MonetaryPolicyRate{}

	collyCollector.OnHTML("#Table779 > tbody", func(h *colly.HTMLElement) {
		column := h.ChildTexts("#Table779 > tbody > tr > td > span > table > tbody > tr > td")

		yearsHeader := column[1:yearDifference]
		years := []string{}

		for _, year := range yearsHeader {
			years = append(years, year[:4])
		}

		result := [][]string{}
		for i := yearDifference; i < len(column); i += (len(years) + 1) {
			row := column[i : i+len(years)+1]
			result = append(result, row)
		}

		var monetaryPolicyRatesHTML []models.MonetaryPolicyRateHTML
		for _, row := range result {
			values := row[1:] // <-- Get rates without first element (the date)
			for i := 0; i < len(values); i++ {
				date := row[0] + " " + years[i]
				value := values[i]
				if value != "" && date != "" {
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
				_ = level.Debug(scrapper.logger).Log("message", "error converting from html to monetary policy rate", "error", err)
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
	fmt.Println(url)
	collyCollector := colly.NewCollector()

	monetaryPolicyRate := models.MonetaryPolicyRate{}

	collyCollector.OnHTML("#Table779", func(h *colly.HTMLElement) {
		valueHTML := h.ChildText("#Table779 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(2) > p")
		dateHTML := h.ChildText("#Table779 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(1) > p")
		yearHTML := h.ChildText("#Table779 > tbody > tr:nth-child(1) > td:nth-child(2) > span > table > tbody > tr > td.celda779 > p")

		if valueHTML == "" || dateHTML == "" || yearHTML == "" {
			_ = level.Error(scrapper.logger).Log("msg", "error getting monetary policy rate from html", "url", url, "date", date)
			return
		}

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
	fmt.Println(url)
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
				if value != "" && date != "" {
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
	fmt.Println(url)
	collyCollector := colly.NewCollector()

	primeRate := models.PrimeRate{}

	collyCollector.OnHTML("#Table60", func(h *colly.HTMLElement) {
		valueHTML := h.ChildText("#Table60 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(2) > p")
		dateHTML := h.ChildText("#Table60 > tbody > tr:nth-child(2) > td:nth-child(2) > span > table > tbody > tr > td:nth-child(1) > p")
		yearHTML := h.ChildText("#Table60 > tbody > tr:nth-child(1) > td:nth-child(2) > span > table > tbody > tr > td.celda60 > p")

		if valueHTML == "" || dateHTML == "" || yearHTML == "" {
			_ = level.Error(scrapper.logger).Log("msg", "error getting prime rate from html", "url", url, "date", date)
			return
		}
		primeRateHTML := models.PrimeRateHTML{
			Value: valueHTML,
			Date:  dateHTML + " " + yearHTML,
		}

		toPrimeRate, err := primeRateHTML.ToPrimeRate()
		if err != nil {
			_ = level.Error(scrapper.logger).Log("msg", "error converting from MonetaryPolicyRateHTML to MonetaryPolicyRate models", "error", err)
			return
		}
		primeRate = toPrimeRate
		primeRate.Date = date
	})

	collyCollector.Visit(url)

	return &primeRate, nil
}

func (scrapper *BCCRScrapper) GetCostaRicaInflationRateByDates(dateFrom time.Time, dateTo time.Time, filter int64) ([]models.CostaRicaInflationRate, error) {
	url := scrapper.getScrappingUrlWithFilter(scrapper.urls.InflationCostaRicaUrl, dateFrom, dateTo, filter)
	fmt.Println(url)
	collyCollector := colly.NewCollector()

	inflationRates := []models.CostaRicaInflationRate{}

	collyCollector.OnHTML("#theTable2732 > tbody", func(h *colly.HTMLElement) {
		columns := h.ChildTexts("#theTable2732 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		interanualVariation := h.ChildTexts("#theTable2732 > tbody > tr:nth-child(2) > td:nth-child(4) > table > tbody > tr > td > table > tbody > tr > td")

		inflationRatesHTML := []models.CostaRicaInflationRateHTML{}
		for index, value := range interanualVariation {
			dateHTML := columns[index]
			if value != "" && dateHTML != "" {
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
				_ = level.Debug(scrapper.logger).Log("msg", "error converting from CostaRicaInflationRateHTML to CostaRicaInflationRate models", "error", err)
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
	url := scrapper.getScrappingUrlWithFilter(scrapper.urls.InflationCostaRicaUrl, dateFrom, dateTo, 0)
	fmt.Println(url)
	collyCollector := colly.NewCollector()

	inflationRate := models.CostaRicaInflationRate{}

	collyCollector.OnHTML("#theTable2732 > tbody", func(h *colly.HTMLElement) {
		dateHTML := h.ChildText("#theTable2732 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		valueHTML := h.ChildText("#theTable2732 > tbody > tr:nth-child(2) > td:nth-child(4) > table > tbody > tr > td > table > tbody > tr > td")

		if valueHTML == "" || dateHTML == "" {
			_ = level.Debug(scrapper.logger).Log("msg", "error getting inflation rate from html", "url", url, "date", date)
			return
		}
		inflationRateHTML := models.CostaRicaInflationRateHTML{
			Value: valueHTML,
			Date:  dateHTML,
		}

		costaRicaInflationRate, err := toCostaRicaInflationRate(inflationRateHTML)
		if err != nil {
			_ = level.Debug(scrapper.logger).Log("msg", "error converting from CostaRicaInflationRateHTML to CostaRicaInflationRate models", "error", err)
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
	collyCollector := colly.NewCollector()

	inflationRates := []models.USAInflationRate{}

	collyCollector.OnHTML("#table0", func(h *colly.HTMLElement) {
		headersHTML := h.ChildText("#table0 > thead > tr")
		yearsHTML := h.ChildTexts("#table0 > tbody > tr")
		valueHTML := h.ChildTexts("#table0 > tbody > tr > td")

		months := strings.Split(headersHTML, "  ")

		years := []string{}
		for _, year := range yearsHTML {
			years = append(years, year[:4])
		}

		indexYear := 0
		indexMonth := 0

		inflationRatesHTML := []models.USAInflationRateHTML{}
		for index, value := range valueHTML {
			if index%14 == 0 && index != 0 {
				indexYear++
			}
			if index%14 == 0 && index != 0 {
				indexMonth = 1
			} else {
				indexMonth++
			}

			if value != "" && months[indexMonth] != "" && years[indexYear] != "" && months[indexMonth] != "HALF1" && months[indexMonth] != "HALF2" {
				inflationRateHTML := models.USAInflationRateHTML{
					Value: value,
					Date:  years[indexYear] + " " + months[indexMonth],
				}

				fmt.Printf("%s-%s: %s\n", years[indexYear], months[indexMonth], value)
				inflationRatesHTML = append(inflationRatesHTML, inflationRateHTML)
			}
		}

		for _, inflationRateHTML := range inflationRatesHTML {
			inflationRate, err := toUSAInflationRate(inflationRateHTML)
			if err != nil {
				_ = level.Debug(scrapper.logger).Log("msg", "error converting from USAInflationRateHTML to USAInflationRate models", "error", err)
				return
			}

			// check if date is between dateFrom and dateTo
			if (inflationRate.Date.After(dateFrom) || inflationRate.Date.Equal(dateFrom)) && (inflationRate.Date.Before(dateTo) || inflationRate.Date.Equal(dateTo)) {
				inflationRates = append(inflationRates, inflationRate)
			}
		}
	})

	collyCollector.Visit(url)

	return inflationRates, nil
}

func (scrapper *BCCRScrapper) GetUSAInflationRateByDate(date time.Time) (*models.USAInflationRate, error) {
	url := scrapper.urls.InflationUSAUrl
	collyCollector := colly.NewCollector()

	inflationRate := models.USAInflationRate{}

	collyCollector.OnHTML("#table0", func(h *colly.HTMLElement) {
		headersHTML := h.ChildText("#table0 > thead > tr")
		yearsHTML := h.ChildTexts("#table0 > tbody > tr")
		valueHTML := h.ChildTexts("#table0 > tbody > tr > td")

		months := strings.Split(headersHTML, "  ")

		years := []string{}
		for _, year := range yearsHTML {
			years = append(years, year[:4])
		}

		indexYear := 0
		indexMonth := 0

		inflationRatesHTML := []models.USAInflationRateHTML{}
		for index, value := range valueHTML {
			if index%14 == 0 && index != 0 {
				indexYear++
			}
			if index%14 == 0 && index != 0 {
				indexMonth = 1
			} else {
				indexMonth++
			}

			if value != "" && months[indexMonth] != "" && years[indexYear] != "" && months[indexMonth] != "HALF1" && months[indexMonth] != "HALF2" {
				inflationRateHTML := models.USAInflationRateHTML{
					Value: value,
					Date:  years[indexYear] + " " + months[indexMonth],
				}

				inflationRatesHTML = append(inflationRatesHTML, inflationRateHTML)
			}
		}

		lastInflationRateHTML := inflationRatesHTML[len(inflationRatesHTML)-1]
		lastInflationRate, err := toUSAInflationRate(lastInflationRateHTML)
		if err != nil {
			_ = level.Debug(scrapper.logger).Log("msg", "error converting from USAInflationRateHTML to USAInflationRate models", "error", err)
			return
		}
		inflationRate = lastInflationRate
	})

	collyCollector.Visit(url)

	return &inflationRate, nil
}

func (scrapper *BCCRScrapper) GetTreasuryRateUSAByDates(dateFrom time.Time, dateTo time.Time) ([]models.TreasuryRateUSA, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.TreasuryRateUSAUrl, dateFrom, dateTo)
	fmt.Println(url)
	collyCollector := colly.NewCollector()

	treasuryRates := []models.TreasuryRateUSA{}

	collyCollector.OnHTML("#theTable677 > tbody", func(h *colly.HTMLElement) {
		columns := h.ChildTexts("#theTable677 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")
		rates := h.ChildTexts("#col_135401 > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr > td")

		treasuryRatesHTML := []models.TreasuryRateUSAHTML{}
		for index, value := range rates {
			dateHTML := columns[index]
			if value != "" && dateHTML != "" {
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
				_ = level.Error(scrapper.logger).Log("msg", "error converting from TreasuryRateUSAHTML to TreasuryRateUSA models", "error", err, "rate", treasuryRateHTML)
				return
			}
			treasuryRates = append(treasuryRates, treasuryRate)
		}
	})

	collyCollector.Visit(url)

	return treasuryRates, nil
}

func (scrapper *BCCRScrapper) GetTreasuryRateUSAByDate(date time.Time) (*models.TreasuryRateUSA, error) {
	url := scrapper.getScrappingUrl(scrapper.urls.TreasuryRateUSAUrl, date, date)
	fmt.Println(url)
	collyCollector := colly.NewCollector()

	treasuryRate := models.TreasuryRateUSA{}

	collyCollector.OnHTML("#theTable677 > tbody", func(h *colly.HTMLElement) {
		valueHTML := h.ChildText("#col_135401 > table > tbody > tr > td > table > tbody > tr > td > table > tbody > tr > td")
		dateHTML := h.ChildText("#theTable677 > tbody > tr:nth-child(2) > td:nth-child(1) > table > tbody > tr > td")

		if valueHTML == "" || dateHTML == "" {
			_ = level.Error(scrapper.logger).Log("msg", "error getting treasury rate from html", "url", url, "date", date)
			return
		}

		treasuryRateHTML := models.TreasuryRateUSAHTML{
			Value: valueHTML,
			Date:  dateHTML,
		}

		rate, err := toTreasuryRateUSA(treasuryRateHTML)
		if err != nil {
			_ = level.Error(scrapper.logger).Log("msg", "error converting from MonetaryPolicyRateHTML to MonetaryPolicyRate models",
				"url", url, "error", err, "date", date)
			return
		}
		treasuryRate = rate
		treasuryRate.Date = date
	})

	collyCollector.Visit(url)

	return &treasuryRate, nil
}
