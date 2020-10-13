package targets

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/vitwit/matic-jagar/config"
)

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

		_ = writeToInfluxDb(c, bp, "matic_bor_current_height", map[string]string{}, map[string]interface{}{"block_height": height, "height_in_hex": cbh.Result})
		log.Printf("Bor Current Block Height: %d", height)
	}

}

// GetBorCurrentBlokHeight returns current block height of bor from db
func GetBorCurrentBlokHeight(cfg *config.Config, c client.Client) string {
	var validatorHeight string
	q := client.NewQuery("SELECT last(height_in_hex) FROM matic_bor_current_height", cfg.InfluxDB.Database, "")
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
