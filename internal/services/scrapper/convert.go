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
