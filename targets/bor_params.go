package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetBorParams to get span duration and producer count
func GetBorParams(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var params BorParams
	err = json.Unmarshal(resp.Body, &params)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	spanDuration := params.Result.SpanDuration
	producerCount := params.Result.ProducerCount

	_ = writeToInfluxDb(c, bp, "matic_bor_params", map[string]string{}, map[string]interface{}{"span_duration": spanDuration, "producer_count": producerCount})
	log.Printf("Span Duration: %d and Producer Count: %d", spanDuration, producerCount)
}
