package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetLatestProposedBlockAndTime to get latest proposed block height and time
func GetLatestProposedBlockAndTime(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var blockResp LastProposedBlockAndTime
	err = json.Unmarshal(resp.Body, &blockResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	blockTime := GetUserDateFormat(blockResp.Block.Header.Time)
	blockHeight := blockResp.Block.Header.Height
	log.Printf("last proposed block time : %s,  height : %s", blockTime, blockHeight)

	if cfg.ValDetails.ValidatorHexAddress == blockResp.Block.Header.ProposerAddress {
		fields := map[string]interface{}{
			"height":     blockResp.Block.Header.Height,
			"block_time": blockTime,
		}
		_ = writeToInfluxDb(c, bp, "matic_last_proposed_block", map[string]string{}, fields)
	}

	// Store chainID in database
	chainID := blockResp.Block.Header.ChainID
	_ = writeToInfluxDb(c, bp, "matic_chain_id", map[string]string{}, map[string]interface{}{"chain_id": chainID})
	log.Printf("Chain ID : %s ", chainID)
}

// GetUserDateFormat to which returns date in a user friendly
func GetUserDateFormat(timeToConvert string) string {
	time, err := time.Parse(time.RFC3339, timeToConvert)
	if err != nil {
		fmt.Println("Error while converting date ", err)
	}
	date := time.Format("Mon Jan _2 15:04:05 2006")
	fmt.Println("Converted time into date format : ", date)
	return date
}
