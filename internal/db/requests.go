package db

import (
	"time"

	"github.com/ansel1/merry"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type Request struct {
	ID              string    `json:"id" db:"id"`
	RequestedAt     time.Time `json:"requestedAt" db:"requested_at"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt" db:"last_updated_at"`
	RequestDuration float64   `json:"requestDuration" db:"request_duration"`
	ResponseStatus  string    `json:"responseStatus" db:"response_status"`
	IsOK            bool      `json:"isOk" db:"is_ok"`
	ErrorMsg        string    `json:"errorMsg" db:"error_msg"`
}

type requestRepo struct {
	db    *goqu.Database
	table exp.IdentifierExpression
}

type RequestRepo interface {
	Begin() (*goqu.TxDatabase, error)
	Create(data Request) (*Request, error)
	CreateWithTx(data Request, tx *goqu.TxDatabase) (*Request, error)
}

func NewRequestRepo(db *goqu.Database) RequestRepo {
	return &requestRepo{
		db:    db,
		table: goqu.T("requests"),
	}
}

func (r *requestRepo) Begin() (*goqu.TxDatabase, error) {
	return r.db.Begin()
}

func (r *requestRepo) Create(data Request) (*Request, error) {
	created := new(Request)

	ok, err := r.db.Insert(r.table).Cols(
		"requested_at",
		"last_updated_at",
		"request_duration",
		"response_status",
		"is_ok",
		"error_msg",
	).Vals(goqu.Vals{
		data.RequestedAt,
		data.LastUpdatedAt,
		data.RequestDuration,
		data.ResponseStatus,
		data.IsOK,
		data.ErrorMsg,
	}).Returning("*").Executor().ScanStruct(created)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	if !ok {
		return nil, nil
	}

	return created, nil
}

func (r *requestRepo) CreateWithTx(data Request, tx *goqu.TxDatabase) (*Request, error) {
	created := new(Request)

	ok, err := tx.Insert(r.table).Cols(
		"requested_at",
		"last_updated_at",
		"request_duration",
		"response_status",
		"is_ok",
		"error_msg",
	).Vals(goqu.Vals{
		data.RequestedAt,
		data.LastUpdatedAt,
		data.RequestDuration,
		data.ResponseStatus,
		data.IsOK,
		data.ErrorMsg,
	}).Returning("*").Executor().ScanStruct(created)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	if !ok {
		return nil, nil
	}

	return created, nil
}
