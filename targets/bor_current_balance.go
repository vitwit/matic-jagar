package targets

import (
	"fmt"
	"log"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/alerter"
	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// CurrentEthBalance is to query the eth_getBalance and stores the current balance in db
// Alerter will alerts if the current balance reaches the configure threshold
func CurrentEthBalance(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	ops.Body.Params = append(ops.Body.Params, cfg.ValDetails.SignerAddress, "latest")
	balance, err := scraper.EthResult(ops)
	if err != nil {
		log.Printf("Error in GetEthBalance method: %v", err)
		return
	}

	if &balance != nil {

		bal, er := HexToBigInt(balance.Result[2:])
		if !er {
			log.Printf("Error conversion from hex to big int : %v", er)
			return
		}

		ethBalance := ConvertWeiToEth(bal)     // current balance
		prevBal := GetBorBalanceFromDB(cfg, c) // previous balance from db
		if prevBal == "" {
			prevBal = "0"
		}
		if prevBal != ethBalance {
			if strings.ToUpper(cfg.AlerterPreferences.BalanceChangeAlerts) == "YES" {
				_ = alerter.SendTelegramAlert(fmt.Sprintf("Bor Balance Change Alert : Your account balance has changed from  %s to %s", prevBal+"ETH", ethBalance+"ETH"), cfg)
				_ = alerter.SendEmailAlert(fmt.Sprintf("Bor Balance Change Alert : Your Bor account balance has changed from  %s to %s", prevBal+"ETH", ethBalance+"ETH"), cfg)
			}
		}

		balThreshold := fmt.Sprintf("%f", cfg.AlertingThresholds.EthBalanceThreshold)

		if ethBalance <= balThreshold {
			if strings.ToUpper(cfg.AlerterPreferences.EthLowBalanceAlert) == "YES" {
				_ = alerter.SendTelegramAlert(fmt.Sprintf("Eth Low Balance Alert : Your account balance has reached to your configured thershold %s", ethBalance+"ETH"), cfg)
				_ = alerter.SendEmailAlert(fmt.Sprintf("Eth Low Balance Alert : Your Bor account balance has  reached to your configured thershold %s", ethBalance+"ETH"), cfg)
			}
		}

		balWithDenom := ethBalance + "ETH"
		_ = writeToInfluxDb(c, bp, "bor_eth_balance", map[string]string{}, map[string]interface{}{"balance": balWithDenom, "amount": ethBalance})
		log.Printf("Eth Current Balance: %s", ethBalance)
	} else {
		log.Println("Got an empty response from eth rpc endpoint !")
		return
	}
}

// GetBorBalanceFRomDB returns bor validator balance from db
func GetBorBalanceFromDB(cfg *config.Config, c client.Client) string {
	var balance string
	q := client.NewQuery("SELECT last(amount) FROM bor_eth_balance", cfg.InfluxDB.Database, "")
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
