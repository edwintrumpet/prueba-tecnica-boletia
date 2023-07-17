package requester

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
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

	Start(ctx, nil, nil)

	calls := httpmock.GetTotalCallCount()

	assert.Equal(t, 1, calls)
}
