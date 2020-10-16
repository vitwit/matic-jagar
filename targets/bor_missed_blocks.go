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

// SendSingleMissedBlockAlert sends missed block alert to telegram and email accounts
func SendBorSingleMissedBlockAlert(ops HTTPOptions, cfg *config.Config, c client.Client, cbh string) error {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
	}

	if cfg.AlertingThresholds.MissedBlocksThreshold == 1 {
		if strings.ToUpper(cfg.ChooseAlerts.MissedBlockAlerts) == "YES" {
			err = SendTelegramAlert(fmt.Sprintf("%s validator on bor node missed a block at block height %s", cfg.ValDetails.ValidatorName, cbh), cfg)
			err = SendEmailAlert(fmt.Sprintf("%s validator on bor node missed a block at block height %s", cfg.ValDetails.ValidatorName, cbh), cfg)
			err = writeToInfluxDb(c, bp, "matic_bor_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": cbh, "range": cbh})
			err = writeToInfluxDb(c, bp, "matic_bor_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
			err = writeToInfluxDb(c, bp, "matic_bor_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
			if err != nil {
				return err
			}
		}

	} else {
		err = writeToInfluxDb(c, bp, "matic_bor_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
		if err != nil {
			return err
		}
	}

	return nil
}

// GetBorMissedBlocks sends alerts of missed blocks according to the threshold given by user
func GetBorMissedBlocks(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	borHeight := GetBorCurrentBlokHeightInHex(cfg, c)

	ops.Body.Params = append(ops.Body.Params, borHeight)

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var signers BorSignersRes
	err = json.Unmarshal(resp.Body, &signers)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	height := HexToIntConversion(borHeight)

	cbh := strconv.Itoa(height)

	if signers.Result != nil {
		addrExists := false

		for _, c := range signers.Result {
			if strings.ToUpper(c) == strings.ToUpper(cfg.ValDetails.SignerAddress) {
				addrExists = true
			}
		}

		log.Println("address exists and height......", addrExists, cbh)

		if !addrExists {

			// // Calling SendEmeregencyAlerts to send emeregency alerts
			// err := SendEmeregencyAlerts(cfg, c, cbh)
			// if err != nil {
			// 	log.Println("Error while sending emeregecny missed block alerts...", err)
			// }

			blocks := GetBorContinuousMissedBlock(cfg, c)
			currentHeightFromDb := GetBorlatestCurrentHeightFromDB(cfg, c)
			blocksArray := strings.Split(blocks, ",")
			fmt.Println("blocks length ", int64(len(blocksArray)), currentHeightFromDb)
			// calling function to store single blocks
			err = SendBorSingleMissedBlockAlert(ops, cfg, c, cbh)
			if err != nil {
				log.Printf("Error while sending missed block alert: %v", err)

			}
			if cfg.AlertingThresholds.MissedBlocksThreshold > 1 && strings.ToUpper(cfg.ChooseAlerts.MissedBlockAlerts) == "YES" {
				if int64(len(blocksArray))-1 >= cfg.AlertingThresholds.MissedBlocksThreshold {
					missedBlocks := strings.Split(blocks, ",")
					_ = SendTelegramAlert(fmt.Sprintf("%s validator on bor node missed blocks from height %s to %s", cfg.ValDetails.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
					_ = SendEmailAlert(fmt.Sprintf("%s validator on bor node missed blocks from height %s to %s", cfg.ValDetails.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
					_ = writeToInfluxDb(c, bp, "matic_bor_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": blocks, "range": missedBlocks[0] + " - " + missedBlocks[len(missedBlocks)-2]})
					_ = writeToInfluxDb(c, bp, "matic_bor_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": "", "current_height": cbh})
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
				_ = writeToInfluxDb(c, bp, "matic_bor_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": blocks, "current_height": cbh})
				return

			}
		} else {
			_ = writeToInfluxDb(c, bp, "matic_bor_val_signed_blocks", map[string]string{}, map[string]interface{}{"signed_block_height": cbh})
		}
	} else {
		log.Println("Got an empty response from the rpc...")
	}
}

// GetContinuousMissedBlock returns the latest missed block from the db
func GetBorContinuousMissedBlock(cfg *config.Config, c client.Client) string {
	var blocks string
	q := client.NewQuery("SELECT last(block_height) FROM matic_bor_missed_blocks", cfg.InfluxDB.Database, "")
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
func GetBorlatestCurrentHeightFromDB(cfg *config.Config, c client.Client) string {
	var currentHeight string
	q := client.NewQuery("SELECT last(current_height) FROM matic_bor_missed_blocks", cfg.InfluxDB.Database, "")
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
