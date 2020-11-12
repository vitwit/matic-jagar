package targets

import (
	"fmt"
	"log"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// GetBlockTimeDifference to calculate block time difference of prev block and current block
func GetBlockTimeDifference(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
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
