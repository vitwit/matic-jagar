package targets

import (
	// "encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// UnconfimedTxns is to get the no of uncofirmed txns and stores it in db
func UnconfimedTxns(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	unconfirmedTxns, err := scraper.GetUnconfirmedTxs(ops)
	if err != nil {
		log.Printf("Error in unconfirmed txs: %v", err)
		return
	}

	if &unconfirmedTxns.Result == nil {
		log.Println("Got an empty response from validator rpc !")
		return
	}

	totalUnconfirmedTxns := unconfirmedTxns.Result.Total

	_ = writeToInfluxDb(c, bp, "heimdall_unconfirmed_txns", map[string]string{}, map[string]interface{}{"unconfirmed_txns": totalUnconfirmedTxns})
	log.Printf("No of unconfirmed txns: %s", totalUnconfirmedTxns)
}

// ValidatorGas is to get validator max tx gas and stores in db
func ValidatorGas(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
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
