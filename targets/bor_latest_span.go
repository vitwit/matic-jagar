package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetBorLatestSpan to get latest span id
func GetBorLatestSpan(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var latestSpan BorLatestSpan
	err = json.Unmarshal(resp.Body, &latestSpan)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	spanID := latestSpan.Result.SpanID

	_ = writeToInfluxDb(c, bp, "matic_bor_latest_span", map[string]string{}, map[string]interface{}{"span_id": spanID})
	log.Printf("Bor Latest Span ID: %d", spanID)
}
