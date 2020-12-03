package targets

import (
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// BorPendingTransactions is to get the pending transactions of bor and stores in db
func BorPendingTransactions(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	txs, err := scraper.BorPendingTransactions(ops)
	if err != nil {
		log.Printf("Error while getting bor pending transactions: %v", err)
		return
	}

	if &txs != nil {
		pendingTxns := len(txs.Result)

		err = db.WriteToInfluxDb(c, bp, "bor_pending_txns", map[string]string{}, map[string]interface{}{"pending_txns": pendingTxns})
		if err != nil {
			log.Printf("Error while writing bor pending txns into db : %v", err)
		}
		log.Printf("Pending Transactions: %d", pendingTxns)
	} else {
		log.Println("Got an empty response from bor rpc !")
		return
	}
}
