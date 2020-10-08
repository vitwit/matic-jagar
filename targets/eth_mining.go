package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// BorEthMining checks if client is actively mining new blocks or not
func BorEthMining(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp.Body != nil {
		var bnl BorBoolResp
		err = json.Unmarshal(resp.Body, &bnl)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		mining := bnl.Result

		_ = writeToInfluxDb(c, bp, "matic_bot_eth_mining", map[string]string{}, map[string]interface{}{"eth_mining": mining})
		log.Println("Bor Eth Mining: ", mining)
	}
}
