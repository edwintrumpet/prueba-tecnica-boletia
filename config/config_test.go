package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	os.Setenv("REQUESTS_TIME", "480")
	os.Setenv("TIMEOUT", "10")
	os.Setenv("CURRENCIES_HOST", "test-example.com")
	os.Setenv("API_KEY", "test-api-key")
	os.Setenv("DB_USER", "test-db-user")
	os.Setenv("DB_PASSWORD", "test-db-password")
	os.Setenv("DB_NAME", "test-db-name")

	err := New()
	assert.NoError(t, err)

	assert.Equal(t, Get.RequestsTime, 480)
	assert.Equal(t, Get.Timeout, 10)
	assert.Equal(t, Get.CurrenciesHost, "test-example.com")
	assert.Equal(t, Get.ApiKey, "test-api-key")
	assert.Equal(t, Get.DB.Host, "localhost")
	assert.Equal(t, Get.DB.Migrate, true)
	assert.Equal(t, Get.DB.Name, "test-db-name")
	assert.Equal(t, Get.DB.Password, "test-db-password")
	assert.Equal(t, Get.DB.Port, "5432")
	assert.Equal(t, Get.DB.SSLMode, "disable")
	assert.Equal(t, Get.DB.User, "test-db-user")

	// reset
	os.Unsetenv("REQUESTS_TIME")
	os.Unsetenv("TIMEOUT")
	os.Unsetenv("CURRENCIES_HOST")
	os.Unsetenv("API_KEY")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
}

func TestNewRequiredVariableIsNotSetError(t *testing.T) {
	err := New()
	assert.ErrorContains(t, err, "is not set")
}

func TestNewMock(t *testing.T) {
	NewMock(1, 30, "test-example.com")

	assert.Equal(t, Get.RequestsTime, 1)
	assert.Equal(t, Get.Timeout, 30)
	assert.Equal(t, Get.CurrenciesHost, "test-example.com")
	assert.Equal(t, Get.ApiKey, "api-key")
}
