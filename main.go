package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/targets"
)

func main() {
	cfg, err := config.ReadFromFile()
	if err != nil {
		log.Fatal(err)
	}

	// str := "000000000000000000000000000000000000000000000030ca024f987b900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000043e400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e4b8e9222704401ad16d4d826732953daf07c7e200000000000000000000000015ed57ca28cbebb58d9c6c62f570046bc089bc660000000000000000000000000000000000000000000000000000000000000001"
	// a := targets.Hex2int(str)
	// a := web3.utils.hexToAscii(hex)

	// targets.CheckData()

	// log.Fatalf("len..", len("000000000000000000000000"))

	// targets.DecodeStringResp()

	m := targets.InitTargets(cfg)
	runner := targets.NewRunner()

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     fmt.Sprintf("http://localhost:%s", cfg.InfluxDB.Port),
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
		go func(target targets.Target) {
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
