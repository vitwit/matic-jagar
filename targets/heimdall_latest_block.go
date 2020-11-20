package targets

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
	"github.com/vitwit/matic-jagar/utils"
)

// LatestProposedBlockAndTime is to get latest proposed block height, time and checks
// whether the validator hex address is equals to proposals address if yes then stores in it db
// Also stores latest block height and chain id in db
func LatestProposedBlockAndTime(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	blockResp, err := scraper.LatestBlock(ops)
	if err != nil {
		log.Printf("Error in latest proposed block and time: %v", err)
		return
	}

	time := blockResp.Block.Header.Time
	blockTime := utils.GetUserDateFormat(time) //convert time to user readable format
	blockHeight := blockResp.Block.Header.Height
	log.Printf("last proposed block time : %s,  height : %s", blockTime, blockHeight)

	if strings.EqualFold(cfg.ValDetails.ValidatorHexAddress, blockResp.Block.Header.ProposerAddress) {
		fields := map[string]interface{}{
			"height":     blockHeight,
			"block_time": blockTime,
		}
		_ = writeToInfluxDb(c, bp, "heimdall_last_proposed_block", map[string]string{}, fields)
	}

	_ = writeToInfluxDb(c, bp, "heimdall_lastest_block", map[string]string{}, map[string]interface{}{"height": blockHeight, "block_time": time})

	// Store chainID in database
	chainID := blockResp.Block.Header.ChainID
	_ = writeToInfluxDb(c, bp, "heimdall_chain_id", map[string]string{}, map[string]interface{}{"chain_id": chainID})
	log.Printf("Chain ID : %s ", chainID)
}

// GetPrevBlockTime returns time of the pevious block
func GetPrevBlockTime(cfg *config.Config, c client.Client, height string) string {
	var t string
	q := client.NewQuery(fmt.Sprintf("SELECT last(block_time) FROM heimdall_lastest_block WHERE height = '%s'", height), cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						value := r.Series[0].Values[0][idx]
						t = fmt.Sprintf("%v", value)
						break
					}
				}
			}
		}
	}
	return t
}

// BlockTimeDifference is to calcualte block time difference of prev block and current block and stores in db
func BlockTimeDifference(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	currentBlockResp, err := scraper.LatestBlock(ops)
	if err != nil {
		log.Printf("Error in block time difference: %v", err)
		return
	}

	currentBlockHeight, _ := strconv.Atoi(currentBlockResp.Block.Header.Height) //covert string to int
	prevBlockHeight := currentBlockHeight - 1
	prevBlockTime := GetPrevBlockTime(cfg, c, strconv.Itoa(prevBlockHeight)) // get previous block time
	currentBlockTime := currentBlockResp.Block.Header.Time

	if currentBlockHeight-prevBlockHeight == 1 {

		convertedCurrentTime, _ := time.Parse(time.RFC3339, currentBlockTime)
		conevrtedPrevBlockTime, _ := time.Parse(time.RFC3339, prevBlockTime)
		timeDiff := convertedCurrentTime.Sub(conevrtedPrevBlockTime)
		diffSeconds := fmt.Sprintf("%.2f", timeDiff.Seconds())

		_ = writeToInfluxDb(c, bp, "heimdall_block_time_diff", map[string]string{}, map[string]interface{}{"time_diff": diffSeconds})
		log.Printf("time diff: %s", diffSeconds)
	}

}
