package targets

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
	"github.com/vitwit/matic-jagar/utils"
)

// TotalCheckPointsCount is to get total no of check points and stores in db
func TotalCheckPointsCount(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	cp, err := scraper.GetTotalCheckPoints(ops)
	if err != nil {
		log.Printf("Error in get total checkPoints: %v", err)
		return
	}

	count := cp.Result.Result

	_ = db.WriteToInfluxDb(c, bp, "heimdall_total_checkpoints", map[string]string{}, map[string]interface{}{"total_count": count})
	log.Printf("Checkpoints total count: %d", count)
}

// LatestCheckpoints is to get latest check point and stores in db
func LatestCheckpoints(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	lcp, err := scraper.GetLatestCheckpoints(ops)
	if err != nil {
		log.Printf("Error while getting latest checkpoints : %v", err)
		return
	}

	startBlock := lcp.Result.StartBlock
	endBlock := lcp.Result.EndBlock

	_ = db.WriteToInfluxDb(c, bp, "heimdall_latest_checkpoint", map[string]string{}, map[string]interface{}{"start_block": startBlock, "end_block": endBlock})
	log.Printf("Latest checkpoint Start Block: %d and End Block: %d", startBlock, endBlock)
}

// CheckpointsDuration is to get checkpoints duration and stores in db
func CheckpointsDuration(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	cpd, err := scraper.GetCheckpointsDuration(ops)
	if err != nil {
		log.Printf("Error in get checkpoints duration: %v", err)
		return
	}

	if &cpd.Result == nil {
		log.Println("Got an empty response of checkpoints duration!")
		return
	}

	duration := cpd.Result.CheckpointBufferTime
	minutes := utils.ConvertNanoSecToMinutes(duration) //covert nano seconds to minutes

	_ = db.WriteToInfluxDb(c, bp, "heimdall_checkpoint_duration", map[string]string{}, map[string]interface{}{"duration": minutes})
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

// ProposedCheckpoints is to get proposed checkpoint, counts no of proposed checkpoints by validator and stores in db
func ProposedCheckpoints(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	// Get latest checkpoint from db
	latestCP := GetLatestCheckPoint(cfg, c)
	ops.Endpoint = ops.Endpoint + latestCP

	proposedCP, err := scraper.GetProposedCheckpoints(ops)
	if err != nil {
		log.Printf("Error in get proposed checkpoints: %v", err)
		return
	}

	if &proposedCP.Result == nil {
		log.Println("Got an empty response of proposed checkpoints!")
		return
	}

	// Get last proposed checkpoint from db
	lastProposedCheckpoint := GetLastProposedCheckpoint(cfg, c)

	if strings.EqualFold(proposedCP.Result.Proposer, cfg.ValDetails.SignerAddress) {
		num := GetProposedCount(cfg, c) // get checkpoints proposed count from db
		count, err := strconv.Atoi(num) // convert string to int
		if err != nil {
			log.Printf("Error while converting proposed checkpoints count to int : %v", err)
			return
		}
		if latestCP != lastProposedCheckpoint {
			count++
		}

		_ = db.WriteToInfluxDb(c, bp, "heimdall_proposed_checkpoint", map[string]string{}, map[string]interface{}{"last_proposed_cp": latestCP, "proposed_count": count})
		log.Printf("Latest Proposed Checkpoint : %s Proposed Count : %d", latestCP, count)
	}
}

// GetProposedCount returns the count of proposed checkpoints from db
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
