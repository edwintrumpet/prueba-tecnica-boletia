package main

import (
	"context"
	"log"

	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/edwintrumpet/prueba-tecnica-boletia/internal/db"
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

	currencies := db.NewCurrencyRepo(goquDB)
	requests := db.NewRequestRepo(goquDB)

	requester.Start(context.Background(), currencies, requests)
}
