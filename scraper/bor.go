package scraper

import (
	"encoding/json"
	"log"

	"github.com/vitwit/matic-jagar/types"
)

// EthResult will request the given endpoint and unmarshals the data
// Returns the EthResult data or error if any
func EthResult(ops types.HTTPOptions) (types.EthResult, error) {
	var result types.EthResult

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while getting response of eth result: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error while unmarshelling eth response: %v", err)
		return result, err
	}

	return result, nil
}

// EthBlockNumber will request the given endpoint and unmarshals the data
// Returns the eth block height or error if any
func EthBlockNumber(ops types.HTTPOptions) (types.BorValHeight, error) {
	var result types.BorValHeight
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while getting response of bor block number: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error while unmarshelling response of bor block: %v", err)
		return result, err
	}

	return result, nil
}

// BorLatestSpan will request the given endpoint and unmarshals the data
// Returns the bor latest span details or error if any
func BorLatestSpan(ops types.HTTPOptions) (types.BorLatestSpan, error) {
	var latestSpan types.BorLatestSpan
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while getting bor latest span: %v", err)
		return latestSpan, err
	}

	err = json.Unmarshal(resp.Body, &latestSpan)
	if err != nil {
		log.Printf("Error while unmarshelling bor latest span: %v", err)
		return latestSpan, err
	}

	return latestSpan, nil
}

// BorSignersRes will request the given endpoint and unmarshals the data
// Returns the bor signer details of given block height or error if any
func BorSignersRes(ops types.HTTPOptions) (types.BorSignersRes, error) {
	var signers types.BorSignersRes
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while getting bor signers: %v", err)
		return signers, err
	}

	err = json.Unmarshal(resp.Body, &signers)
	if err != nil {
		log.Printf("Error while unmarshelling bor signers: %v", err)
		return signers, err
	}

	return signers, nil
}

// BorValidatorHeight will request the given endpoint and unmarshals the data
// Returns the bor validator height or error if any
func BorValidatorHeight(ops types.HTTPOptions) (types.BorValHeight, error) {
	var result types.BorValHeight
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while getting bor validtaor height: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error while unmarshelling bor validator height: %v", err)
		return result, err
	}
	return result, nil
}

// BorParams will request the given endpoint and unmarshals the data
// Returns the bor params such as span duration producer count etc or error if any
func BorParams(ops types.HTTPOptions) (types.BorParams, error) {
	var params types.BorParams
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while gettng bor params: %v", err)
		return params, err
	}

	err = json.Unmarshal(resp.Body, &params)
	if err != nil {
		log.Printf("Error while unmarshelling bor params: %v", err)
		return params, err
	}

	return params, nil
}

// BorPendingTransactions will request the given endpoint and unmarshals the data
// Returns the eth pending transaction details or error if any
func BorPendingTransactions(ops types.HTTPOptions) (types.EthPendingTransactions, error) {
	var txs types.EthPendingTransactions
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while getting bor pending txs: %v", err)
		return txs, err
	}

	err = json.Unmarshal(resp.Body, &txs)
	if err != nil {
		log.Printf("Error while unmarshelling bor pedning txs: %v", err)
		return txs, err
	}

	return txs, nil
}

// GetSpanProducers will request the given endpoint and unmarshals the data
// Returns span producer details or error if any
func GetSpanProducers(ops types.HTTPOptions) (types.BorSpanProducers, error) {
	var spanProducers types.BorSpanProducers
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error while getting span producers: %v", err)
		return spanProducers, err
	}

	err = json.Unmarshal(resp.Body, &spanProducers)
	if err != nil {
		log.Printf("Error while unmarshelling span producers: %v", err)
		return spanProducers, err
	}
	return spanProducers, nil
}
