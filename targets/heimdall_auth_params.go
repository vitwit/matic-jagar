package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetValidatorFeeAndGas is to get validator vee and max tx gas
func GetValidatorFeeAndGas(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

	valFee := authParam.Result.TxFees
	maxTxGas := authParam.Result.MaxTxGas

	_ = writeToInfluxDb(c, bp, "heimdall_auth_params", map[string]string{}, map[string]interface{}{"validator_fee": valFee, "max_tx_gas": maxTxGas})
	log.Printf("Val fee: %s and max tx gas: %d\n", valFee, maxTxGas)
}
