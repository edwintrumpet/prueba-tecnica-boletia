package api

import (
	"net/http"

	"github.com/ansel1/merry"
	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/edwintrumpet/prueba-tecnica-boletia/internal/currencies"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type api struct {
	service currencies.Service
}

type API interface {
	Start()
}

func New(s currencies.Service) API {
	return &api{
		service: s,
	}
}

func (a *api) Start() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet},
	}))

	e.HTTPErrorHandler = httpErrorHandler

	e.GET("/currencies/:code", a.handler)

	e.Logger.Fatal(e.Start(config.Port()))
}

func (a *api) handler(c echo.Context) error {
	req := currencies.Request{
		Code:  c.Param("code"),
		Finit: c.QueryParam("finit"),
		Fend:  c.QueryParam("fend"),
	}

	res, err := a.service.Historial(req)
	if err != nil {
		return merry.Wrap(err)
	}

	return c.JSON(http.StatusOK, res)
}

func httpErrorHandler(err error, c echo.Context) {
	userMessage := merry.UserMessage(err)
	statusCode := merry.HTTPCode(err)

	if userMessage == "" {
		userMessage = "Internal server error"
	}

	logrus.WithFields(logrus.Fields{
		"stack":       merry.Stacktrace(err),
		"userMessage": userMessage,
		"statusCode":  statusCode,
		"error":       err.Error(),
	}).Info("api error")

	// nolint:all
	c.JSON(statusCode, ErrorResponse{Error: userMessage})
}
