package config

import (
	"fmt"

	"github.com/ansel1/merry"
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
		return merry.Wrap(err)
	}

	return nil
}

func NewMock(requestTime, timeout int, currenciesHost, dbPassword, dbHost, dbMigrate string) error {
	err := env.ParseWithOptions(&Get, env.Options{
		Environment: map[string]string{
			"REQUESTS_TIME":   "1",
			"TIMEOUT":         "30",
			"CURRENCIES_HOST": currenciesHost,
			"API_KEY":         "api-key",
			"DB_USER":         "test-user",
			"DB_PASSWORD":     dbPassword,
			"DB_NAME":         "test_boletia_db",
			"DB_MIGRATE":      dbMigrate,
			"DB_HOST":         dbHost,
		},
	})

	return merry.Wrap(err)
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
