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
func SendSingleMissedBlockAlert(ops HTTPOptions, cfg *config.Config, c client.Client, cbh string) error {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
	}

	if cfg.MissedBlocksThreshold == 1 {
		err = SendTelegramAlert(fmt.Sprintf("%s validator missed a block at block height %s", cfg.ValidatorName, cbh), cfg)
		err = SendEmailAlert(fmt.Sprintf("%s validator missed a block at block height %s", cfg.ValidatorName, cbh), cfg)
		err = writeToInfluxDb(c, bp, "matic_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": cbh, "range": cbh})
		err = writeToInfluxDb(c, bp, "matic_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
		err = writeToInfluxDb(c, bp, "matic_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})

		return err
	}
	err = writeToInfluxDb(c, bp, "matic_total_missed_blocks", map[string]string{}, map[string]interface{}{"block_height": cbh, "current_height": cbh})
	if err != nil {
		return err
	}

	return nil
}

// GetMissedBlocks sends alerts of missed blocks according to the threshold given by user
func GetMissedBlocks(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var networkLatestBlock Status
	err = json.Unmarshal(resp.Body, &networkLatestBlock)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	cbh := networkLatestBlock.Result.SyncInfo.LatestBlockHeight

	resp, err = HitHTTPTarget(HTTPOptions{
		Endpoint:    cfg.ExternalRPC + "/block",
		QueryParams: QueryParams{"height": cbh},
		Method:      "GET",
	})
	if err != nil {
		log.Printf("Error getting details of current block: %v", err)
		return
	}

	var b CurrentBlockWithHeight
	err = json.Unmarshal(resp.Body, &b)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if &b != nil {
		addrExists := false
		for _, c := range b.Result.Block.LastCommit.Precommits {
			if c.ValidatorAddress == cfg.ValidatorHexAddress {
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

			blocks := GetContinuousMissedBlock(cfg, c)
			currentHeightFromDb := GetlatestCurrentHeightFromDB(cfg, c)
			blocksArray := strings.Split(blocks, ",")
			fmt.Println("blocks length ", int64(len(blocksArray)), currentHeightFromDb)
			// calling function to store single blocks
			err = SendSingleMissedBlockAlert(ops, cfg, c, cbh)
			if err != nil {
				log.Printf("Error while sending missed block alert: %v", err)

			}
			if cfg.MissedBlocksThreshold > 1 {
				if int64(len(blocksArray))-1 >= cfg.MissedBlocksThreshold {
					missedBlocks := strings.Split(blocks, ",")
					_ = SendTelegramAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
					_ = SendEmailAlert(fmt.Sprintf("%s validator missed blocks from height %s to %s", cfg.ValidatorName, missedBlocks[0], missedBlocks[len(missedBlocks)-2]), cfg)
					_ = writeToInfluxDb(c, bp, "matic_continuous_missed_blocks", map[string]string{}, map[string]interface{}{"missed_blocks": blocks, "range": missedBlocks[0] + " - " + missedBlocks[len(missedBlocks)-2]})
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
