package requester

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
)

func Start() {
	c := &http.Client{
		Timeout: time.Second * time.Duration(config.Get.Timeout),
	}

	for {
		go func() {
			start := time.Now()

			path := config.Get.CurrenciesHost
			res, err := c.Get(path)
			if err != nil {
				// TODO handle error, save in db
				// TODO handle timeout error
				log.Fatal(err)
			}
			defer res.Body.Close()

			responseTime := time.Since(start)

			body, err := io.ReadAll(res.Body)
			if err != nil {
				// TODO handle error, save in db
				log.Fatal(err)
			}

			// TODO parse into a struct
			// TODO save in db

			log.Println("body", string(body))
			log.Println("response time", responseTime)

		}()

		time.Sleep(time.Minute * time.Duration(config.Get.RequestsTime))
	}
}
