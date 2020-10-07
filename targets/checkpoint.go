package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetTotalCheckPointsCount to get total no of check points
func GetTotalCheckPointsCount(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var cp TotalCheckpoints
	err = json.Unmarshal(resp.Body, &cp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	count := cp.Result.Result

	_ = writeToInfluxDb(c, bp, "matic_total_checkpoints", map[string]string{}, map[string]interface{}{"total_count": count})
	log.Printf("Checkpoints total count: %d", count)
}
