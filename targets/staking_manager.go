package targets

import (
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
	"github.com/vitwit/matic-jagar/utils"
)

// ContractAddress is to get the validator share contract address, self stake and stores it in db
func ContractAddress(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	subStr := GetEncodedData(ops, cfg, c, "validators(uint256)")
	if subStr == "" {
		return
	}
	valID := GetValID(cfg, c)

	ID, _ := strconv.ParseInt(valID, 10, 64)
	hexNum := strconv.FormatInt(ID, 16)
	log.Printf("hex number of validator ID", hexNum)
	// n := len(subStr) + len(valID)
	for i := 0; i < 64-len(hexNum); i++ {
		subStr = subStr + "0"
	}

	dataHash := subStr + hexNum

	if dataHash != "" {
		data := types.Payload{
			Jsonrpc: "2.0",
			Method:  "eth_call",
			Params: []interface{}{
				types.Params{
					To:   cfg.ValDetails.StakeManagerContract,
					Data: dataHash,
				},
				"latest",
			},
			ID: 1,
		}

		ops.Body = data
		hexData, err := scraper.GetHexData(ops)
		if err != nil {
			log.Printf("Error in contract address: %v", err)
			return
		}

		if hexData.Result != "" {

			valResp := utils.DecodeEthCallResult(hexData.Result)
			contractAddress := "0x" + valResp[6][24:]
			stakeAmount, er := utils.HexToBigInt(valResp[0][24:])
			if !er {
				return
			}

			amount := utils.FixSelfStakeDecimals(stakeAmount) + utils.MaticDenom

			_ = db.WriteToInfluxDb(c, bp, "heimdall_contract_details", map[string]string{}, map[string]interface{}{"self_stake": amount, "contract_address": contractAddress})
			log.Printf("Contract Address: %s and Self Stake Amount : %s", contractAddress, amount)
		} else {
			log.Println("Got an empty response from eth rpc endpoint ! ")
			return
		}
	}
}

// GetCommissionRate is to get the commission rate
// by calling method commissionRate() of validator share contract
func GetCommissionRate(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	subStr := GetEncodedData(ops, cfg, c, "commissionRate()")
	if subStr == "" {
		return
	}
	n := len(subStr)
	for i := 0; i < 66-n; i++ {
		subStr = subStr + "0"
	}
	dataHash := subStr
	if dataHash != "" {
		result := EthCall(ops, cfg, c, dataHash)
		if result.Result != "" {
			commissionRate, er := utils.HexToBigInt(result.Result[2:])
			if !er {
				return
			}

			var fee float64
			f, err := strconv.ParseFloat(commissionRate.String(), 64)
			if err != nil {
				fee = 0
			} else {
				fee = f * 100
			}

			_ = db.WriteToInfluxDb(c, bp, "heimdall_commission_rate", map[string]string{}, map[string]interface{}{"commission_rate": fee})
			log.Printf("Contract Rate: %f", fee)
		}
	}
}

// GetValidatorRewards is to get the rewards
// by calling method validatorRewards() of validator share contract
func GetValidatorRewards(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	subStr := GetEncodedData(ops, cfg, c, "validatorRewards()")
	if subStr == "" {
		return
	}
	n := len(subStr)
	for i := 0; i < 66-n; i++ {
		subStr = subStr + "0"
	}
	dataHash := subStr
	if dataHash != "" {
		result := EthCall(ops, cfg, c, dataHash)
		if result.Result != "" {
			rewards, er := utils.HexToBigInt(result.Result[2:])
			if !er {
				return
			}
			rewradsInEth := utils.ConvertWeiToEth(rewards) + utils.MaticDenom

			_ = db.WriteToInfluxDb(c, bp, "heimdall_validator_rewards", map[string]string{}, map[string]interface{}{"val_rewards": rewradsInEth})
			log.Printf("Validator Rewards: %s", rewradsInEth)
		}
	}
}

// GetValContractAddress returns validator share contract address from db
func GetValContractAddress(cfg *config.Config, c client.Client) string {
	var address string
	q := client.NewQuery("SELECT last(contract_address) FROM heimdall_contract_details", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						heightValue := r.Series[0].Values[0][idx]
						address = fmt.Sprintf("%v", heightValue)
						break
					}
				}
			}
		}
	}
	return address
}

// GetEncodedData returns the sha3Hash
//which will be calculated for the partcular method signature and by making a call of web3_sha3
func GetEncodedData(ops types.HTTPOptions, cfg *config.Config, c client.Client, methodSignature string) string {
	signature := methodSignature

	bytesData := []byte(signature)
	hex := utils.EncodeToHex(bytesData)
	ops.Body.Params = append(ops.Body.Params, hex)
	ops.Body.Method = "web3_sha3"

	hexData, err := scraper.GetHexData(ops)
	if err != nil {
		log.Printf("Error in get encoded data: %v", err)
		return ""
	}

	log.Printf("hex data : %v of signature : %s", hexData, signature)

	if hexData.Result == "" {
		log.Println("Response of web3_sha3 is empty!")
		return ""
	}

	sha3Hash := hexData.Result
	subStr := sha3Hash[:10]

	return subStr
}

// EthCall will returns the validator share contract method response
func EthCall(ops types.HTTPOptions, cfg *config.Config, c client.Client, dataHash string) (eth types.EthResult) {
	contractAddress := GetValContractAddress(cfg, c)

	data := types.Payload{
		Jsonrpc: "2.0",
		Method:  "eth_call",
		Params: []interface{}{
			types.Params{
				To:   contractAddress,
				Data: dataHash,
			},
			"latest",
		},
		ID: 1,
	}

	ops.Body = data

	result, err := scraper.GetHexData(ops)
	if err != nil {
		log.Printf("Error in eth call: %v", err)
		return
	}

	return result
}
