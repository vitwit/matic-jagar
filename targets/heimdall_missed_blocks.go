package targets

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/alerter"
	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// SendSingleMissedBlockAlert is to alert about single missed block if threshold value is 1 and also stores it in db
func SendSingleMissedBlockAlert(ops types.HTTPOptions, cfg *config.Config, c client.Client, cbh string) error {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
	}

	if cfg.AlertingThresholds.MissedBlocksThreshold == 1 && strings.ToUpper(cfg.AlerterPreferences.MissedBlockAlerts) == "YES" {
		err = alerter.SendTelegramAlert(fmt.Sprintf("%s validator missed a block at block height %s", cfg.ValDetails.ValidatorName, cbh), cfg)
		err = alerter.SendEmailAlert(fmt.Sprintf("%s validator missed a block at block height %s", cfg.ValDetails.ValidatorName, cbh), cfg)
		err = writeToInfluxDb(c, bp, "heimdall_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": cbh, "range": cbh})
		err = writeToInfluxDb(c, bp, "matic_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
		err = writeToInfluxDb(c, bp, "heimdall_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})

		return err
	}
	err = writeToInfluxDb(c, bp, "heimdall_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh})
	err = writeToInfluxDb(c, bp, "heimdall_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
	if err != nil {
		return err
	}

	return nil
}

// MissedBlocks is to get the current block precommits and checks whether the validator is signed the block or not
// if not signed then it will be considered as missed block and stores it in db
// Alerter will notify when the missed blocks count reaches to the configured threshold
func MissedBlocks(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	b, err := scraper.LatestBlock(ops)
	if err != nil {
		log.Printf("Error in get missed blocks: %v", err)
		return
	}

	if &b != nil {
		addrExists := false
		for _, c := range b.Block.LastCommit.Precommits {
			if strings.EqualFold(c.ValidatorAddress, cfg.ValDetails.ValidatorHexAddress) {
				addrExists = true
			}
		}

		cbh := b.Block.Header.Height

		log.Printf("Address exists :%v, and height : %s", addrExists, cbh)

		if !addrExists {

			blocks := GetContinuousMissedBlock(cfg, c)
			currentHeightFromDb := GetlatestCurrentHeightFromDB(cfg, c)
			blocksArray := strings.Split(blocks, ",")
			fmt.Println("blocks length ", int64(len(blocksArray)), currentHeightFromDb)
			// calling function to store single blocks
			err = SendSingleMissedBlockAlert(ops, cfg, c, cbh)
			if err != nil {
				log.Printf("Error while sending missed block alert: %v", err)

			}
			if cfg.AlertingThresholds.MissedBlocksThreshold > 1 && strings.ToUpper(cfg.AlerterPreferences.MissedBlockAlerts) == "YES" {
				if int64(len(blocksArray))-1 >= cfg.AlertingThresholds.MissedBlocksThreshold {
					missedBlocks := strings.Split(blocks, ",")
					_ = alerter.SendTelegramAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValDetails.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
					_ = alerter.SendEmailAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValDetails.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
					_ = writeToInfluxDb(c, bp, "heimdall_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": blocks, "range": missedBlocks[0] + " - " + missedBlocks[len(missedBlocks)-2]})
					_ = writeToInfluxDb(c, bp, "matic_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": "", "current_height": cbh})
					return
				}
				if len(blocksArray) == 1 {
					blocks = cbh + ","
				} else {
					rpcBlockHeight, _ := strconv.Atoi(cbh)
					dbBlockHeight, _ := strconv.Atoi(currentHeightFromDb)
					diff := rpcBlockHeight - dbBlockHeight
					if diff == 1 {
						blocks = blocks + cbh + ","
					} else if diff > 1 {
						blocks = ""
					}
				}
				_ = writeToInfluxDb(c, bp, "matic_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": blocks, "current_height": cbh})
				return

			}
		}
	} else {
		log.Println("Got an empty response from external rpc block dataa...")
	}
}

// GetContinuousMissedBlock returns the latest missed block from the db
func GetContinuousMissedBlock(cfg *config.Config, c client.Client) string {
	var blocks string
	q := client.NewQuery("SELECT last(block_height) FROM matic_missed_blocks", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						blocks = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}
	return blocks
}

// GetlatestCurrentHeightFromDB returns latest current height from db
func GetlatestCurrentHeightFromDB(cfg *config.Config, c client.Client) string {
	var currentHeight string
	q := client.NewQuery("SELECT last(current_height) FROM matic_missed_blocks", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						currentHeight = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}
	return currentHeight
}
