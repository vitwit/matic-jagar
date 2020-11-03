package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

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

	_ = writeToInfluxDb(c, bp, "heimdall_total_checkpoints", map[string]string{}, map[string]interface{}{"total_count": count})
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

	_ = writeToInfluxDb(c, bp, "heimdall_latest_checkpoint", map[string]string{}, map[string]interface{}{"start_block": startBlock, "end_block": endBlock})
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

	_ = writeToInfluxDb(c, bp, "heimdall_checkpoint_duration", map[string]string{}, map[string]interface{}{"duration": minutes})
	log.Printf("Checkpoints Duration in nano seconds: %d", duration)
}

// GetLatestCheckPoint returns the latest checkpoint from db
func GetLatestCheckPoint(cfg *config.Config, c client.Client) string {
	var count string
	q := client.NewQuery("SELECT last(total_count) FROM heimdall_total_checkpoints", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						value := r.Series[0].Values[0][idx]
						count = fmt.Sprintf("%v", value)
						break
					}
				}
			}
		}
	}
	return count
}

// GetProposedCheckpoints to get proposed checkpoint and count
func GetProposedCheckpoints(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	// Get latest checkpoint from db
	latestCP := GetLatestCheckPoint(cfg, c)

	ops.Endpoint = ops.Endpoint + latestCP

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var proposedCP ProposedCheckpoints
	err = json.Unmarshal(resp.Body, &proposedCP)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// Get last proposed checkpoint from db
	lastProposedCheckpoint := GetLastProposedCheckpoint(cfg, c)

	if strings.EqualFold(proposedCP.Result.Proposer, cfg.ValDetails.SignerAddress) {
		num := GetProposedCount(cfg, c)
		count, _ := strconv.Atoi(num)
		if latestCP != lastProposedCheckpoint {
			count++
		}

		_ = writeToInfluxDb(c, bp, "heimdall_proposed_checkpoint", map[string]string{}, map[string]interface{}{"last_proposed_cp": latestCP, "proposed_count": count})
		log.Fatalf("Latest Proposed Checkpoint : %s Proposed Count : %d", latestCP, count)
	}
}

// GetProposedCount returns the count of proposed checkpoints
func GetProposedCount(cfg *config.Config, c client.Client) string {
	var count string
	q := client.NewQuery("SELECT last(proposed_count) FROM heimdall_proposed_checkpoint", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						value := r.Series[0].Values[0][idx]
						count = fmt.Sprintf("%v", value)
						break
					}
				}
			}
		}
	}
	return count
}

// GetLastProposedCheckpoint returns the last proposed checkpoint from db
func GetLastProposedCheckpoint(cfg *config.Config, c client.Client) string {
	var cp string
	q := client.NewQuery("SELECT last(last_proposed_cp) FROM heimdall_proposed_checkpoint", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						value := r.Series[0].Values[0][idx]
						cp = fmt.Sprintf("%v", value)
						break
					}
				}
			}
		}
	}
	return cp
}
