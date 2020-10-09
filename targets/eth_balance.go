package targets

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetEthBalance to get eth balance
func GetEthBalance(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	ops.Body.Params = append(ops.Body.Params, cfg.SignerAddress, "latest")
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

		bal, er := HexToBigInt(balance.Result)
		if !er {
			log.Printf("Error conversion from hex to big int : %v", er)
			return
		}

		ethBalance := ConvertWeiToEth(bal)

		if ethBalance < cfg.EthBalanceThreshold {
			_ = SendTelegramAlert(fmt.Sprintf("Your Eth current balance has dropped below %s", cfg.EthBalanceThreshold), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Your Eth current balance has dropped below %s", cfg.EthBalanceThreshold), cfg)
		}

		_ = writeToInfluxDb(c, bp, "matic_eth_balance", map[string]string{}, map[string]interface{}{"balance": ethBalance + "ETH"})
		log.Printf("Eth Current Balance: %s", ethBalance)
	}

}
