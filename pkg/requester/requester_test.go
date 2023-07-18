package requester

import (
	"context"
	"errors"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/c2fo/testify/mock"
	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/edwintrumpet/prueba-tecnica-boletia/internal/db"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestStartSuccess(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	body, err := os.ReadFile("./samples/response.json")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			apiKey := r.Header.Get("apiKey")
			assert.Equal(t, config.Get.ApiKey, apiKey)
			return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Begin").Return(mockTx, nil)
	mockTx.On("Commit").Return(nil)
	mockRequests.On("CreateWithTx", mock.Anything, mock.Anything).Return(&db.Request{
		ID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
	}, nil)
	mockCurrencies.On("Create", []db.SaveCurrency{
		{
			Code:      "COP",
			Value:     4080.0095483066,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
	}, mock.Anything).Return(true, nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 1)
	mockRequests.AssertNumberOfCalls(t, "Begin", 1)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 1)
	mockRequests.AssertNumberOfCalls(t, "Create", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 1)
	mockCurrencies.AssertCalled(t, "Create", []db.SaveCurrency{
		{
			Code:      "COP",
			Value:     4080.0095483066,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
	}, mock.Anything)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartTimeoutExceeded(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return nil, errors.New("Timeout exceeded")
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Create", mock.Anything).Return(&db.Request{}, nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 0)
	mockRequests.AssertNumberOfCalls(t, "Create", 1)
	mockRequests.AssertNumberOfCalls(t, "Begin", 0)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 0)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartAPIUnauthorized(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusUnauthorized, `{"msg":"Unauthorized"}`), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Create", mock.Anything).Return(&db.Request{}, nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 0)
	mockRequests.AssertNumberOfCalls(t, "Create", 1)
	mockRequests.AssertNumberOfCalls(t, "Begin", 0)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 0)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartBadResponse(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, `{{}`), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Create", mock.Anything).Return(&db.Request{}, nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 0)
	mockRequests.AssertNumberOfCalls(t, "Create", 1)
	mockRequests.AssertNumberOfCalls(t, "Begin", 0)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 0)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartErrorOnInitTx(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	body, err := os.ReadFile("./samples/response.json")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Begin").Return(mockTx, errors.New("test error"))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockRequests.AssertNumberOfCalls(t, "Begin", 1)
	mockRequests.AssertNumberOfCalls(t, "Create", 0)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 0)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartErrorCreatingRequestsWithTx(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	body, err := os.ReadFile("./samples/response.json")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Begin").Return(mockTx, nil)
	mockRequests.On("CreateWithTx", mock.Anything, mock.Anything).Return(&db.Request{
		ID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
	}, errors.New("test error"))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 0)
	mockRequests.AssertNumberOfCalls(t, "Begin", 1)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 1)
	mockRequests.AssertNumberOfCalls(t, "Create", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 0)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartRequestWithTxNotSaved(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	body, err := os.ReadFile("./samples/response.json")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Begin").Return(mockTx, nil)
	mockRequests.On("CreateWithTx", mock.Anything, mock.Anything).Return(nil, nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 0)
	mockRequests.AssertNumberOfCalls(t, "Begin", 1)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 1)
	mockRequests.AssertNumberOfCalls(t, "Create", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 0)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartErrorCreatingCurrency(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	body, err := os.ReadFile("./samples/response.json")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
		},
	)
	defer httpmock.DeactivateAndReset()

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
	}, mock.Anything).Return(false, errors.New("test error"))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 0)
	mockRequests.AssertNumberOfCalls(t, "Begin", 1)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 1)
	mockRequests.AssertNumberOfCalls(t, "Create", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 1)
	mockCurrencies.AssertCalled(t, "Create", []db.SaveCurrency{
		{
			Code:      "COP",
			Value:     4080.0095483066,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
	}, mock.Anything)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartCurrenciesNotSaved(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	body, err := os.ReadFile("./samples/response.json")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			apiKey := r.Header.Get("apiKey")
			assert.Equal(t, config.Get.ApiKey, apiKey)
			return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Begin").Return(mockTx, nil)
	mockTx.On("Commit").Return(nil)
	mockRequests.On("CreateWithTx", mock.Anything, mock.Anything).Return(&db.Request{
		ID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
	}, nil)
	mockCurrencies.On("Create", []db.SaveCurrency{
		{
			Code:      "COP",
			Value:     4080.0095483066,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
	}, mock.Anything).Return(false, nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 0)
	mockRequests.AssertNumberOfCalls(t, "Begin", 1)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 1)
	mockRequests.AssertNumberOfCalls(t, "Create", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 1)
	mockCurrencies.AssertCalled(t, "Create", []db.SaveCurrency{
		{
			Code:      "COP",
			Value:     4080.0095483066,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
	}, mock.Anything)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartErrorOnCommit(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	body, err := os.ReadFile("./samples/response.json")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, string(body)), nil
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Begin").Return(mockTx, nil)
	mockTx.On("Commit").Return(errors.New("test-error"))
	mockRequests.On("CreateWithTx", mock.Anything, mock.Anything).Return(&db.Request{
		ID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
	}, nil)
	mockCurrencies.On("Create", []db.SaveCurrency{
		{
			Code:      "COP",
			Value:     4080.0095483066,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
	}, mock.Anything).Return(true, nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 1)
	mockRequests.AssertNumberOfCalls(t, "Begin", 1)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 1)
	mockRequests.AssertNumberOfCalls(t, "Create", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 1)
	mockCurrencies.AssertCalled(t, "Create", []db.SaveCurrency{
		{
			Code:      "COP",
			Value:     4080.0095483066,
			RequestID: "ec4b371a-379d-4a8d-bfdd-2c933fce0267",
		},
	}, mock.Anything)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartErrorCreatingRequest(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return nil, errors.New("Internal test server error")
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Create", mock.Anything).Return(nil, errors.New("test-error"))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 0)
	mockRequests.AssertNumberOfCalls(t, "Create", 1)
	mockRequests.AssertNumberOfCalls(t, "Begin", 0)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 0)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}

func TestStartRequestNotSaved(t *testing.T) {
	mockCurrencies := new(db.MockCurrenciesRepo)
	mockRequests := new(db.MockRequestsRepo)
	mockTx := new(db.MockTx)

	mockUrl := "example.com"
	err := config.NewMock(1, 10, mockUrl, "", "", "")
	assert.NoError(t, err)

	httpmock.Activate()
	httpmock.RegisterResponder(
		http.MethodGet,
		mockUrl,
		func(r *http.Request) (*http.Response, error) {
			return nil, errors.New("Internal test server error")
		},
	)
	defer httpmock.DeactivateAndReset()

	mockRequests.On("Create", mock.Anything).Return(nil, nil)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()
	Start(ctx, mockCurrencies, mockRequests)

	mockTx.AssertNumberOfCalls(t, "Commit", 0)
	mockRequests.AssertNumberOfCalls(t, "Create", 1)
	mockRequests.AssertNumberOfCalls(t, "Begin", 0)
	mockRequests.AssertNumberOfCalls(t, "CreateWithTx", 0)
	mockCurrencies.AssertNumberOfCalls(t, "Create", 0)

	calls := httpmock.GetTotalCallCount()
	assert.Equal(t, 1, calls)
}
