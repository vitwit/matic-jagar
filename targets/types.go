package targets

import (
	"time"

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

	Status struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			NodeInfo struct {
				ProtocolVersion interface{} `json:"protocol_version"`
				ID              string      `json:"id"`
				ListenAddr      string      `json:"listen_addr"`
				Network         string      `json:"network"`
				Version         string      `json:"version"`
				Channels        string      `json:"channels"`
				Moniker         string      `json:"moniker"`
				Other           struct {
					TxIndex    string `json:"tx_index"`
					RPCAddress string `json:"rpc_address"`
				} `json:"other"`
			} `json:"node_info"`
			SyncInfo struct {
				LatestBlockHash   string `json:"latest_block_hash"`
				LatestAppHash     string `json:"latest_app_hash"`
				LatestBlockHeight string `json:"latest_block_height"`
				LatestBlockTime   string `json:"latest_block_time"`
				CatchingUp        bool   `json:"catching_up"`
			} `json:"sync_info"`
			ValidatorInfo struct {
				Address string `json:"address"`
				PubKey  struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"pub_key"`
				VotingPower string `json:"voting_power"`
			} `json:"validator_info"`
		} `json:"result"`
	}

	NetInfo struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			Listening bool     `json:"listening"`
			Listeners []string `json:"listeners"`
			NPeers    string   `json:"n_peers"`
			Peers     []struct {
				NodeInfo struct {
					ProtocolVersion struct {
						P2P   string `json:"p2p"`
						Block string `json:"block"`
						App   string `json:"app"`
					} `json:"protocol_version"`
					ID         string `json:"id"`
					ListenAddr string `json:"listen_addr"`
					Network    string `json:"network"`
					Version    string `json:"version"`
					Channels   string `json:"channels"`
					Moniker    string `json:"moniker"`
					Other      struct {
						TxIndex    string `json:"tx_index"`
						RPCAddress string `json:"rpc_address"`
					} `json:"other"`
				} `json:"node_info"`
				IsOutbound       bool        `json:"is_outbound"`
				ConnectionStatus interface{} `json:"connection_status"`
				RemoteIP         string      `json:"remote_ip"`
			} `json:"peers"`
		} `json:"result"`
	}

	ValidatorsHeight struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			BlockHeight string `json:"block_height"`
			Validators  []struct {
				Address string `json:"address"`
				PubKey  struct {
					Type  string `json:"type"`
					Value string `json:"value"`
				} `json:"pub_key"`
				VotingPower      string `json:"voting_power"`
				ProposerPriority string `json:"proposer_priority"`
			} `json:"validators"`
		} `json:"result"`
	}

	CurrentBlockWithHeight struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			BlockMeta interface{} `json:"block_meta"`
			Block     struct {
				Header struct {
					Height string `json:"height"`
					Time   string `json:"time"`
				} `json:"header"`
				Data struct {
					Txs interface{} `json:"txs"`
				} `json:"data"`
				Evidence struct {
					Evidence interface{} `json:"evidence"`
				} `json:"evidence"`
				LastCommit struct {
					BlockID    interface{} `json:"block_id"`
					Precommits []struct {
						Type             int         `json:"type"`
						Height           string      `json:"height"`
						Round            string      `json:"round"`
						BlockID          interface{} `json:"block_id"`
						Timestamp        time.Time   `json:"timestamp"`
						ValidatorAddress string      `json:"validator_address"`
						ValidatorIndex   string      `json:"validator_index"`
						Signature        string      `json:"signature"`
						SideTxResults    interface{} `json:"side_tx_results"`
					} `json:"precommits"`
				} `json:"last_commit"`
			} `json:"block"`
		} `json:"result"`
	}

	// UnconfirmedTxns struct which holds the parameters of unconfirmed txns
	UnconfirmedTxns struct {
		Jsonrpc string `json:"jsonrpc"`
		Result  struct {
			NTxs       string      `json:"n_txs"`
			Total      string      `json:"total"`
			TotalBytes string      `json:"total_bytes"`
			Txs        interface{} `json:"txs"`
		} `json:"result"`
	}

	ApplicationInfo struct {
		NodeInfo           interface{} `json:"node_info"`
		ApplicationVersion struct {
			Name       string `json:"name"`
			ServerName string `json:"server_name"`
			ClientName string `json:"client_name"`
			Version    string `json:"version"`
			Commit     string `json:"commit"`
			BuildTags  string `json:"build_tags"`
			Go         string `json:"go"`
		} `json:"application_version"`
	}

	AuthParams struct {
		Height string `json:"height"`
		Result struct {
			MaxMemoCharacters      int    `json:"max_memo_characters"`
			TxSigLimit             int    `json:"tx_sig_limit"`
			TxSizeCostPerByte      int    `json:"tx_size_cost_per_byte"`
			SigVerifyCostEd25519   int    `json:"sig_verify_cost_ed25519"`
			SigVerifyCostSecp256K1 int    `json:"sig_verify_cost_secp256k1"`
			MaxTxGas               int    `json:"max_tx_gas"`
			TxFees                 string `json:"tx_fees"`
		} `json:"result"`
	}

	ValStatusResp struct {
		Height string `json:"height"`
		Result struct {
			ID          int    `json:"ID"`
			StartEpoch  int    `json:"startEpoch"`
			EndEpoch    int    `json:"endEpoch"`
			Nonce       int    `json:"nonce"`
			Power       int    `json:"power"`
			PubKey      string `json:"pubKey"`
			Signer      string `json:"signer"`
			LastUpdated string `json:"last_updated"`
			Jailed      bool   `json:"jailed"`
			Accum       int    `json:"accum"`
		} `json:"result"`
	}
)
