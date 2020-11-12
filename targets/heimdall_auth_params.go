package targets

import (
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// GetValidatorGas is to get validator max tx gas
func GetValidatorGas(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	authParam, err := scraper.AuthParams(ops)
	if err != nil {
		log.Printf("Error in validator gas: %v", err)
		return
	}

	maxTxGas := authParam.Result.MaxTxGas

	_ = writeToInfluxDb(c, bp, "heimdall_auth_params", map[string]string{}, map[string]interface{}{"max_tx_gas": maxTxGas})
	log.Printf("Max tx gas: %d\n", maxTxGas)
}
