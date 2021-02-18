package targets

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/alerter"
	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// NodeVersion is to get application version and stores in db
func NodeVersion(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	applicationInfo, err := scraper.GetVersion(ops)
	if err != nil {
		log.Printf("Error in node version: %v", err)
		return
	}

	appVersion := applicationInfo.ApplicationVersion.Version
	if appVersion == "" {
		return
	}

	_ = db.WriteToInfluxDb(c, bp, "heimdall_version", map[string]string{}, map[string]interface{}{"v": appVersion})
	log.Printf("Version: %s", appVersion)
}

// ValidatorCaughtUp is to get validator syncing status and stores it in db
// Alerter will alerts when the node is not synced
func ValidatorCaughtUp(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	sync, err := scraper.GetCaughtUpStatus(ops)
	if err != nil {
		log.Printf("Error in validator caughtup: %v", err)
		return
	}

	var synced int
	caughtUp := !sync.Syncing
	if !caughtUp {
		if strings.ToUpper(cfg.AlerterPreferences.NodeSyncAlert) == "YES" {
			_ = alerter.SendTelegramAlert("⚠️ Your heimdall validator node is not synced!", cfg)
			_ = alerter.SendEmailAlert("⚠️ Your heimdall validator node is not synced!", cfg)
		}
		synced = 0
	} else {
		synced = 1
	}

	_ = db.WriteToInfluxDb(c, bp, "heimdall_node_synced", map[string]string{}, map[string]interface{}{"synced": synced})
	log.Printf("Heimdall Valiator Caught UP: %v", sync.Syncing)
}

// Status is to get response from rpc /status endpoint and stores node status
// block height and operator info
// Alerter will notify about the node status i.e., validator instance is running or not by checking status resonse
func Status(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	status, err := scraper.GetStatus(ops)
	if err != nil {
		_ = db.WriteToInfluxDb(c, bp, "heimdall_node_status", map[string]string{}, map[string]interface{}{"status": 0})

		log.Printf("Validator Error: %v", err)
		return
	}

	if &status.Result == nil {
		if strings.ToUpper(cfg.AlerterPreferences.NodeStatusAlert) == "YES" {
			_ = alerter.SendTelegramAlert("⚠️ Your heimdall validator instance is not running", cfg)
			_ = alerter.SendEmailAlert("⚠️ Your heimdall validator instance is not running", cfg)
		}
		_ = db.WriteToInfluxDb(c, bp, "heimdall_node_status", map[string]string{}, map[string]interface{}{"status": 0})
		return
	}

	err = db.WriteToInfluxDb(c, bp, "heimdall_node_status", map[string]string{}, map[string]interface{}{"status": 1})
	if err != nil {
		log.Printf("Error while writing node status into db : %v ", err)
	}

	var bh int
	currentBlockHeight := status.Result.SyncInfo.LatestBlockHeight
	if currentBlockHeight != "" {
		bh, _ = strconv.Atoi(currentBlockHeight)
		err = db.WriteToInfluxDb(c, bp, "heimdall_current_block_height", map[string]string{}, map[string]interface{}{"height": bh})
		if err != nil {
			log.Printf("Error while stroing current block height : %v", err)
		}
	}

	// Store validator details such as moniker, signer address and hex address
	moniker := status.Result.NodeInfo.Moniker
	hexAddress := status.Result.ValidatorInfo.Address
	signerAddress := cfg.ValDetails.SignerAddress
	_ = db.WriteToInfluxDb(c, bp, "heimdall_val_desc", map[string]string{}, map[string]interface{}{"moniker": moniker, "hex_address": hexAddress, "signer_address": signerAddress, "address": signerAddress[2:]})

	log.Printf("Moniker:%s ", moniker)
}

// GetValidatorBlock returns validator current block height from db
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

// GetNodeSync returns the syncing status of a node from db
func GetNodeSync(cfg *config.Config, c client.Client) string {
	var status, sync string
	q := client.NewQuery("SELECT last(synced) FROM heimdall_node_synced", cfg.InfluxDB.Database, "")
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
