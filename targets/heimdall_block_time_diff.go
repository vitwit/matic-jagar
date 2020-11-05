package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetBlockTimeDifference to calculate block time difference of prev block and current block
func GetBlockTimeDifference(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	currResp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var currentBlockResp LatestBlock
	err = json.Unmarshal(currResp.Body, &currentBlockResp)
	if err != nil {
		log.Printf("Error: %v", err)
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
