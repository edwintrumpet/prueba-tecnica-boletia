package requester

import (
	"context"
	"encoding/json"
	"fmt"
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

func Start(ctx context.Context) {
	projectId := "4ugjhdsfb"

	fmt.Println(projectId)

	c := &http.Client{
		Timeout: time.Second * time.Duration(config.Get.Timeout),
	}

	req, err := http.NewRequest(http.MethodGet, config.Get.CurrenciesHost, nil)
	if err != nil {
		// TODO handle error
		log.Fatal(err)
	}

	req.Header.Set("apikey", config.Get.ApiKey)

	ticker := time.NewTicker(time.Minute * time.Duration(config.Get.RequestsTime))
	defer ticker.Stop()

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

			log.Println("response", res.Data["USD"])
			log.Println("response time", responseTime)

		}()

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
