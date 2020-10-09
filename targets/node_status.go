package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetNodeStatus to get reponse of validator status like
//current block height and node status
func GetNodeStatus(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}
	var pts []*client.Point

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var status Status
	err = json.Unmarshal(resp.Body, &status)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var bh int
	currentBlockHeight := status.Result.SyncInfo.LatestBlockHeight
	if currentBlockHeight != "" {
		bh, _ = strconv.Atoi(currentBlockHeight)
		p2, err := createDataPoint("matic_current_block_height", map[string]string{}, map[string]interface{}{"height": bh})
		if err == nil {
			pts = append(pts, p2)
		}
	}

	var synced int
	caughtUp := !status.Result.SyncInfo.CatchingUp
	if !caughtUp {
		_ = SendTelegramAlert("Your validator node is not synced!", cfg)
		_ = SendEmailAlert("Your validator node is not synced!", cfg)
		synced = 0
	} else {
		synced = 1
	}
	p3, err := createDataPoint("matic_node_synced", map[string]string{}, map[string]interface{}{"status": synced})
	if err == nil {
		pts = append(pts, p3)
	}

	bp.AddPoints(pts)
	_ = writeBatchPoints(c, bp)
	log.Printf("\nCurrent Block Height: %s \nCaught Up? %t \n",
		currentBlockHeight, caughtUp)

	// Store validator details such as moniker, signer address and hex address
	moniker := status.Result.NodeInfo.Moniker
	hexAddress := status.Result.ValidatorInfo.Address
	signerAddress := cfg.SignerAddress
	_ = writeToInfluxDb(c, bp, "matic_val_desc", map[string]string{}, map[string]interface{}{"moniker": moniker, "hex_address": hexAddress, "signer_address": signerAddress})

	log.Printf("Moniker:%s ", moniker)
}

// GetValidatorBlock returns validator current block height
func GetValidatorBlock(cfg *config.Config, c client.Client) string {
	var validatorHeight string
	q := client.NewQuery("SELECT last(height) FROM matic_current_block_height", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						validatorHeight = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}
	return validatorHeight
}

func BorCurrentHeight(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		var cbh BorResult
		err = json.Unmarshal(resp.Body, &cbh)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		height := HexToIntConversion(cbh.Result)

		_ = writeToInfluxDb(c, bp, "matic_bor_current_height", map[string]string{}, map[string]interface{}{"block_height": height})
		log.Printf("Bor Current Block Height: %d", height)
	}

}

func BorEthSyncing(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	// log.Fatalf("resp..", resp)

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		var sync BorResult
		err = json.Unmarshal(resp.Body, &sync)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		// log.Fatalf("result..", resp.Body)

		var synced int

		_ = writeToInfluxDb(c, bp, "matic_bor_node_synced", map[string]string{}, map[string]interface{}{"status": synced})
		log.Printf("Bor Syncing Status: %d", synced)
	}
}
