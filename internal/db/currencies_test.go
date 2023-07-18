package db

import (
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/stretchr/testify/assert"
)

func TestNewCurrencyRepo(t *testing.T) {
	repo := NewCurrencyRepo(&goqu.Database{})
	assert.NotNil(t, repo)
}

func TestCreateCurrencies(t *testing.T) {
	currenciesRepo, tx, requestID, err := databasePreparationForCurrenciesTest(t)
	assert.NoError(t, err)

	testCases := [...]struct {
		name  string
		data  []SaveCurrency
		ok    bool
		error string
	}{
		{
			name: "success",
			data: []SaveCurrency{
				{
					Code:      "COP",
					Value:     4080.0095483066,
					RequestID: requestID,
				},
			},
			ok: true,
		},
		{
			name: "not existent request",
			data: []SaveCurrency{
				{
					Code:      "COP",
					Value:     4080.0095483066,
					RequestID: "e18dfbab-8cfc-4943-b450-940930db7c94",
				},
			},
			ok:    false,
			error: "violates foreign key constraint",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ok, err := currenciesRepo.Create(tc.data, tx)
			if err != nil {
				assert.NotEmpty(t, tc.error)
				assert.ErrorContains(t, err, tc.error)
				return
			}

			assert.Empty(t, tc.error)

			assert.Equal(t, tc.ok, ok)
		})
	}
}
func databasePreparationForCurrenciesTest(t *testing.T) (CurrencyRepo, Tx, string, error) {
	err := config.NewMock(1, 1, "", "test-password", "localhost", "true")
	assert.NoError(t, err)

	db, err := New()
	assert.NoError(t, err)

	currenciesRepo := NewCurrencyRepo(db)
	requestsRepo := NewRequestRepo(db)

	request, err := requestsRepo.Create(Request{
		RequestedAt:     time.Now(),
		RequestDuration: 1,
		ResponseStatus:  "200",
		IsOK:            true,
	})
	assert.NoError(t, err)

	tx, err := requestsRepo.Begin()
	assert.NoError(t, err)

	return currenciesRepo, tx, request.ID, nil
}
