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
				Endpoint: cfg.ValidatorRPCEndpoint + "/net_info?",
				Method:   http.MethodGet,
			},
			Func:        GetNetInfo,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "cmd",
			Name:          "Get Node Status",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.ValidatorRPCEndpoint + "/status?",
				Method:   http.MethodGet,
			},
			Func:        GetNodeStatus,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http", // Confirmation about alerting
			Name:          "Get Heimdall Current Balanace",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/bank/balances/" + cfg.SignerAddress,
				Method:   http.MethodGet,
			},
			Func:        GetHeimdallCurrentBal,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Node Version",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/node_info",
				Method:   http.MethodGet,
			},
			Func:        NodeVersion,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Proposals",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/gov/proposals",
				Method:   http.MethodGet,
			},
			Func:        GetProposals,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Last proposed block and time",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/blocks/latest",
				Method:   http.MethodGet,
			},
			Func:        GetLatestProposedBlockAndTime,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Network Latest Block",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.ExternalRPC + "/status?",
				Method:   http.MethodGet,
			},
			Func:        GetNetworkLatestBlock,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Validator Voting Power",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/staking/signer/" + cfg.SignerAddress,
				Method:   http.MethodGet,
			},
			Func:        GetValidatorVotingPower,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Block Time Difference",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.ValidatorRPCEndpoint + "/block",
				Method:   http.MethodGet,
			},
			Func:        GetBlockTimeDifference,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Current Block Height",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.ExternalRPC + "/status",
				Method:   http.MethodGet,
			},
			Func:        GetMissedBlocks,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get no of unconfirmed txns",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.ValidatorRPCEndpoint + "/num_unconfirmed_txs?",
				Method:   http.MethodGet,
			},
			Func:        GetUnconfimedTxns,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Validator fee and gas",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/auth/params",
				Method:   http.MethodGet,
			},
			Func:        GetValidatorFeeAndGas,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Validator status",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/staking/signer/" + cfg.SignerAddress,
				Method:   http.MethodGet,
			},
			Func:        ValidatorStatusAlert,
			ScraperRate: cfg.Scraper.ValidatorRate,
		},
		{
			ExecutionType: "http",
			Name:          "Get bor params",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/bor/params",
				Method:   http.MethodGet,
			},
			Func:        GetBorParams,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get total no of checkpoints",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/checkpoints/count",
				Method:   http.MethodGet,
			},
			Func:        GetTotalCheckPointsCount,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Latest Checkpoints",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/checkpoints/latest",
				Method:   http.MethodGet,
			},
			Func:        GetLatestCheckpoints,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "http",
			Name:          "Get Checkpoints Duration",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.LCDEndpoint + "/checkpoints/params",
				Method:   http.MethodGet,
			},
			Func:        GetCheckpointsDuration,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Peers count of bor",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.BorEndPoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "net_peerCount", ID: 74},
			},
			Func:        BorPeersCount,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Current Block Height",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.BorEndPoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "eth_blockNumber", ID: 83},
			},
			Func:        BorCurrentHeight,
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Bor Balance",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.BorEndPoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "eth_getBalance", ID: 1},
			},
			Func:        BorLatestBalance, // clarity on wei to eth conversion
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Eth Syncing",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.BorEndPoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "eth_syncing", ID: 1},
			},
			Func:        BorEthSyncing, // Clarity about metric
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Bor Net Listening",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.BorEndPoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "net_listening", ID: 67},
			},
			Func:        BorNetListening, // Clarity about the metric and alerting
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Bor Eth Mining",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.BorEndPoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "eth_mining", ID: 71},
			},
			Func:        BorEthMining, // Clarity about the metric and alerting
			ScraperRate: cfg.Scraper.Rate,
		},
		{
			ExecutionType: "curl cmd",
			Name:          "Get Eth Balance",
			HTTPOptions: HTTPOptions{
				Endpoint: cfg.EthRPCEndPoint,
				Method:   http.MethodPost,
				Body:     Payload{Jsonrpc: "2.0", Method: "eth_getBalance", ID: 1},
			},
			Func:        GetEthBalance,
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
