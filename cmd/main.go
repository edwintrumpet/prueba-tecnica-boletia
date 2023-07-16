package main

import (
	"context"
	"log"

	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/edwintrumpet/prueba-tecnica-boletia/internal/db"
	"github.com/edwintrumpet/prueba-tecnica-boletia/pkg/requester"
)

func main() {
	if err := config.New(); err != nil {
		log.Fatal(err)
	}

	if err := db.New(); err != nil {
		log.Fatal(err)
	}

	requester.Start(context.Background())
}
