package targets

import (
	"fmt"
	"log"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/alerter"
	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
	"github.com/vitwit/matic-jagar/utils"
)

// HeimdallCurrentBal is to get current balance information using signer address and stores it in db
// Alerter will alerts whenever there is a change in balance
func HeimdallCurrentBal(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	accResp, err := scraper.HeimdallCurrentBal(ops)
	if err != nil {
		log.Printf("Error in heimdall current balance: %v", err)
		_ = db.WriteToInfluxDb(c, bp, "heimdall_current_balance", map[string]string{}, map[string]interface{}{"current_balance": "NA"})
		return
	}

	if len(accResp.Result) > 0 {
		amount := utils.ConvertToMatic(accResp.Result[0].Amount) // curent amount
		prevAmount := GetAccountBalFromDb(cfg, c)                // amount from db

		if prevAmount == "" {
			prevAmount = "0"
		}

		if prevAmount != amount {
			denom := utils.MaticDenom
			if strings.ToUpper(cfg.AlerterPreferences.BalanceChangeAlerts) == "YES" {
				_ = alerter.SendTelegramAlert(fmt.Sprintf("Heimdall Balance Change Alert : Your account balance has changed from  %s to %s", prevAmount+denom, amount+denom), cfg)
				_ = alerter.SendEmailAlert(fmt.Sprintf("Heimdall Balance Change Alert : Your account balance has changed from  %s to %s", prevAmount+denom, amount+denom), cfg)
			}
		}

		addressBalance := utils.ConvertToCommaSeparated(amount) + strings.ToUpper(accResp.Result[0].Denom)
		_ = db.WriteToInfluxDb(c, bp, "heimdall_current_balance", map[string]string{}, map[string]interface{}{"current_balance": addressBalance, "balance": amount})
		log.Printf("Address Balance: %s", addressBalance)
	}
}

// GetAccountBalFromDb returns account balance from db
func GetAccountBalFromDb(cfg *config.Config, c client.Client) string {
	var balance string
	q := client.NewQuery("SELECT last(balance) FROM heimdall_current_balance", cfg.InfluxDB.Database, "")
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
	q := client.NewQuery("SELECT last(current_balance) FROM heimdall_current_balance", cfg.InfluxDB.Database, "")
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
