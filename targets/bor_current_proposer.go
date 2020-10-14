package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetBorCurrentProposer to get current proposer and calculate no of blocks produced
func GetBorCurrentProposer(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var currentProposer EthResult
	err = json.Unmarshal(resp.Body, &currentProposer)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	prevCount := GetBlocksProducedCountFromDB(cfg, c)

	count, _ := strconv.Atoi(prevCount)

	proposer := currentProposer.Result

	if proposer == cfg.SignerAddress {
		count = count + 1
	}

	_ = writeToInfluxDb(c, bp, "matic_current_proposer", map[string]string{}, map[string]interface{}{"blocks_produced": count, "current_proposer": proposer})
	log.Printf("No of Blocks Proposed: %d", count)
}

// GetBlocksProducedCountFromDB returns the no of blocks produced from db
func GetBlocksProducedCountFromDB(cfg *config.Config, c client.Client) string {
	var count string
	q := client.NewQuery("SELECT last(blocks_produced) FROM matic_current_proposer", cfg.InfluxDB.Database, "")
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
