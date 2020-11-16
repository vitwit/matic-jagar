package scraper

import (
	"encoding/json"
	"log"

	"github.com/vitwit/matic-jagar/types"
)

// HeimdallCurrentBal will request the given endpoint and unmarshals the data
// and returns the account bal response or error if any
func HeimdallCurrentBal(ops types.HTTPOptions) (types.AccountBalResp, error) {
	var accResp types.AccountBalResp
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return accResp, err
	}

	err = json.Unmarshal(resp.Body, &accResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return accResp, err
	}

	return accResp, nil
}

// AuthParams will request the given endpoint and unmarshals the data
// and returns the auth params such as max tax gas etc or error if any
func AuthParams(ops types.HTTPOptions) (types.AuthParams, error) {
	var authParam types.AuthParams
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return authParam, err
	}

	err = json.Unmarshal(resp.Body, &authParam)
	if err != nil {
		log.Printf("Error: %v", err)
		return authParam, err
	}

	return authParam, nil
}

// LatestBlock will request the given endpoint and unmarshals the data
// and returns the latest block details or error if any
func LatestBlock(ops types.HTTPOptions) (types.LatestBlock, error) {
	var result types.LatestBlock
	currResp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	err = json.Unmarshal(currResp.Body, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	return result, nil
}

// GetTotalCheckPoints will request the given endpoint and unmarshals the data
// and returns the latest block details or error if any
func GetTotalCheckPoints(ops types.HTTPOptions) (types.TotalCheckpoints, error) {
	var result types.TotalCheckpoints
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

// GetLatestCheckpoints will request the given endpoint and unmarshals the data
// and returns the latest checkpoint details or error if any
func GetLatestCheckpoints(ops types.HTTPOptions) (types.LatestCheckpoints, error) {
	var lcp types.LatestCheckpoints
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return lcp, err
	}

	err = json.Unmarshal(resp.Body, &lcp)
	if err != nil {
		log.Printf("Error: %v", err)
		return lcp, err
	}

	return lcp, nil
}

// GetCheckpointsDuration will request the given endpoint and unmarshals the data
// and returns the latest checkpoint details or error if any
func GetCheckpointsDuration(ops types.HTTPOptions) (types.CheckpointsDuration, error) {
	var cpd types.CheckpointsDuration
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return cpd, err
	}

	err = json.Unmarshal(resp.Body, &cpd)
	if err != nil {
		log.Printf("Error: %v", err)
		return cpd, err
	}

	return cpd, nil
}

// GetProposedCheckpoints will request the given endpoint and unmarshals the data
// and returns the proposed checkpoint details or error if any
func GetProposedCheckpoints(ops types.HTTPOptions) (types.ProposedCheckpoints, error) {
	var proposedCP types.ProposedCheckpoints
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return proposedCP, err
	}

	err = json.Unmarshal(resp.Body, &proposedCP)
	if err != nil {
		log.Printf("Error: %v", err)
		return proposedCP, err
	}
	return proposedCP, nil
}

// GetNetInfo will request the given endpoint and unmarshals the data
// and returns the net info details such as no.of peers and addresses or error if any
func GetNetInfo(ops types.HTTPOptions) (types.NetInfo, error) {
	var result types.NetInfo
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error getting net info: %v", err)
		return result, err
	}

	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		log.Printf("Error: %v", err)
		return result, err
	}

	return result, nil
}

// GetStatus will request the given endpoint and unmarshals the data
// and returns the status and addresses or error if any
func GetStatus(ops types.HTTPOptions) (types.Status, error) {
	var result types.Status
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

// GetCaughtUpStatus will request the given endpoint and unmarshals the data
// and returns the validator caughtup status or error if any
func GetCaughtUpStatus(ops types.HTTPOptions) (types.Caughtup, error) {
	var sync types.Caughtup
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return sync, err
	}

	err = json.Unmarshal(resp.Body, &sync)
	if err != nil {
		log.Printf("Error: %v", err)
		return sync, err
	}

	return sync, nil
}

// GetVersion will request the given endpoint and unmarshals the data
// and returns the application info such as version or error if any
func GetVersion(ops types.HTTPOptions) (types.ApplicationInfo, error) {
	var applicationInfo types.ApplicationInfo
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return applicationInfo, err
	}

	err = json.Unmarshal(resp.Body, &applicationInfo)
	if err != nil {
		log.Printf("Error: %v", err)
		return applicationInfo, err
	}

	return applicationInfo, nil
}

// GetProposals will request the given endpoint and unmarshals the data
// and returns the proposals or error if any
func GetProposals(ops types.HTTPOptions) (types.Proposals, error) {
	var p types.Proposals
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return p, err
	}

	err = json.Unmarshal(resp.Body, &p)
	if err != nil {
		log.Printf("Error: %v", err)
		return p, err
	}

	return p, nil
}

// GetProposalVoters will request the given endpoint and unmarshals the data
// and returns the voters of proposal or error if any
func GetProposalVoters(ops types.HTTPOptions) (types.ProposalVoters, error) {
	var voters types.ProposalVoters
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return voters, err
	}

	err = json.Unmarshal(resp.Body, &voters)
	if err != nil {
		log.Printf("Error: %v", err)
		return voters, err
	}

	return voters, nil
}

// GetProposalDepositors will request the given endpoint and unmarshals the data
// and returns the deposits of a proposals or error if any
func GetProposalDepositors(ops types.HTTPOptions) (types.Depositors, error) {
	var depositors types.Depositors
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return depositors, err
	}

	err = json.Unmarshal(resp.Body, &depositors)
	if err != nil {
		log.Printf("Error: %v", err)
		return depositors, err
	}

	return depositors, nil
}

// GetUnconfirmedTxs will request the given endpoint and unmarshals the data
// and returns the unconfirmed transactions or error if any
func GetUnconfirmedTxs(ops types.HTTPOptions) (types.UnconfirmedTxns, error) {
	var txs types.UnconfirmedTxns
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

// GetValStatus will request the given endpoint and unmarshals the data
// and returns thevalidator status whether it's voting/jailed or error if any
func GetValStatus(ops types.HTTPOptions) (types.ValStatusResp, error) {
	var result types.ValStatusResp
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
