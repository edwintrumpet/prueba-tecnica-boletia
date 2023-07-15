package main

import (
	"log"

	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
	"github.com/edwintrumpet/prueba-tecnica-boletia/pkg/requester"
)

func main() {
	if err := config.New(); err != nil {
		log.Fatal(err)
	}

	requester.Start()
}
