package repositories

import (
	"time"

	"github.com/jrodolforojas/libertadfinanciera-backend/internal/models"
)

type Repository interface {
	SaveLatestExchangeRate(exchangeRate models.ExchangeRate) (*models.ExchangeRate, error)
	GetExchangeRates() ([]models.ExchangeRate, error)
	GetExchangeRateByDate(dateFrom time.Time, dateTo time.Time) (models.ExchangeRate, error)
}
