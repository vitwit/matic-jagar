package targets

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/alerter"
	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
	"github.com/vitwit/matic-jagar/utils"
)

// BorNetworkHeight is to get the network height of bor and stores in db
// Alerter will notify whenever bock diff of network and validator reaches configured threshold
func BorNetworkHeight(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	cbh, err := scraper.BorValidatorHeight(ops)
	if err != nil {
		log.Printf("Error in BorNetworkHeight: %v", err)
		return
	}

	if &cbh != nil {

		networkHeight, err := utils.HexToIntConversion(cbh.Result)
		if err != nil {
			log.Printf("Error while converting bor n/w height from hex to int : %v", err)
			return
		}

		_ = db.WriteToInfluxDb(c, bp, "bor_network_height", map[string]string{}, map[string]interface{}{"block_height": networkHeight, "height_in_hex": cbh.Result})
		log.Printf("Bor Network Block Height: %d", networkHeight)

		// Get validator block height
		HTTPOptions := types.HTTPOptions{
			Endpoint: cfg.Endpoints.BorRPCEndpoint,
			Method:   http.MethodPost,
			Body:     types.Payload{Jsonrpc: "2.0", Method: "eth_blockNumber", ID: 83},
		}

		cbh, err := scraper.EthBlockNumber(HTTPOptions)
		if err != nil {
			log.Printf("Error in bor validator height : %v", err)
			return
		}

		if &cbh == nil || cbh.Result == "" {
			log.Println("Got an empty response from bor rpc !")
			return
		}

		vaidatorBlockHeight, err := utils.HexToIntConversion(cbh.Result)
		if err != nil {
			log.Printf("Error while converting bor current height from hex to int : %v", err)
			return
		}

		heightDiff := networkHeight - vaidatorBlockHeight

		_ = db.WriteToInfluxDb(c, bp, "bor_height_difference", map[string]string{}, map[string]interface{}{"difference": heightDiff})
		log.Printf("BOR :: Network height: %d and Validator Height: %d", networkHeight, vaidatorBlockHeight)

		// Send alert
		if int64(heightDiff) >= cfg.AlertingThresholds.BlockDiffThreshold && strings.ToUpper(cfg.AlerterPreferences.BlockDiffAlerts) == "YES" {
			_ = alerter.SendTelegramAlert(fmt.Sprintf("Bor Block Difference Alert: Block Difference between network and validator has exceeded %d", cfg.AlertingThresholds.BlockDiffThreshold), cfg)
			_ = alerter.SendEmailAlert(fmt.Sprintf("Bor Block Difference Alert : Block difference between network and validator has exceeded %d", cfg.AlertingThresholds.BlockDiffThreshold), cfg)
			log.Println("Sent alert of bor block height difference")
		}
	} else {
		log.Println("Got an empty response from bor external rpc !")
		return
	}

}
