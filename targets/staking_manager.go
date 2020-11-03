package targets

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetContractAddress is to get the validator share contract address and self stake
func GetContractAddress(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	subStr := GetEncodedData(ops, cfg, c, "validators(uint256)")
	if subStr == "" {
		return
	}
	valID := GetValID(cfg, c)
	// n := len(subStr) + len(valID)
	for i := 0; i < 64-len(valID); i++ {
		subStr = subStr + "0"
	}

	dataHash := subStr + valID

	if dataHash != "" {
		data := Payload{
			Jsonrpc: "2.0",
			Method:  "eth_call",
			Params: []interface{}{
				Params{
					To:   cfg.ValDetails.StakeManagerContract,
					Data: dataHash,
				},
				"latest",
			},
			ID: 1,
		}

		ops.Body = data

		resp, err := HitHTTPTarget(ops)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		if resp.Body != nil {
			var hexData EthResult
			err = json.Unmarshal(resp.Body, &hexData)
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Println("hex data of eth_call error if any..", hexData.Error)
			if hexData.Result == "" {
				return
			}

			valResp := DecodeEthCallResult(hexData.Result)
			contractAddress := "0x" + valResp[6][24:]
			stakeAmount, er := HexToBigInt(valResp[0][24:])
			if !er {
				return
			}
			amount := ConvertWeiToEth(stakeAmount) + MaticDenom

			_ = writeToInfluxDb(c, bp, "heimdall_contract_details", map[string]string{}, map[string]interface{}{"self_stake": amount, "contract_address": contractAddress})
			log.Printf("Contract Address: %s and Self Stake Amount : %s", contractAddress, amount)
		}
	}
}

// GetCommissionRate is to get the commission rate
// by calling method commissionRate() of validator share contract
func GetCommissionRate(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
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
			commissionRate, er := HexToBigInt(result.Result[2:])
			if !er {
				return
			}

			rate := ConvertWeiToEth(commissionRate) + MaticDenom

			_ = writeToInfluxDb(c, bp, "heimdall_commission_rate", map[string]string{}, map[string]interface{}{"commission_rate": rate})
			log.Printf("Contract Rate: %s", rate)
		}
	}
}

// GetValidatorRewards is to get the rewards
// by calling method validatorRewards() of validator share contract
func GetValidatorRewards(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
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
			rewards, er := HexToBigInt(result.Result[2:])
			if !er {
				return
			}
			rewradsInEth := ConvertWeiToEth(rewards) + MaticDenom

			_ = writeToInfluxDb(c, bp, "heimdall_validator_rewards", map[string]string{}, map[string]interface{}{"val_rewards": rewradsInEth})
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
func GetEncodedData(ops HTTPOptions, cfg *config.Config, c client.Client, methodSignature string) string {
	signature := methodSignature

	bytesData := []byte(signature)
	hex := EncodeToHex(bytesData)
	ops.Body.Params = append(ops.Body.Params, hex)
	ops.Body.Method = "web3_sha3"

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return ""
	}

	var hexData EthResult
	err = json.Unmarshal(resp.Body, &hexData)
	if err != nil {
		log.Printf("Error: %v", err)
		return ""
	}

	if hexData.Result == "" {
		return ""
	}

	sha3Hash := hexData.Result
	subStr := sha3Hash[:10]

	return subStr
}

// EthCall will returns the validator share contract method response
func EthCall(ops HTTPOptions, cfg *config.Config, c client.Client, dataHash string) (eth EthResult) {
	contractAddress := GetValContractAddress(cfg, c)
	data := Payload{
		Jsonrpc: "2.0",
		Method:  "eth_call",
		Params: []interface{}{
			Params{
				To:   contractAddress,
				Data: dataHash,
			},
			"latest",
		},
		ID: 1,
	}

	ops.Body = data

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return eth
	}

	var result EthResult
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return eth
	}

	return result
}
