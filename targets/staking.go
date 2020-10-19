package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

func GetEncodedData(ops HTTPOptions, cfg *config.Config, c client.Client) {
	signature := "validators(uint256)"

	bytesData := []byte(signature)
	hex := EncodeToHex(bytesData)
	ops.Body.Params = append(ops.Body.Params, hex)

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

	sha3Hash := hexData.Result
	subStr := sha3Hash[:10]

	valID := GetValID(cfg, c)

	for i := 0; i < 64-len(valID); i++ {
		subStr = subStr + "0"
	}

	dataHash := subStr + valID

	log.Println("hex...", sha3Hash)
	log.Println("prefix..", dataHash)
}

func GetContractAddress(ops HTTPOptions, cfg *config.Config, c client.Client) {

	ops.Body.Params = append(ops.Body.Params, "")
}
