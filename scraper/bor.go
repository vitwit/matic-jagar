package scraper

import (
	"encoding/json"
	"log"

	"github.com/vitwit/matic-jagar/types"
)

func EthResult(ops types.HTTPOptions) (types.EthResult, error) {
	var result types.EthResult

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	return result, nil
}

func EthBlockNumber(ops types.HTTPOptions) (types.BorValHeight, error) {
	var result types.BorValHeight
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	return result, nil
}

func BorLatestSpan(ops types.HTTPOptions) (types.BorLatestSpan, error) {
	var latestSpan types.BorLatestSpan
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return latestSpan, err
	}

	err = json.Unmarshal(resp.Body, &latestSpan)
	if err != nil {
		log.Printf("Error: %v", err)
		return latestSpan, err
	}

	return latestSpan, nil
}

func BorSignersRes(ops types.HTTPOptions) (types.BorSignersRes, error) {
	var signers types.BorSignersRes
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return signers, err
	}

	err = json.Unmarshal(resp.Body, &signers)
	if err != nil {
		log.Printf("Error: %v", err)
		return signers, err
	}

	return signers, nil
}

func BorValidatorHeight(ops types.HTTPOptions) (types.BorValHeight, error) {
	var result types.BorValHeight
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}
	return result, nil
}

func BorParams(ops types.HTTPOptions) (types.BorParams, error) {
	var params types.BorParams
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return params, err
	}

	err = json.Unmarshal(resp.Body, &params)
	if err != nil {
		log.Printf("Error: %v", err)
		return params, err
	}

	return params, nil
}

func BorPendingTransactions(ops types.HTTPOptions) (types.EthPendingTransactions, error) {
	var txs types.EthPendingTransactions
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return txs, err
	}

	err = json.Unmarshal(resp.Body, &txs)
	if err != nil {
		log.Printf("Error: %v", err)
		return txs, err
	}

	return txs, nil
}
