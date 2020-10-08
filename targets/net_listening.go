package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// BorNetListening checks if client is actively listening for network connections
func BorNetListening(ops HTTPOptions, cfg *config.Config, c client.Client) {
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
		var bnl NetListening
		err = json.Unmarshal(resp.Body, &bnl)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		netListen := bnl.Result

		_ = writeToInfluxDb(c, bp, "matic_bot_net_listening", map[string]string{}, map[string]interface{}{"net_listen": netListen})
		log.Println("Bor Net Listening: ", netListen)
	}
}
