package targets

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/types"
)

type targetRunner struct{}

// NewRunner returns targetRunner
func NewRunner() *targetRunner {
	return &targetRunner{}
}

// Run to run the request
func (m targetRunner) Run(function func(ops types.HTTPOptions, cfg *config.Config, c client.Client), ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	function(ops, cfg, c)
}

// InitTargets which returns the targets
//can write all the endpoints here
func InitTargets(cfg *config.Config) *types.Targets {
	return &types.Targets{List: []types.Target{
		{
			ExecutionType: "http",
			Name:          "Net Info URL",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallRPCEndpoint + "/net_info?",
				Method:   http.MethodGet,
			},
			Func:        NetInfo,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "cmd",
			Name:          "Get Node Status",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallRPCEndpoint + "/status?",
				Method:   http.MethodGet,
			},
			Func:        Status,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Query Matic Contract",
			Name:          "Get Heimdall Current Balanace",
			HTTPOptions: types.HTTPOptions{
				Method:   http.MethodPost,
				Endpoint: cfg.Endpoints.EthRPCEndpoint,
			},
			Func:        HeimdallCurrentBal,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Node Version",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/node_info",
				Method:   http.MethodGet,
			},
			Func:        NodeVersion,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Proposals",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/gov/proposals",
				Method:   http.MethodGet,
			},
			Func:        Proposals,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Last proposed block and time",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/blocks/latest",
				Method:   http.MethodGet,
			},
			Func:        LatestProposedBlockAndTime,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Network Latest Block",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallExternalRPC + "/status?",
				Method:   http.MethodGet,
			},
			Func:        NetworkLatestBlock,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Validator Voting Power",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/staking/signer/" + cfg.ValDetails.SignerAddress,
				Method:   http.MethodGet,
			},
			Func:        ValidatorVotingPower,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Calcualte Block Time Difference",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/blocks/latest",
				Method:   http.MethodGet,
			},
			Func:        BlockTimeDifference,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Missed Blocks and send alerts",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/blocks/latest",
				Method:   http.MethodGet,
			},
			Func:        MissedBlocks,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get no of unconfirmed txns",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallRPCEndpoint + "/num_unconfirmed_txs?",
				Method:   http.MethodGet,
			},
			Func:        UnconfimedTxns,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Validator gas",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/auth/params",
				Method:   http.MethodGet,
			},
			Func:        ValidatorGas,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Validator Status Alerts",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/staking/signer/" + cfg.ValDetails.SignerAddress,
				Method:   http.MethodGet,
			},
			Func:        ValidatorStatusAlert,
			ScraperRate: cfg.Scraper.ValidatorRate,
		},
		{
			ExecutionType: "http",
			Name:          "Get total no of checkpoints",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/checkpoints/count",
				Method:   http.MethodGet,
			},
			Func:        TotalCheckPointsCount,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Latest Checkpoints",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/checkpoints/latest",
				Method:   http.MethodGet,
			},
			Func:        LatestCheckpoints,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Checkpoints Duration",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/checkpoints/params",
				Method:   http.MethodGet,
			},
			Func:        CheckpointsDuration,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get bor params",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/bor/params",
				Method:   http.MethodGet,
			},
			Func:        BorParams,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get bor latest span",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/bor/latest-span",
				Method:   http.MethodGet,
			},
			Func:        BorLatestSpan,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Current Block Height of Bor Node",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.BorRPCEndpoint,
				Method:   http.MethodPost,
				Body:     types.Payload{Jsonrpc: "2.0", Method: "eth_blockNumber", ID: 83},
			},
			Func:        CurrentBlockNumber,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Missed Blocks",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.BorRPCEndpoint,
				Method:   http.MethodPost,
				Body:     types.Payload{Jsonrpc: "2.0", Method: "bor_getSigners", ID: 1},
			},
			Func:        BorMissedBlocks,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Eth Balance",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.EthRPCEndpoint,
				Method:   http.MethodPost,
				Body:     types.Payload{Jsonrpc: "2.0", Method: "eth_getBalance", ID: 1},
			},
			Func:        CurrentEthBalance,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Bor Current Proposer",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.BorRPCEndpoint,
				Method:   http.MethodPost,
				Body:     types.Payload{Jsonrpc: "2.0", Method: "bor_getCurrentProposer", ID: 1},
			},
			Func:        BorCurrentProposer,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Telegram command",
			Name:          "command based alerts",
			Func:          TelegramAlerting,
			ScraperRate:   cfg.Scraper.CommandsRate,
		},
		{
			ExecutionType: "http",
			Name:          "Get and Store Validator Share Contract Address",
			HTTPOptions: types.HTTPOptions{
				Method:   http.MethodPost,
				Endpoint: cfg.Endpoints.EthRPCEndpoint,
			},
			Func:        ContractAddress,
			ScraperRate: cfg.Scraper.ContractRate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Commission Rate",
			HTTPOptions: types.HTTPOptions{
				Method:   http.MethodPost,
				Endpoint: cfg.Endpoints.EthRPCEndpoint,
			},
			Func:        GetCommissionRate,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Validator Rewards",
			HTTPOptions: types.HTTPOptions{
				Method:   http.MethodPost,
				Endpoint: cfg.Endpoints.EthRPCEndpoint,
			},
			Func:        GetValidatorRewards,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Bor Pending Transactions",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.BorRPCEndpoint,
				Method:   http.MethodPost,
				Body:     types.Payload{Jsonrpc: "2.0", Method: "eth_pendingTransactions", ID: 64},
			},
			Func:        BorPendingTransactions,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Check weather validator is part of block producer",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/bor/span/",
				Method:   http.MethodGet,
			},
			Func:        BlockProducer,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get proposed checkpoints",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/checkpoints/",
				Method:   http.MethodGet,
			},
			Func:        ProposedCheckpoints,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Network Height of Bor",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.BorExternalRPC,
				Method:   http.MethodPost,
				Body:     types.Payload{Jsonrpc: "2.0", Method: "eth_blockNumber", ID: 83},
			},
			Func:        BorNetworkHeight,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Validator Caught UP",
			HTTPOptions: types.HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/syncing",
				Method:   http.MethodGet,
			},
			Func:        ValidatorCaughtUp,
			ScraperRate: cfg.Scraper.Rate,
		},
	}}
}

func addQueryParameters(req *http.Request, queryParams types.QueryParams) {
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	req.URL.RawQuery = params.Encode()
}

//newHTTPRequest to make a new http request
func newHTTPRequest(ops types.HTTPOptions) (*http.Request, error) {
	// make new request
	payloadBytes, _ := json.Marshal(ops.Body)
	req, err := http.NewRequest(ops.Method, ops.Endpoint, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	// Add any query parameters to the URL.
	if len(ops.QueryParams) != 0 {
		addQueryParameters(req, ops.QueryParams)
	}

	return req, nil
}

func makeResponse(res *http.Response) (*types.PingResp, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &types.PingResp{}, err
	}

	response := &types.PingResp{
		StatusCode: res.StatusCode,
		Body:       body,
	}
	_ = res.Body.Close()
	return response, nil
}

// HitHTTPTarget to hit the target and get response
func HitHTTPTarget(ops types.HTTPOptions) (*types.PingResp, error) {
	req, err := newHTTPRequest(ops)
	if err != nil {
		return nil, err
	}

	httpcli := http.Client{Timeout: time.Duration(5 * time.Second)}
	resp, err := httpcli.Do(req)
	if err != nil {
		return nil, err
	}

	res, err := makeResponse(resp)
	if err != nil {
		return nil, err
	}

	return res, nil
}
