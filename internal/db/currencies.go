package db

import (
	"time"

	"github.com/ansel1/merry"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type Currency struct {
	ID        string  `json:"id" db:"id"`
	Code      string  `json:"code" db:"code"`
	Value     float64 `json:"value" db:"value"`
	RequestID string  `json:"requestId" db:"request_id"`
}

type SaveCurrency struct {
	Code      string  `json:"code" db:"code"`
	Value     float64 `json:"value" db:"value"`
	RequestID string  `json:"requestId" db:"request_id"`
}

type FindCurrenciesRequest struct {
	Code  string     `json:"code"`
	Finit *time.Time `json:"finit"`
	Fend  *time.Time `json:"fend"`
}

type FindCurrenciesResponse struct {
	Code            string    `json:"code" db:"code"`
	Value           float64   `json:"value" db:"value"`
	RequestedAt     time.Time `json:"requestedAt" db:"requested_at"`
	LastUpdatedAt   time.Time `json:"lastUpdatedAt" db:"last_updated_at"`
	RequestDuration float64   `json:"requestDuration" db:"request_duration"`
}

type currencyRepo struct {
	db    *goqu.Database
	table exp.IdentifierExpression
}

type CurrencyRepo interface {
	Create(data []SaveCurrency, tx Tx) (bool, error)
	Find(req FindCurrenciesRequest) ([]FindCurrenciesResponse, error)
}

func NewCurrencyRepo(db *goqu.Database) CurrencyRepo {
	return &currencyRepo{
		db:    db,
		table: goqu.T("currencies"),
	}
}

func (r *currencyRepo) Create(data []SaveCurrency, tx Tx) (bool, error) {
	res, err := tx.Insert(r.table).
		Rows(data).
		Returning("*").
		Executor().Exec()
	if err != nil {
		return false, merry.Wrap(err)
	}

	affectedRows, err := res.RowsAffected()
	if err != nil {
		return false, merry.Wrap(err)
	}

	if affectedRows < 1 {
		return false, nil
	}

	return true, nil
}

func (r *currencyRepo) Find(req FindCurrenciesRequest) ([]FindCurrenciesResponse, error) {
	res := []FindCurrenciesResponse{}

	q := r.db.From(r.table).
		Select(
			"currencies.code",
			"currencies.value",
			"requests.requested_at",
			"requests.last_updated_at",
			"requests.request_duration",
		).
		LeftJoin(goqu.T("requests"), goqu.On(goqu.Ex{"currencies.request_id": goqu.I("requests.id")}))

	if req.Code != "" {
		q = q.Where(goqu.C("code").Eq(req.Code))
	}

	if req.Finit != nil {
		q = q.Where(goqu.C("requested_at").Gt(*req.Finit))
	}

	if req.Fend != nil {
		q = q.Where(goqu.C("requested_at").Lt(*req.Fend))
	}

	err := q.Order(goqu.I("requests.requested_at").Desc(), goqu.I("currencies.code").Asc()).
		Executor().
		ScanStructs(&res)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	return res, nil
}
