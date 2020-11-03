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
)

type targetRunner struct{}

// NewRunner returns targetRunner
func NewRunner() *targetRunner {
	return &targetRunner{}
}

// Run to run the request
func (m targetRunner) Run(function func(ops HTTPOptions, cfg *config.Config, c client.Client), ops HTTPOptions, cfg *config.Config, c client.Client) {
	function(ops, cfg, c)
}

// InitTargets which returns the targets
//can write all the endpoints here
func InitTargets(cfg *config.Config) *Targets {
	return &Targets{List: []Target{
		{
			ExecutionType: "http",
			Name:          "Net Info URL",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallRPCEndpoint + "/net_info?",
				Method:   http.MethodGet,
			},
			Func:        GetNetInfo,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "cmd",
			Name:          "Get Node Status",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallRPCEndpoint + "/status?",
				Method:   http.MethodGet,
			},
			Func:        GetNodeStatus,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http", // Confirmation about alerting
			Name:          "Get Heimdall Current Balanace",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/bank/balances/" + cfg.ValDetails.SignerAddress,
				Method:   http.MethodGet,
			},
			Func:        GetHeimdallCurrentBal,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Node Version",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/node_info",
				Method:   http.MethodGet,
			},
			Func:        NodeVersion,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Proposals",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/gov/proposals",
				Method:   http.MethodGet,
			},
			Func:        GetProposals,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Last proposed block and time",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/blocks/latest",
				Method:   http.MethodGet,
			},
			Func:        GetLatestProposedBlockAndTime,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Network Latest Block",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallExternalRPC + "/status?",
				Method:   http.MethodGet,
			},
			Func:        GetNetworkLatestBlock,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Validator Voting Power",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/staking/signer/" + cfg.ValDetails.SignerAddress,
				Method:   http.MethodGet,
			},
			Func:        GetValidatorVotingPower,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Block Time Difference",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/blocks/latest",
				Method:   http.MethodGet,
			},
			Func:        GetBlockTimeDifference,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Missed Blocks",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/blocks/latest",
				Method:   http.MethodGet,
			},
			Func:        GetMissedBlocks,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get no of unconfirmed txns",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallRPCEndpoint + "/num_unconfirmed_txs?",
				Method:   http.MethodGet,
			},
			Func:        GetUnconfimedTxns,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Validator fee and gas",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/auth/params", //take confirmation about validator fee
				Method:   http.MethodGet,
			},
			Func:        GetValidatorFeeAndGas,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Validator status",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/staking/signer/" + cfg.ValDetails.SignerAddress,
				Method:   http.MethodGet,
			},
			Func:        ValidatorStatusAlert,
			ScraperRate: cfg.Scraper.ValidatorRate,
		},
		{
			ExecutionType: "http",
			Name:          "Get total no of checkpoints",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/checkpoints/count",
				Method:   http.MethodGet,
			},
			Func:        GetTotalCheckPointsCount,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Latest Checkpoints",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/checkpoints/latest",
				Method:   http.MethodGet,
			},
			Func:        GetLatestCheckpoints,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Checkpoints Duration",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/checkpoints/params",
				Method:   http.MethodGet,
			},
			Func:        GetCheckpointsDuration,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get bor params",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/bor/params",
				Method:   http.MethodGet,
			},
			Func:        GetBorParams,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get bor latest span",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/bor/latest-span",
				Method:   http.MethodGet,
			},
			Func:        GetBorLatestSpan,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Current Block Height of Bor Node",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.BorRPCEndpoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "eth_blockNumber", ID: 83},
			},
			Func:        BorCurrentHeight,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Missed Blocks",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.BorRPCEndpoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "bor_getSigners", ID: 1},
			},
			Func:        GetBorMissedBlocks,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Eth Balance",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.EthRPCEndpoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "eth_getBalance", ID: 1},
			},
			Func:        GetEthBalance,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Bor Current Proposer",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.BorRPCEndpoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "bor_getCurrentProposer", ID: 1},
			},
			Func:        GetBorCurrentProposer,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "Telegram command",
			Name:          "command based alerts",
			Func:          TelegramAlerting,
			ScraperRate:   "2s",
		},
		{
			ExecutionType: "curl",
			Name:          "Get Contract Address",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.EthRPCEndpoint,
			},
			Func:        GetContractAddress,
			ScraperRate: cfg.Scraper.ContractRate,
		},
		{
			ExecutionType: "curl",
			Name:          "Get Commission Rate",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.EthRPCEndpoint,
			},
			Func:        GetCommissionRate,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl",
			Name:          "Get Validator Rewards",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.EthRPCEndpoint,
			},
			Func:        GetValidatorRewards,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Bor Pending Transactions",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.BorRPCEndpoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "eth_pendingTransactions", ID: 64},
			},
			Func:        GetBorPendingTransactions,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Check weather validator is part of block producer",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/bor/span/",
				Method:   http.MethodGet,
			},
			Func:        GetBlockProducer,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get proposed checkpoints",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/checkpoints/",
				Method:   http.MethodGet,
			},
			Func:        GetProposedCheckpoints,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Network Height of Bor",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.BorExternalRPC,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "eth_blockNumber", ID: 83},
			},
			Func:        BorNetworkHeight,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Validator Caught UP",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.Endpoints.HeimdallLCDEndpoint + "/syncing",
				Method:   http.MethodGet,
			},
			Func:        ValidatorCaughtUp,
			ScraperRate: cfg.Scraper.Rate,
		},
	}}
}

func addQueryParameters(req *http.Request, queryParams QueryParams) {
	params := url.Values{}
	for key, value := range queryParams {
		params.Add(key, value)
	}
	req.URL.RawQuery = params.Encode()
}

//newHTTPRequest to make a new http request
func newHTTPRequest(ops HTTPOptions) (*http.Request, error) {
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

func makeResponse(res *http.Response) (*PingResp, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return &PingResp{}, err
	}

	response := &PingResp{
		StatusCode: res.StatusCode,
		Body:       body,
	}
	_ = res.Body.Close()
	return response, nil
}

// HitHTTPTarget to hit the target and get response
func HitHTTPTarget(ops HTTPOptions) (*PingResp, error) {
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
