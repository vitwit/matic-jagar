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

// BorNetworkHeight which returns the network height of bor
func BorNetworkHeight(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp.Body != nil {
		var cbh BorValHeight
		err = json.Unmarshal(resp.Body, &cbh)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		networkHeight, err := HexToIntConversion(cbh.Result)
		if err != nil {
			log.Printf("Error while converting bor n/w height from hex to int : %v", err)
			return
		}

		_ = writeToInfluxDb(c, bp, "bor_network_height", map[string]string{}, map[string]interface{}{"block_height": networkHeight, "height_in_hex": cbh.Result})
		log.Printf("Bor Network Block Height: %d", networkHeight)

		// Calling function to get validator latest
		// block height
		validatorHeight := GetBorCurrentBlokHeight(cfg, c)
		if validatorHeight == "" {
			log.Println("Error while fetching validator block height of bor from db ", validatorHeight)
			return
		}

		vaidatorBlockHeight, _ := strconv.Atoi(validatorHeight)
		heightDiff := networkHeight - vaidatorBlockHeight

		_ = writeToInfluxDb(c, bp, "bor_height_difference", map[string]string{}, map[string]interface{}{"difference": heightDiff})
		log.Printf("BOR :: Network height: %d and Validator Height: %d", networkHeight, vaidatorBlockHeight)

		// Send alert
		if int64(heightDiff) >= cfg.AlertingThresholds.BlockDiffThreshold && strings.ToUpper(cfg.ChooseAlerts.BlockDiffAlerts) == "YES" {
			_ = SendTelegramAlert(fmt.Sprintf("Bor Block Difference Alert: Block Difference between network and validator has exceeded %d", cfg.AlertingThresholds.BlockDiffThreshold), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Bor Block Difference Alert : Block difference between network and validator has exceeded %d", cfg.AlertingThresholds.BlockDiffThreshold), cfg)
			log.Println("Sent alert of block height difference")
		}
	}

}
