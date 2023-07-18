package requester

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ansel1/merry"
	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/edwintrumpet/prueba-tecnica-boletia/internal/db"
	"github.com/sirupsen/logrus"
)

type response struct {
	Meta struct {
		LastUpdatedAt time.Time `json:"last_updated_at"`
	} `jsob:"meta"`
	Data map[string]struct {
		Code  string  `json:"code"`
		Value float64 `json:"value"`
	} `json:"data"`
}

func Start(ctx context.Context, currencies db.CurrencyRepo, requests db.RequestRepo) {
	c := &http.Client{
		Timeout: time.Second * time.Duration(config.Get.Timeout),
	}

	req, err := http.NewRequest(http.MethodGet, config.Get.CurrenciesHost, nil)
	if err != nil {
		err = merry.Wrap(err)
		logrus.WithFields(logrus.Fields{
			"stack": merry.Stacktrace(err),
		}).Fatal(err)
	}

	req.Header.Set("apikey", config.Get.ApiKey)

	ticker := time.NewTicker(time.Minute * time.Duration(config.Get.RequestsTime))
	defer ticker.Stop()

	for {
		go makeRequest(c, req, currencies, requests)

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func makeRequest(c *http.Client, req *http.Request, currencies db.CurrencyRepo, requests db.RequestRepo) {
	start := time.Now()

	httpRes, err := c.Do(req)
	if err != nil {
		statusCode := "500"
		if strings.Contains(err.Error(), "Timeout exceeded") {
			statusCode = "408"
		}

		request := db.Request{
			RequestedAt:     start.UTC(),
			RequestDuration: time.Since(start).Seconds(),
			ResponseStatus:  statusCode,
			ErrorMsg:        err.Error(),
			IsOK:            false,
		}
		saveRequest(request, requests)
		return
	}
	defer httpRes.Body.Close()

	responseTime := time.Since(start)

	statusCode := strconv.Itoa(httpRes.StatusCode)

	if statusCode != "200" {
		request := db.Request{
			RequestedAt:     start.UTC(),
			RequestDuration: responseTime.Seconds(),
			ResponseStatus:  statusCode,
			ErrorMsg:        "status not ok",
			IsOK:            false,
		}
		saveRequest(request, requests)
		return
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		request := db.Request{
			RequestedAt:     start.UTC(),
			RequestDuration: responseTime.Seconds(),
			ResponseStatus:  statusCode,
			ErrorMsg:        err.Error(),
			IsOK:            false,
		}
		saveRequest(request, requests)
		return
	}

	var res response
	if err = json.Unmarshal(body, &res); err != nil {
		request := db.Request{
			RequestedAt:     start.UTC(),
			RequestDuration: responseTime.Seconds(),
			ResponseStatus:  statusCode,
			ErrorMsg:        err.Error(),
			IsOK:            false,
		}
		saveRequest(request, requests)
		return
	}

	request := db.Request{
		RequestedAt:     start.UTC(),
		LastUpdatedAt:   res.Meta.LastUpdatedAt,
		RequestDuration: responseTime.Seconds(),
		ResponseStatus:  statusCode,
		IsOK:            true,
	}

	tx, err := requests.Begin()
	if err != nil {
		err = merry.Wrap(err)
		logrus.WithFields(logrus.Fields{
			"requestedAt":    request.RequestedAt,
			"responseStatus": request.ResponseStatus,
			"isOk":           request.IsOK,
			"errorMsg":       request.ErrorMsg,
			"stack":          merry.Stacktrace(err),
		}).Error(err)
		return
	}

	createdRequest, err := requests.CreateWithTx(request, tx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestedAt":    request.RequestedAt,
			"responseStatus": request.ResponseStatus,
			"isOk":           request.IsOK,
			"errorMsg":       request.ErrorMsg,
			"stack":          merry.Stacktrace(err),
		}).Error(err)
		return
	}

	if createdRequest == nil {
		err := merry.New("request not saved on db")
		logrus.WithFields(logrus.Fields{
			"requestedAt":    request.RequestedAt,
			"responseStatus": request.ResponseStatus,
			"isOk":           request.IsOK,
			"errorMsg":       request.ErrorMsg,
			"stack":          merry.Stacktrace(err),
		}).Error(err)
		return
	}

	listOfCurrencies := []db.SaveCurrency{}
	for _, val := range res.Data {
		currency := db.SaveCurrency{
			Code:      val.Code,
			Value:     val.Value,
			RequestID: createdRequest.ID,
		}

		listOfCurrencies = append(listOfCurrencies, currency)
	}

	ok, err := currencies.Create(listOfCurrencies, tx)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestId":      createdRequest.ID,
			"requestedAt":    request.RequestedAt,
			"responseStatus": request.ResponseStatus,
			"isOk":           request.IsOK,
			"errorMsg":       request.ErrorMsg,
			"stack":          merry.Stacktrace(err),
		}).Error(err)
		return
	}

	if !ok {
		err := merry.New("currencies not saved on db")
		logrus.WithFields(logrus.Fields{
			"requestId":      createdRequest.ID,
			"requestedAt":    request.RequestedAt,
			"responseStatus": request.ResponseStatus,
			"isOk":           request.IsOK,
			"errorMsg":       request.ErrorMsg,
			"stack":          merry.Stacktrace(err),
		}).Error(err)
		return
	}

	if err := tx.Commit(); err != nil {
		logrus.WithFields(logrus.Fields{
			"requestId":      createdRequest.ID,
			"requestedAt":    request.RequestedAt,
			"responseStatus": request.ResponseStatus,
			"isOk":           request.IsOK,
			"errorMsg":       request.ErrorMsg,
			"stack":          merry.Stacktrace(err),
		}).Error(err)
	}
}

func saveRequest(request db.Request, repo db.RequestRepo) string {
	createdRequest, err := repo.Create(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"requestedAt":    request.RequestedAt,
			"responseStatus": request.ResponseStatus,
			"isOk":           request.IsOK,
			"errorMsg":       request.ErrorMsg,
			"stack":          merry.Stacktrace(err),
		}).Error(err)
		return ""
	}

	if createdRequest == nil {
		err := merry.New("request not saved on db")
		logrus.WithFields(logrus.Fields{
			"requestedAt":    request.RequestedAt,
			"responseStatus": request.ResponseStatus,
			"isOk":           request.IsOK,
			"errorMsg":       request.ErrorMsg,
			"stack":          merry.Stacktrace(err),
		}).Error(err)
		return ""
	}

	return createdRequest.ID
}
