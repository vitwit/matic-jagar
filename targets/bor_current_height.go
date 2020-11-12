package targets

import (
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// BorCurrentHeight which returns the current height of bor validator
func BorCurrentHeight(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	cbh, err := scraper.EthBlockNumber(ops)
	if err != nil {
		log.Printf("Error in BorCurrentHeight: %v", err)
		return
	}

	if &cbh != nil {

		height, err := HexToIntConversion(cbh.Result)
		if err != nil {
			log.Printf("Error while converting bor current height from hex to int : %v", err)
			return
		}

		_ = writeToInfluxDb(c, bp, "bor_current_height", map[string]string{}, map[string]interface{}{"block_height": height, "height_in_hex": cbh.Result})
		log.Printf("Bor Current Block Height: %d", height)
	}

}

// GetBorCurrentBlokHeightInHex returns current block height of bor from db
func GetBorCurrentBlokHeightInHex(cfg *config.Config, c client.Client) string {
	var validatorHeight string
	q := client.NewQuery("SELECT last(height_in_hex) FROM bor_current_height", cfg.InfluxDB.Database, "")
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

// GetBorCurrentBlokHeight returns current block height of bor from db
func GetBorCurrentBlokHeight(cfg *config.Config, c client.Client) string {
	var validatorHeight string
	q := client.NewQuery("SELECT last(block_height) FROM bor_current_height", cfg.InfluxDB.Database, "")
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
