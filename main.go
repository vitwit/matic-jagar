package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/targets"
	"github.com/vitwit/matic-jagar/types"
)

func main() {
	cfg, err := config.ReadFromFile()
	if err != nil {
		log.Fatal(err)
	}

	m := targets.InitTargets(cfg)
	runner := targets.NewRunner()

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("%s:%s", cfg.InfluxDB.IP, cfg.InfluxDB.Port),
		Username: cfg.InfluxDB.Username,
		Password: cfg.InfluxDB.Password,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	var wg sync.WaitGroup
	for _, tg := range m.List {
		wg.Add(1)
		go func(target types.Target) {
			scrapeRate, err := time.ParseDuration(target.ScraperRate)
			if err != nil {
				log.Fatal(err)
			}
			for {
				runner.Run(target.Func, target.HTTPOptions, cfg, c)
				time.Sleep(scrapeRate)
			}
		}(tg)
	}
	wg.Wait()
}
