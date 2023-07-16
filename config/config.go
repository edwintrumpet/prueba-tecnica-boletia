package config

import (
	"fmt"

	"github.com/caarlos0/env/v9"
)

type config struct {
	RequestsTime   int    `env:"REQUESTS_TIME,required"`
	Timeout        int    `env:"TIMEOUT,required"`
	CurrenciesHost string `env:"CURRENCIES_HOST,required"`
	ApiKey         string `env:"API_KEY,required"`
	DB             struct {
		User     string `env:"DB_USER,required"`
		Password string `env:"DB_PASSWORD,required"`
		Name     string `env:"DB_NAME,required"`
		Host     string `env:"DB_HOST" envDefault:"localhost"`
		Port     string `env:"DB_PORT" envDefault:"5432"`
		SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
		Migrate  bool   `env:"DB_MIGRATE" envDefault:"true"`
	}
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
		DB: struct {
			User     string "env:\"DB_USER,required\""
			Password string "env:\"DB_PASSWORD,required\""
			Name     string "env:\"DB_NAME,required\""
			Host     string "env:\"DB_HOST\" envDefault:\"localhost\""
			Port     string "env:\"DB_PORT\" envDefault:\"5432\""
			SSLMode  string "env:\"DB_SSL_MODE\" envDefault:\"disable\""
			Migrate  bool   "env:\"DB_MIGRATE\" envDefault:\"true\""
		}{
			User:     "test-user",
			Password: "test-password",
			Name:     "test_boletia_db",
			Host:     "localhost",
			Port:     "5433",
			SSLMode:  "disable",
			Migrate:  true,
		},
	}
}

func DBdsn() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=%s",
		Get.DB.User,
		Get.DB.Password,
		Get.DB.Name,
		Get.DB.Host,
		Get.DB.Port,
		Get.DB.SSLMode,
	)
}
