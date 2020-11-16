package targets

import (
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// BorPendingTransactions is to get the pending transactions of bor and stores in db
func BorPendingTransactions(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	txs, err := scraper.BorPendingTransactions(ops)
	if err != nil {
		log.Printf("Error in bor pending transactions: %v", err)
		return
	}

	if &txs != nil {
		pendingTxns := len(txs.Result)

		_ = writeToInfluxDb(c, bp, "bor_pending_txns", map[string]string{}, map[string]interface{}{"pending_txns": pendingTxns})
		log.Printf("Pending Transactions: %d", pendingTxns)
	}
}
