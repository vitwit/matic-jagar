package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetHeimdallCurrentBal to get current balance information using signer address
func GetHeimdallCurrentBal(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		_ = writeToInfluxDb(c, bp, "matic_heimdall_current_balance", map[string]string{}, map[string]interface{}{"current_balance": "NA"})
		return
	}

	var accResp AccountBalResp
	err = json.Unmarshal(resp.Body, &accResp)
	if err != nil {
		log.Printf("Error: %v", err)
		_ = writeToInfluxDb(c, bp, "matic_heimdall_current_balance", map[string]string{}, map[string]interface{}{"current_balance": "NA"})
		return
	}

	if len(accResp.Result) > 0 {
		addressBalance := convertToCommaSeparated(ConvertToMatic(accResp.Result[0].Amount)) + accResp.Result[0].Denom
		_ = writeToInfluxDb(c, bp, "matic_heimdall_current_balance", map[string]string{}, map[string]interface{}{"current_balance": addressBalance})
		log.Printf("Address Balance: %s", addressBalance)
	}
}
