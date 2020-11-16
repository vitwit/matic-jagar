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

// GetNetworkLatestBlock to get latest block height of a network
func GetNetworkLatestBlock(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	networkBlock, err := scraper.GetStatus(ops)
	if err != nil {
		log.Printf("Error in network latest block: %v", err)
		return
	}

	if &networkBlock != nil {

		networkBlockHeight, err := strconv.Atoi(networkBlock.Result.SyncInfo.LatestBlockHeight)
		if err != nil {
			log.Println("Error while converting network height from string to int ", err)
		}
		_ = writeToInfluxDb(c, bp, "heimdall_network_latest_block", map[string]string{}, map[string]interface{}{"block_height": networkBlockHeight})
		log.Printf("Network height: %d", networkBlockHeight)

		// Calling function to get validator latest
		// block height
		validatorHeight := GetValidatorBlock(cfg, c)
		if validatorHeight == "" {
			log.Println("Error while fetching validator block height from db ", validatorHeight)
			return
		}

		vaidatorBlockHeight, _ := strconv.Atoi(validatorHeight)
		heightDiff := networkBlockHeight - vaidatorBlockHeight

		_ = writeToInfluxDb(c, bp, "heimdall_height_difference", map[string]string{}, map[string]interface{}{"difference": heightDiff})
		log.Printf("Network height: %d and Validator Height: %d", networkBlockHeight, vaidatorBlockHeight)

		// Send alert
		if int64(heightDiff) >= cfg.AlertingThresholds.BlockDiffThreshold && strings.ToUpper(cfg.AlerterPreferences.BlockDiffAlerts) == "YES" {
			_ = SendTelegramAlert(fmt.Sprintf("Heimdall Block Difference Alert: Block Difference between network and validator has exceeded %d", cfg.AlertingThresholds.BlockDiffThreshold), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Heimdall Block Difference Alert : Block difference between network and validator has exceeded %d", cfg.AlertingThresholds.BlockDiffThreshold), cfg)
			log.Println("Sent alert of block height difference")
		}
	}
}

// GetNetworkBlock returns network current block height
func GetNetworkBlock(cfg *config.Config, c client.Client) string {
	var networkHeight string
	q := client.NewQuery("SELECT last(block_height) FROM heimdall_network_latest_block", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						networkHeight = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}

	return networkHeight
}
