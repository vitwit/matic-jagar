package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

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
		amount := ConvertToMatic(accResp.Result[0].Amount) // curent amount
		prevAmount := GetAccountBalFromDb(cfg, c)          // amount from db

		if prevAmount != amount {
			if strings.ToUpper(cfg.ChooseAlerts.BalanceChangeAlerts) == "YES" {
				_ = SendTelegramAlert(fmt.Sprintf("Heimdall Balance Change Alert : Your account balance has changed from  %s to %s", prevAmount, amount), cfg)
				_ = SendEmailAlert(fmt.Sprintf("Heimdall Balance Change Alert : Your account balance has changed from  %s to %s", prevAmount, amount), cfg)
			}
		}

		addressBalance := convertToCommaSeparated(amount) + accResp.Result[0].Denom
		_ = writeToInfluxDb(c, bp, "matic_heimdall_current_balance", map[string]string{}, map[string]interface{}{"current_balance": addressBalance, "balance": amount})
		log.Printf("Address Balance: %s", addressBalance)
	}
}

// GetAccountBalFromDb returns account balance from db
func GetAccountBalFromDb(cfg *config.Config, c client.Client) string {
	var balance string
	q := client.NewQuery("SELECT last(balance) FROM matic_heimdall_current_balance", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						amount := r.Series[0].Values[0][idx]
						balance = fmt.Sprintf("%v", amount)
						break
					}
				}
			}
		}
	}
	return balance
}

// GetAccountBalWithDenomFromdb returns account balance from db
func GetAccountBalWithDenomFromdb(cfg *config.Config, c client.Client) string {
	var balance string
	q := client.NewQuery("SELECT last(current_balance) FROM matic_heimdall_current_balance", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						amount := r.Series[0].Values[0][idx]
						balance = fmt.Sprintf("%v", amount)
						break
					}
				}
			}
		}
	}
	return balance
}
