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

// ValidatorCaughtUp is to get validator syncing status
func ValidatorCaughtUp(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var sync Caughtup
	err = json.Unmarshal(resp.Body, &sync)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var synced int
	caughtUp := !sync.Syncing
	if !caughtUp {
		if strings.ToUpper(cfg.ChooseAlerts.NodeSyncAlert) == "YES" {
			_ = SendTelegramAlert("Your validator node is not synced!", cfg)
			_ = SendEmailAlert("Your validator node is not synced!", cfg)
		}
		synced = 0
	} else {
		synced = 1
	}

	_ = writeToInfluxDb(c, bp, "heimdall_val_caughtup", map[string]string{}, map[string]interface{}{"synced": synced})
	log.Printf("Heimdall Valiator Caught UP: %v", sync.Syncing)
}

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
		p2, err := createDataPoint("heimdall_current_block_height", map[string]string{}, map[string]interface{}{"height": bh})
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

	p3, err := createDataPoint("heimdall_node_synced", map[string]string{}, map[string]interface{}{"status": synced})
	if err == nil {
		pts = append(pts, p3)
	}

	bp.AddPoints(pts)
	_ = writeBatchPoints(c, bp)
	log.Printf("\nCurrent Block Height: %s", currentBlockHeight)

	// Store validator details such as moniker, signer address and hex address
	moniker := status.Result.NodeInfo.Moniker
	hexAddress := status.Result.ValidatorInfo.Address
	signerAddress := cfg.ValDetails.SignerAddress
	_ = writeToInfluxDb(c, bp, "heimdall_val_desc", map[string]string{}, map[string]interface{}{"moniker": moniker, "hex_address": hexAddress, "signer_address": signerAddress, "address": signerAddress[2:]})

	log.Printf("Moniker:%s ", moniker)
}

// GetValidatorBlock returns validator current block height
func GetValidatorBlock(cfg *config.Config, c client.Client) string {
	var validatorHeight string
	q := client.NewQuery("SELECT last(height) FROM heimdall_current_block_height", cfg.InfluxDB.Database, "")
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

// GetNodeSync returns the syncing status of a node
func GetNodeSync(cfg *config.Config, c client.Client) string {
	var status, sync string
	q := client.NewQuery("SELECT last(status) FROM heimdall_node_synced", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						s := r.Series[0].Values[0][idx]
						sync = fmt.Sprintf("%v", s)
						break
					}
				}
			}
		}
	}

	if sync == "1" {
		status = "synced"
	} else {
		status = "not synced"
	}

	return status
}
