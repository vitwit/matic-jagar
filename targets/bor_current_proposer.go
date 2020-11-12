package targets

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// GetBorCurrentProposer to get current proposer and calculate no of blocks produced
func GetBorCurrentProposer(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	currentProposer, err := scraper.EthResult(ops)
	if err != nil {
		log.Printf("Error in GetBorCurrentProposer: %v", err)
		return
	}

	prevCount := GetBlocksProducedCountFromDB(cfg, c)
	if prevCount == "" {
		prevCount = "0"
	}
	count, err := strconv.Atoi(prevCount)
	if err != nil {
		log.Printf("Error in conversion from string to int of produced count : %v", err)
		return
	}
	proposer := currentProposer.Result

	if strings.EqualFold(proposer, cfg.ValDetails.SignerAddress) {
		count = count + 1
	}

	_ = writeToInfluxDb(c, bp, "bor_current_proposer", map[string]string{}, map[string]interface{}{"blocks_produced": count, "current_proposer": proposer, "proposer": proposer[2:]})
	log.Printf("No of Blocks Proposed: %d and currnt Proposer : %s", count, proposer)
}

// GetBlocksProducedCountFromDB returns the no of blocks produced from db
func GetBlocksProducedCountFromDB(cfg *config.Config, c client.Client) string {
	var count string
	q := client.NewQuery("SELECT last(blocks_produced) FROM bor_current_proposer", cfg.InfluxDB.Database, "")
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
