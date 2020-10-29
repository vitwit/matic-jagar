package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

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

	sha3Hash := hexData.Result
	subStr := sha3Hash[:10]

	valID := GetValID(cfg, c)

	for i := 0; i < 64-len(valID); i++ {
		subStr = subStr + "0"
	}

	dataHash := subStr + valID

	log.Println("hex...", sha3Hash)
	log.Println("prefix..", dataHash)

	return dataHash
}

func GetContractAddress(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	dataHash := GetEncodedData(ops, cfg, c, "validators(uint256)")

	if dataHash != "" {
		data := Payload{
			Jsonrpc: "2.0",
			Method:  "eth_call",
			Params: []interface{}{
				Params{
					To:   cfg.ValDetails.ContractAddress,
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

		var hexData EthResult
		err = json.Unmarshal(resp.Body, &hexData)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}

		log.Println("hex data of eth_call error if any..", hexData.Error)

		valResp := DecodeEthCallResult(hexData.Result)

		contractAddress := "0x" + valResp[6][24:]

		stakeAmount, _ := strconv.ParseInt(valResp[0], 16, 64)
		value := ConvertValueToEth(stakeAmount)
		amount := fmt.Sprintf("%.6f", value) + MaticDenom

		_ = writeToInfluxDb(c, bp, "heimdall_contract_details", map[string]string{}, map[string]interface{}{"self_stake": amount, "contract_address": contractAddress})
		log.Printf("Contract Address: %s and Self Stake Amount : %d", contractAddress, amount)

	}
}

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

func GetCommissionRate(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	dataHash := GetEncodedData(ops, cfg, c, "commissionRate()")
	if dataHash != "" {
		result := EthCall(ops, cfg, c, dataHash)
		if result.Result != "" {
			commissionRate, _ := strconv.ParseInt(result.Result[2:], 16, 64)
			value := ConvertValueToEth(commissionRate)
			rate := fmt.Sprintf("%.2f", value)

			_ = writeToInfluxDb(c, bp, "heimdall_commission_rate", map[string]string{}, map[string]interface{}{"commission_rate": rate})
			log.Printf("Contract Rate: %d", rate)
		}
	}

}

func GetValidatorRewards(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	dataHash := GetEncodedData(ops, cfg, c, "validatorRewards()")
	if dataHash != "" {

		result := EthCall(ops, cfg, c, dataHash)

		if result.Result != "" {
			rewards, _ := strconv.ParseInt(result.Result[2:], 16, 64)

			rewardsEth := ConvertValueToEth(rewards)
			ether := fmt.Sprintf("%.8f", rewardsEth) + MaticDenom

			_ = writeToInfluxDb(c, bp, "heimdall_validator_rewards", map[string]string{}, map[string]interface{}{"val_rewards": ether})
			log.Printf("Validator Rewards: %s", ether)
		}
	}
}

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
