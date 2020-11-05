package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetValidatorGas is to get validator max tx gas
func GetValidatorGas(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var authParam AuthParams
	err = json.Unmarshal(resp.Body, &authParam)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	maxTxGas := authParam.Result.MaxTxGas

	_ = writeToInfluxDb(c, bp, "heimdall_auth_params", map[string]string{}, map[string]interface{}{"max_tx_gas": maxTxGas})
	log.Printf("Max tx gas: %d\n", maxTxGas)
}
