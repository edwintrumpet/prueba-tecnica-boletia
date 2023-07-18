package main

import (
	"context"
	"log"

	"github.com/edwintrumpet/prueba-tecnica-boletia/api"
	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/edwintrumpet/prueba-tecnica-boletia/internal/db"
	"github.com/edwintrumpet/prueba-tecnica-boletia/pkg/currencies"
	"github.com/edwintrumpet/prueba-tecnica-boletia/pkg/requester"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})

	if err := config.New(); err != nil {
		log.Fatal(err)
	}

	goquDB, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	currenciesRepo := db.NewCurrencyRepo(goquDB)
	requestsRepo := db.NewRequestRepo(goquDB)

	currenciesService := currencies.New(currenciesRepo)

	a := api.New(currenciesService)
	go a.Start()
	requester.Start(context.Background(), currenciesRepo, requestsRepo)
}
