package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

func GetBorPendingTransactions(ops HTTPOptions, cfg *config.Config, c client.Client) {
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
		var txs EthPendingTransactions
		err = json.Unmarshal(resp.Body, &txs)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		pendingTxns := len(txs.Result)

		_ = writeToInfluxDb(c, bp, "matic_bor_pending_txns", map[string]string{}, map[string]interface{}{"pending_txns": pendingTxns})
		log.Printf("Pending Transactions: %d", pendingTxns)
	}
}
