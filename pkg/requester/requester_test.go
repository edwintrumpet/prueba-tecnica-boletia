package requester

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/edwintrumpet/prueba-tecnica-boletia/internal/db"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStart(t *testing.T) {
	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	body, err := os.ReadFile("./samples/response.json")
	assert.NoError(t, err)

	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			apiKey := r.Header.Get("apiKey")
			assert.Equal(t, config.Get.ApiKey, apiKey)
			return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
		},
	)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()

	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockTx.On("Commit").Return(nil)
	mockRequests.On("Begin").Return(mockTx, nil)
	mockRequests.On("CreateWithTx", mock.Anything, mock.Anything).Return(&db.Request{
		ID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
	}, nil)
	mockCurrencies.On("Create", []db.SaveCurrency{
		{
			Code:      "COP",
			Value:     4080.0095483066,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
		{
			Code:      "MXN",
			Value:     16.7513508952,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
		{
			Code:      "USD",
			Value:     1.000002931,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
	}, mock.Anything).Return(true, nil)

	Start(ctx, mockCurrencies, mockRequests)

	calls := httpmock.GetTotalCallCount()

	assert.Equal(t, 1, calls)
}

/*
Cosas que hay que probar
Se hace un request
Y guarda llama el repositorio

Tocaría mockear el repositorio
*/
