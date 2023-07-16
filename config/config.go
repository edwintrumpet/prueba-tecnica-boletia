package config

import "github.com/caarlos0/env/v9"

type config struct {
	RequestsTime   int    `env:"REQUESTS_TIME,required"`
	Timeout        int    `env:"TIMEOUT,required"`
	CurrenciesHost string `env:"CURRENCIES_HOST,required"`
	ApiKey         string `env:"API_KEY,required"`
}

var Get config

func New() error {
	if err := env.Parse(&Get); err != nil {
		return err
	}

	return nil
}

func NewMock(requestTime, timeout int, currenciesHost string) {
	Get = config{
		RequestsTime:   1,
		Timeout:        30,
		CurrenciesHost: currenciesHost,
		ApiKey:         "api-key",
	}
}
