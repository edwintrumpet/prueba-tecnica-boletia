package currencies

import (
	"errors"
	"testing"

	"github.com/edwintrumpet/prueba-tecnica-boletia/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNew(t *testing.T) {
	repo := New(&db.MockCurrenciesRepo{})
	assert.NotNil(t, repo)
}

func TestHistorialSuccessFindCop(t *testing.T) {
	repo := new(db.MockCurrenciesRepo)

	s := service{
		repo: repo,
	}

	repo.On("Find", mock.Anything).Return([]db.FindCurrenciesResponse{}, nil)

	_, err := s.Historial(Request{
		Code:  "COP",
		Finit: "2023-07-18T20:15:00",
		Fend:  "2023-07-18T20:16:00",
	})
	assert.NoError(t, err)

	repo.AssertNumberOfCalls(t, "Find", 1)
}

func TestHistorialSuccessFindAll(t *testing.T) {
	repo := new(db.MockCurrenciesRepo)

	s := service{
		repo: repo,
	}

	repo.On("Find", mock.Anything).Return([]db.FindCurrenciesResponse{}, nil)

	_, err := s.Historial(Request{
		Code:  "all",
		Finit: "2023-07-18T20:15:00",
		Fend:  "2023-07-18T20:16:00",
	})
	assert.NoError(t, err)

	repo.AssertNumberOfCalls(t, "Find", 1)
}

func TestHistorialErrors(t *testing.T) {
	repo := new(db.MockCurrenciesRepo)

	s := service{
		repo: repo,
	}

	testsCases := [...]struct {
		name    string
		request Request
		mocks   func()
		asserts func()
		err     string
	}{
		{
			name: "wrong code",
			request: Request{
				Code:  "CO",
				Finit: "2023-07-18T20:15:00",
				Fend:  "2023-07-18T20:16:00",
			},
			mocks: func() {},
			asserts: func() {
				repo.AssertNumberOfCalls(t, "Find", 0)
			},
			err: "code: the length must be between 3 and 5.",
		},
		{
			name: "wrong finit",
			request: Request{
				Code:  "COP",
				Finit: "2023-07-18T20:15:00Z",
				Fend:  "2023-07-18T20:16:00",
			},
			mocks: func() {},
			asserts: func() {
				repo.AssertNumberOfCalls(t, "Find", 0)
			},
			err: "finit: must be a valid date.",
		},
		{
			name: "wrong fend",
			request: Request{
				Code:  "COP",
				Finit: "2023-07-18T20:15:00",
				Fend:  "date",
			},
			mocks: func() {},
			asserts: func() {
				repo.AssertNumberOfCalls(t, "Find", 0)
			},
			err: "fend: must be a valid date.",
		},
		{
			name: "error repository",
			request: Request{
				Code:  "COP",
				Finit: "2023-07-18T20:15:00",
				Fend:  "2023-07-18T20:16:00",
			},
			mocks: func() {
				repo.On("Find", mock.Anything).Return(nil, errors.New("mock error"))
			},
			asserts: func() {
				repo.AssertNumberOfCalls(t, "Find", 0)
			},
			err: "mock error",
		},
	}

	for _, tc := range testsCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mocks()

			_, err := s.Historial(tc.request)
			if err != nil {
				assert.NotEmpty(t, tc.err)
				assert.ErrorContains(t, err, tc.err)
				return
			}

			assert.Empty(t, tc.err)

			tc.asserts()
		})
	}
}
