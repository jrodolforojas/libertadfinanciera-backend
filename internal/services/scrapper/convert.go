package scrapper

import (
	"strconv"
	"strings"
	"time"

	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
)

var months = map[string]int{
	"Enero":      int(time.January),
	"Febrero":    int(time.February),
	"Marzo":      int(time.March),
	"Abril":      int(time.April),
	"Mayo":       int(time.May),
	"Junio":      int(time.June),
	"Julio":      int(time.July),
	"Agosto":     int(time.August),
	"Septiembre": int(time.September),
	"Octubre":    int(time.October),
	"Noviembre":  int(time.November),
	"Diciembre":  int(time.December),
}

var prefixMonths = map[string]int{
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

func toCostaRicaInflationRate(costaRicaInflationRate models.CostaRicaInflationRateHTML) (models.CostaRicaInflationRate, error) {
	dateArray := strings.Split(costaRicaInflationRate.Date, "/") // 0: Month, 1: Year
	year, err := strconv.Atoi(dateArray[1])
	if err != nil {
		return models.CostaRicaInflationRate{}, err
	}
	month := months[dateArray[0]]
	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	valueHTML := strings.ReplaceAll(costaRicaInflationRate.Value, ",", ".")
	value, error := strconv.ParseFloat(valueHTML, 64)
	if error != nil {
		return models.CostaRicaInflationRate{}, error
	}
	inflationRate := models.CostaRicaInflationRate{
		Value: value,
		Date:  date,
	}

	return inflationRate, nil
}

func toTreasuryRateUSA(treasuryRateUSAHTML models.TreasuryRateUSAHTML) (models.TreasuryRateUSA, error) {
	dateArray := strings.Split(treasuryRateUSAHTML.Date, " ") // 0: Month, 1: Day, 2: Year
	year, err := strconv.Atoi(dateArray[2])
	if err != nil {
		return models.TreasuryRateUSA{}, err
	}
	month := prefixMonths[dateArray[1]]
	day, err := strconv.Atoi(dateArray[0])
	if err != nil {
		return models.TreasuryRateUSA{}, err
	}
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

	valueHTML := strings.ReplaceAll(treasuryRateUSAHTML.Value, ",", ".")
	value, error := strconv.ParseFloat(valueHTML, 64)
	if error != nil {
		return models.TreasuryRateUSA{}, error
	}
	treasuryRateUSA := models.TreasuryRateUSA{
		Value: value,
		Date:  date,
	}

	return treasuryRateUSA, nil
}

func toUSAInflationRate(inflationRateHTML models.USAInflationRateHTML) (models.USAInflationRate, error) {
	dateArray := strings.Split(inflationRateHTML.Date, " ") // 0: Year 1: Month
	year, err := strconv.Atoi(dateArray[0])
	if err != nil {
		return models.USAInflationRate{}, err
	}
	month, err := strconv.Atoi(dateArray[1])
	if err != nil {
		return models.USAInflationRate{}, err
	}
	date := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	lastDayOfMonth := date.AddDate(0, 1, -1)
	valueHTML := strings.ReplaceAll(inflationRateHTML.Value, ",", ".")
	value, error := strconv.ParseFloat(valueHTML, 64)
	if error != nil {
		return models.USAInflationRate{}, error
	}
	inflationRate := models.USAInflationRate{
		Value: value,
		Date:  lastDayOfMonth,
	}
	return inflationRate, nil
}
