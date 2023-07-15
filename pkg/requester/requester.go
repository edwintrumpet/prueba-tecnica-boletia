package requester

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/edwintrumpet/prueba-tecnica-boletia/config"
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

func Start() {
	c := &http.Client{
		Timeout: time.Second * time.Duration(config.Get.Timeout),
	}

	req, err := http.NewRequest(http.MethodGet, config.Get.CurrenciesHost, nil)
	if err != nil {
		// TODO handle error
		log.Fatal(err)
	}

	req.Header.Set("apikey", config.Get.ApiKey)

	for {
		go func() {
			start := time.Now()

			httpRes, err := c.Do(req)
			if err != nil {
				// TODO handle error, save in db
				// TODO handle timeout error
				log.Fatal(err)
			}
			defer httpRes.Body.Close()

			responseTime := time.Since(start)

			body, err := io.ReadAll(httpRes.Body)
			if err != nil {
				// TODO handle error, save in db
				log.Fatal(err)
			}

			var res response
			if err = json.Unmarshal(body, &res); err != nil {
				// TODO handle error, save in db
				log.Fatal(err)
			}

			// TODO parse into a struct
			// TODO save in db

			log.Println("response", res)
			log.Println("response time", responseTime)

		}()

		time.Sleep(time.Minute * time.Duration(config.Get.RequestsTime))
	}
}
