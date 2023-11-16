package configuration

import "github.com/caarlos0/env"

type ServerConfig struct {
	Address  AddressConfig
	Scrapper ScrapperConfig
	Database DatabaseConfig
}

type AddressConfig struct {
	Port           string   `env:"PORT"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS"`
}

type ScrapperConfig struct {
	ExchangeRateUrl       string `env:"EXCHANGE_RATE_URL"`
	BasicPassiveRateUrl   string `env:"TBP_URL"`
	MonetaryPolicyRateUrl string `env:"MONETARY_POLICY_RATE_URL"`
	PrimeRateUrl          string `env:"PRIME_RATE_URL"`
	InflationCostaRicaUrl string `env:"INFLATION_COSTA_RICA_URL"`
	InflationUSAUrl       string `env:"INFLATION_USA_URL"`
	TreasuryRateUSAUrl		string `env:"TREASURY_RATE_USA_URL"`
}

type DatabaseConfig struct {
	SupabaseUrl string `env:"SUPABASE_URL"`
	SupabaseKey string `env:"SUPABASE_KEY"`
}

func Read() (*ServerConfig, error) {
	config := ServerConfig{}
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	if err := env.Parse(&config.Address); err != nil {
		return nil, err
	}
	if err := env.Parse(&config.Scrapper); err != nil {
		return nil, err
	}
	if err := env.Parse(&config.Database); err != nil {
		return nil, err
	}
	return &config, nil
}
