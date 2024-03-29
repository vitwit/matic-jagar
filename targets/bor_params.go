package targets

import (
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// BorParams is to get span duration, producer count and stores it in db
func BorParams(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error while creating batch points : %v", err)
		return
	}

	params, err := scraper.BorParams(ops)
	if err != nil {
		log.Printf("Error while getting bor params: %v", err)
		return
	}

	if &params.Result == nil {
		log.Printf("Got an empty response of bor params : %v", err)
		return
	}

	spanDuration := params.Result.SpanDuration

	err = db.WriteToInfluxDb(c, bp, "bor_params", map[string]string{}, map[string]interface{}{"span_duration": spanDuration})
	if err != nil {
		log.Printf("Error while storing span duration : %v", err)
	}
	log.Printf("Span Duration: %d ", spanDuration)
}
