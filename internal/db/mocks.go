package db

import (
	"github.com/c2fo/testify/mock"
	"github.com/doug-martin/goqu/v9"
)

type MockRequestsRepo struct {
	mock.Mock
}

type MockCurrenciesRepo struct {
	mock.Mock
}

type MockTx struct {
	mock.Mock
}

func (m *MockRequestsRepo) Begin() (Tx, error) {
	args := m.Called()

	return args.Get(0).(Tx), args.Error(1)
}

func (m *MockRequestsRepo) Create(data Request) (*Request, error) {
	args := m.Called(data)

	return args.Get(0).(*Request), args.Error(1)
}

func (m *MockRequestsRepo) CreateWithTx(data Request, tx Tx) (*Request, error) {
	args := m.Called(data, tx)

	return args.Get(0).(*Request), args.Error(1)
}

func (m *MockCurrenciesRepo) Create(data []SaveCurrency, tx Tx) (bool, error) {
	args := m.Called(data, tx)

	return args.Bool(0), args.Error(1)
}

func (m *MockTx) Commit() error {
	args := m.Called()

	return args.Error(0)
}

func (m *MockTx) Insert(table interface{}) *goqu.InsertDataset {
	args := m.Called(table)
	return args.Get(0).(*goqu.InsertDataset)
}
