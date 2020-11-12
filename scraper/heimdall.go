package scraper

import (
	"encoding/json"
	"log"

	"github.com/vitwit/matic-jagar/types"
)

func HeimdallCurrentBal(ops types.HTTPOptions) (types.AccountBalResp, error) {
	var accResp types.AccountBalResp
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		// _ = writeToInfluxDb(c, bp, "heimdall_current_balance", map[string]string{}, map[string]interface{}{"current_balance": "NA"})
		return accResp, err
	}

	err = json.Unmarshal(resp.Body, &accResp)
	if err != nil {
		log.Printf("Error: %v", err)
		// _ = writeToInfluxDb(c, bp, "heimdall_current_balance", map[string]string{}, map[string]interface{}{"current_balance": "NA"})
		return accResp, err
	}

	return accResp, nil
}

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
