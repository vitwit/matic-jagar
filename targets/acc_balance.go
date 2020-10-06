package targets

import (
	"encoding/json"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/vitwit/matic-jagar/config"
)

// GetAccountBal to get account balance information using signer address
func GetAccountBal(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		_ = writeToInfluxDb(c, bp, "matic_account_balance", map[string]string{}, map[string]interface{}{"balance": "NA"})
		return
	}

	var accResp AccountBalResp
	err = json.Unmarshal(resp.Body, &accResp)
	if err != nil {
		log.Printf("Error: %v", err)
		_ = writeToInfluxDb(c, bp, "matic_account_balance", map[string]string{}, map[string]interface{}{"balance": "NA"})
		return
	}

	if len(accResp.Result) > 0 {
		addressBalance := convertToCommaSeparated(accResp.Result[0].Amount) + accResp.Result[0].Denom
		_ = writeToInfluxDb(c, bp, "matic_account_balance", map[string]string{}, map[string]interface{}{"balance": addressBalance})
		log.Printf("Address Balance: %s", addressBalance)
	}
}

func convertToCommaSeparated(amt string) string {
	a, err := strconv.Atoi(amt)
	if err != nil {
		return amt
	}
	p := message.NewPrinter(language.English)
	return p.Sprintf("%d", a)
}
