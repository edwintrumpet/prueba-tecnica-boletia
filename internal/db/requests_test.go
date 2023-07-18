package db

import (
	"testing"
	"time"

	"github.com/doug-martin/goqu/v9"
	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/stretchr/testify/assert"
)

func TestNewRequestRepo(t *testing.T) {
	repo := NewRequestRepo(&goqu.Database{})
	assert.NotNil(t, repo)
}

func TestBegin(t *testing.T) {
	repo := databasePreparationForRequestsTest(t)

	tx, err := repo.Begin()
	assert.NoError(t, err)

	err = tx.Commit()
	assert.NoError(t, err)
}

func TestCreateRequest(t *testing.T) {
	repo := databasePreparationForRequestsTest(t)

	testCases := [...]struct {
		name  string
		data  Request
		error string
	}{
		{
			name: "success",
			data: Request{
				RequestedAt:     time.Now(),
				RequestDuration: 1,
				ResponseStatus:  "200",
				IsOK:            true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			created, err := repo.Create(tc.data)
			if err != nil {
				assert.NotEmpty(t, tc.error)
				assert.ErrorContains(t, err, tc.error)
				return
			}

			assert.Empty(t, tc.error)
			assert.NotNil(t, created)
		})
	}
}

func TestCreateRequestWithTx(t *testing.T) {
	repo := databasePreparationForRequestsTest(t)

	testCases := [...]struct {
		name  string
		data  Request
		error string
	}{
		{
			name: "success",
			data: Request{
				RequestedAt:     time.Now(),
				RequestDuration: 1,
				ResponseStatus:  "200",
				IsOK:            true,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tx, err := repo.Begin()
			assert.NoError(t, err)

			created, err := repo.CreateWithTx(tc.data, tx)
			if err != nil {
				assert.NotEmpty(t, tc.error)
				assert.ErrorContains(t, err, tc.error)
				return
			}

			assert.Empty(t, tc.error)

			err = tx.Commit()
			assert.NoError(t, err)

			assert.NotNil(t, created)
		})
	}
}

func databasePreparationForRequestsTest(t *testing.T) RequestRepo {
	err := config.NewMock(1, 1, "", "test-password", "localhost", "true")
	assert.NoError(t, err)

	db, err := NewMockDB()
	assert.NoError(t, err)

	return NewRequestRepo(db)
}
