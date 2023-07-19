package currencies

import (
	"net/http"
	"strings"
	"time"

	"github.com/ansel1/merry"
	"github.com/edwintrumpet/prueba-tecnica-boletia/internal/db"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Request struct {
	Code  string `json:"code"`
	Finit string `json:"finit"`
	Fend  string `json:"fend"`
}

type service struct {
	repo db.CurrencyRepo
}

type Service interface {
	Historial(req Request) ([]db.FindCurrenciesResponse, error)
}

func New(repo db.CurrencyRepo) Service {
	return &service{
		repo: repo,
	}
}

func (req Request) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Code, validation.Required, is.Alpha, validation.Length(3, 5)),
		validation.Field(&req.Finit, validation.Date("2006-01-02T15:04:05")),
		validation.Field(&req.Fend, validation.Date("2006-01-02T15:04:05")),
	)
}

func (s *service) Historial(req Request) ([]db.FindCurrenciesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, merry.Wrap(err).
			WithHTTPCode(http.StatusBadRequest).
			WithUserMessage(err.Error())
	}

	data := db.FindCurrenciesRequest{}

	if strings.ToLower(req.Code) == "all" {
		data.Code = ""
	} else {
		data.Code = strings.ToUpper(req.Code)
	}

	if req.Finit != "" {
		finit, err := time.Parse("2006-01-02T15:04:05", req.Finit)
		if err != nil {
			return nil, merry.Wrap(err).
				WithHTTPCode(http.StatusBadRequest).
				WithUserMessage("wrong date format for finit")
		}
		data.Finit = &finit
	}

	if req.Fend != "" {
		fend, err := time.Parse("2006-01-02T15:04:05", req.Fend)
		if err != nil {
			return nil, merry.Wrap(err).
				WithHTTPCode(http.StatusBadRequest).
				WithUserMessage("wrong date format for fend")
		}
		data.Fend = &fend
	}

	res, err := s.repo.Find(data)
	if err != nil {
		return nil, merry.Wrap(err)
	}

	return res, nil
}
