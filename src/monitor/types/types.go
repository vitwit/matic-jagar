package types

import (
	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

type (
	// QueryParams map of strings
	QueryParams map[string]string

	// HTTPOptions is a structure that holds all http options parameters
	HTTPOptions struct {
		Endpoint    string
		QueryParams QueryParams
		Body        []byte
		Method      string
	}

	// Target is a structure which holds all the parameters of a target
	//this could be used to write endpoints for each functionality
	Target struct {
		ExecutionType string
		HTTPOptions   HTTPOptions
		Name          string
		Func          func(m HTTPOptions, cfg *config.Config, c client.Client)
		ScraperRate   string
	}

	// Targets list of all the targets
	Targets struct {
		List []Target
	}

	// PingResp is a structure which holds the options of a response
	PingResp struct {
		StatusCode int
		Body       []byte
	}

	// Peer is a structure which holds the info about a peer address
	Peer struct {
		RemoteIP         string      `json:"remote_ip"`
		ConnectionStatus interface{} `json:"connection_status"`
		IsOutbound       bool        `json:"is_outbound"`
		NodeInfo         struct {
			Moniker string `json:"moniker"`
			Network string `json:"network"`
		} `json:"node_info"`
	}

	// NetInfoResult struct
	NetInfoResult struct {
		Listening bool          `json:"listening"`
		Listeners []interface{} `json:"listeners"`
		NumPeers  string        `json:"n_peers"`
		Peers     []Peer        `json:"peers"`
	}

	// NetInfo is a structre which holds the details of address
	NetInfo struct {
		JSONRpc string        `json:"jsonrpc"`
		ID      string        `json:"id"`
		Result  NetInfoResult `json:"result"`
	}
)
