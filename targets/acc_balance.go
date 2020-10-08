package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/big"
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

func BorLatestBalance(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		var balance BorResult
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

		_ = writeToInfluxDb(c, bp, "matic_bor_balance", map[string]string{}, map[string]interface{}{"current_balance": bal})
		log.Printf("Bor Current Balance: %d", bal)
	}

}

func ConvertWeiToEth(num *big.Int) {
	wei := num.String()

	f, _ := strconv.ParseFloat(wei, 64)
	eth := f / math.Pow(10, 18)
	ether := fmt.Sprintf("%.8f", eth)

	log.Println("eth..", ether)
}

func HexToBigInt(hex string) (*big.Int, bool) {
	n := new(big.Int)
	n2, err := n.SetString(hex[2:], 16)

	return n2, err
}