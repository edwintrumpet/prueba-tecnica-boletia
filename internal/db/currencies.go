package db

import (
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

type currencyRepo struct {
	db    *goqu.Database
	table exp.IdentifierExpression
}

type CurrencyRepo interface {
	Create(data []SaveCurrency, tx Tx) (bool, error)
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
