package scraper

import (
	"encoding/json"
	"log"

	"github.com/vitwit/matic-jagar/types"
)

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
