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
	"github.com/vitwit/matic-jagar/utils"
)

// SendSingleMissedBlockAlert is to send signle missed block alerts and stores it in db
func SendBorSingleMissedBlockAlert(ops types.HTTPOptions, cfg *config.Config, c client.Client, cbh string) error {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
	}

	if cfg.AlertingThresholds.MissedBlocksThreshold == 1 {
		if strings.ToUpper(cfg.AlerterPreferences.MissedBlockAlerts) == "YES" {
			err = alerter.SendTelegramAlert(fmt.Sprintf("⚠️ Bor Missed Block Alert: %s validator on bor node missed a block at block height %s", cfg.ValDetails.ValidatorName, cbh), cfg)
			if err != nil {
				log.Printf("Error while sending missed blocks telegram alert : %v", err)
				return err
			}
			err = alerter.SendEmailAlert(fmt.Sprintf("⚠️ Bor Missed Block Alert: %s validator on bor node missed a block at block height %s", cfg.ValDetails.ValidatorName, cbh), cfg)
			if err != nil {
				log.Printf("Error while sending missed blocks email alert : %v", err)
				return err
			}
			err = db.WriteToInfluxDb(c, bp, "bor_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": cbh, "range": cbh})
			if err != nil {
				log.Printf("Error while storing continuous missed blocks : %v", err)
				return err
			}
			err = db.WriteToInfluxDb(c, bp, "matic_bor_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
			if err != nil {
				log.Printf("Error while storing missed blocks : %v", err)
				return err
			}
			err = db.WriteToInfluxDb(c, bp, "bor_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
			if err != nil {
				log.Printf("Error while stroing missed blocks : %v", err)
				return err
			}
		}

	} else {
		err = db.WriteToInfluxDb(c, bp, "bor_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh})
		if err != nil {
			log.Printf("Error while stroing missed blocks : %v", err)
			return err
		}
		err = db.WriteToInfluxDb(c, bp, "bor_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
		if err != nil {
			log.Printf("Error while stroing total missed blocks : %v", err)
			return err
		}
	}

	return nil
}

// BorMissedBlocks is to get the current block precommits and checks whether the validator is signed the block or not
// if not signed then it will be considered as missed block and stores it in db
// Alerter will notify when the missed blocks count reaches to the configured threshold
func BorMissedBlocks(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	borHeight := GetBorCurrentBlokHeightInHex(cfg, c)
	if borHeight == "" {
		log.Println("Got empty block height of bor from db")
		return
	}

	ops.Body.Params = append(ops.Body.Params, borHeight)

	signers, err := scraper.BorSignersRes(ops)
	if err != nil {
		log.Printf("Error in BorMissedBlocks: %v", err)
		return
	}

	height, err := utils.HexToIntConversion(borHeight)
	if err != nil {
		log.Printf("Error while converting bor height from hex to int : %v", err)
		return
	}

	cbh := strconv.Itoa(height)

	if signers.Result != nil {
		isSigned := false

		for _, addr := range signers.Result {
			if strings.EqualFold(addr, cfg.ValDetails.SignerAddress) {
				isSigned = true
			}
		}

		log.Printf("block signed status : %v, and height : %s ", isSigned, cbh)

		if !isSigned {
			blocks := GetBorContinuousMissedBlock(cfg, c)
			currentHeightFromDb := GetBorlatestCurrentHeightFromDB(cfg, c)
			blocksArray := strings.Split(blocks, ",")
			// calling function to store single blocks
			err = SendBorSingleMissedBlockAlert(ops, cfg, c, cbh)
			if err != nil {
				log.Printf("Error while sending missed block alert: %v", err)

			}
			if cfg.AlertingThresholds.MissedBlocksThreshold > 1 && strings.ToUpper(cfg.AlerterPreferences.MissedBlockAlerts) == "YES" {
				if int64(len(blocksArray))-1 >= cfg.AlertingThresholds.MissedBlocksThreshold {
					missedBlocks := strings.Split(blocks, ",")
					_ = alerter.SendTelegramAlert(fmt.Sprintf("⚠️ Bor Missed Blocks Alert: %s validator on bor node missed blocks from height %s to %s", cfg.ValDetails.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
					_ = alerter.SendEmailAlert(fmt.Sprintf("⚠️ Bor Missed Blocks Alert: %s validator on bor node missed blocks from height %s to %s", cfg.ValDetails.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
					_ = db.WriteToInfluxDb(c, bp, "bor_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": blocks, "range": missedBlocks[0] + " - " + missedBlocks[len(missedBlocks)-2]})
					_ = db.WriteToInfluxDb(c, bp, "matic_bor_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": "", "current_height": cbh})
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
				err = db.WriteToInfluxDb(c, bp, "matic_bor_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": blocks, "current_height": cbh})
				if err != nil {
					log.Printf("Error while storing missed blocks : %v ", err)
					return
				}
			}
		} else {
			err = db.WriteToInfluxDb(c, bp, "bor_val_signed_blocks", map[string]string{}, map[string]interface{}{"signed_block_height": cbh})
			if err != nil {
				log.Printf("Error while storing validator signed blocks : %v ", err)
				return
			}
		}
	} else {
		log.Println("Got an empty response from bor rpc !")
		return
	}
}

// GetContinuousMissedBlock returns the latest missed block from db
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
