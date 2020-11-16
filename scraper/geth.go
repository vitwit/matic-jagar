package scraper

import (
	"encoding/json"
	"log"

	"github.com/vitwit/matic-jagar/types"
)

// GetHexData will request the given endpoint of web3_sha3 or eth_call and unmarshals the data
//and returns the EthResult hex data or error if any
func GetHexData(ops types.HTTPOptions) (types.EthResult, error) {
	var hexData types.EthResult
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return hexData, err
	}
	err = json.Unmarshal(resp.Body, &hexData)
	if err != nil {
		log.Printf("Error: %v", err)
		return hexData, err
	}

	return hexData, nil
}
