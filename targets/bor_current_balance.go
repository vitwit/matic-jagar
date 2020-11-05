package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetEthBalance to get eth balance
func GetEthBalance(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	ops.Body.Params = append(ops.Body.Params, cfg.ValDetails.SignerAddress, "latest")
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	if resp != nil {
		var balance EthResult
		err = json.Unmarshal(resp.Body, &balance)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

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
			if strings.ToUpper(cfg.ChooseAlerts.BalanceChangeAlerts) == "YES" {
				_ = SendTelegramAlert(fmt.Sprintf("Bor Balance Change Alert : Your account balance has changed from  %s to %s", prevBal+"ETH", ethBalance), cfg)
				_ = SendEmailAlert(fmt.Sprintf("Bor Balance Change Alert : Your Bor account balance has changed from  %s to %s", prevBal+"ETH", ethBalance), cfg)
			}
		}

		balWithDenom := ethBalance + "ETH"
		_ = writeToInfluxDb(c, bp, "matic_eth_balance", map[string]string{}, map[string]interface{}{"balance": balWithDenom, "amount": ethBalance})
		log.Printf("Eth Current Balance: %s", ethBalance)
	}

}

// GetBorBalanceFRomDB returns bor validator balance from db
func GetBorBalanceFromDB(cfg *config.Config, c client.Client) string {
	var balance string
	q := client.NewQuery("SELECT last(amount) FROM matic_eth_balance", cfg.InfluxDB.Database, "")
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
