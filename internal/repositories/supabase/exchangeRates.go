package supabase

import (
	"time"

	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
)

const tableName = "exchange_rates"

func (supa *Supabase) SaveLatestExchangeRate(exchangeRate models.ExchangeRate) (*models.ExchangeRate, error) {

	var results []models.ExchangeRate
	err := supa.Client.DB.From(tableName).Insert(exchangeRate).Execute(&results)
	if err != nil {
		return nil, err
	}

	return &results[0], nil
}

func (supa *Supabase) GetExchangeRates() ([]models.ExchangeRate, error) {
	var result []models.ExchangeRate
	err := supa.Client.DB.From(tableName).Select("*").Execute(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (supa *Supabase) GetExchangeRateByDate(dateFrom time.Time, dateTo time.Time) (models.ExchangeRate, error) {
	result := models.ExchangeRate{}
	return result, nil
}
