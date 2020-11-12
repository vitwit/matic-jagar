package targets

import (
	// "encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// GetUnconfimedTxns to get the no of uncofirmed txns
func GetUnconfimedTxns(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
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

	totalUnconfirmedTxns := unconfirmedTxns.Result.Total

	_ = writeToInfluxDb(c, bp, "heimdall_unconfirmed_txns", map[string]string{}, map[string]interface{}{"unconfirmed_txns": totalUnconfirmedTxns})
	log.Printf("No of unconfirmed txns: %s", totalUnconfirmedTxns)
}
