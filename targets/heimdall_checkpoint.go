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

// GetTotalCheckPointsCount to get latest check points details
func GetLatestCheckpoints(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var lcp LatestCheckpoints
	err = json.Unmarshal(resp.Body, &lcp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	startBlock := lcp.Result.StartBlock
	endBlock := lcp.Result.EndBlock

	_ = writeToInfluxDb(c, bp, "matic_latest_checkpoint", map[string]string{}, map[string]interface{}{"start_block": startBlock, "end_block": endBlock})
	log.Printf("Latest checkpoint Start Block: %d and End Block: %d", startBlock, endBlock)
}

// GetCheckpointsDuration to get checkpoints duration
func GetCheckpointsDuration(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var cpd CheckpointsDuration
	err = json.Unmarshal(resp.Body, &cpd)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	duration := cpd.Result.CheckpointBufferTime
	minutes := ConvertNanoSecToMinutes(duration)

	_ = writeToInfluxDb(c, bp, "matic_checkpoint_duration", map[string]string{}, map[string]interface{}{"duration": minutes})
	log.Printf("Checkpoints Duration in nano seconds: %d", duration)
}
