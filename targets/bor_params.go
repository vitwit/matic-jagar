package targets

import (
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// BorParams is to get span duration, producer count and stores it in db
func BorParams(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	params, err := scraper.BorParams(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	spanDuration := params.Result.SpanDuration

	err = writeToInfluxDb(c, bp, "bor_params", map[string]string{}, map[string]interface{}{"span_duration": spanDuration})
	if err != nil {
		log.Printf("Error while storing span duration : %v", err)
	}
	log.Printf("Span Duration: %d ", spanDuration)
}
