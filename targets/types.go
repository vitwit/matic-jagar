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
		Body        Payload
		Method      string
	}

	Payload struct {
		Jsonrpc string        `json:"jsonrpc"`
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		ID      int           `json:"id"`
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
				Moniker string `json:"moniker"`
			} `json:"node_info"`
			SyncInfo struct {
				LatestBlockHeight string `json:"latest_block_height"`
				LatestBlockTime   string `json:"latest_block_time"`
				CatchingUp        bool   `json:"catching_up"`
			} `json:"sync_info"`
			ValidatorInfo struct {
				Address     string `json:"address"`
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
					Moniker    string `json:"moniker"`
				} `json:"node_info"`
				RemoteIP string `json:"remote_ip"`
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

	LastProposedBlockAndTime struct {
		BlockMeta interface{} `json:"block_meta"`
		Block     struct {
			Header struct {
				ChainID         string `json:"chain_id"`
				Height          string `json:"height"`
				Time            string `json:"time"`
				NumTxs          string `json:"num_txs"`
				TotalTxs        string `json:"total_txs"`
				LastCommitHash  string `json:"last_commit_hash"`
				ProposerAddress string `json:"proposer_address"`
			} `json:"header"`
		} `json:"block"`
	}

	AccountBalResp struct {
		Height string `json:"height"`
		Result []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"result"`
	}

	Proposals struct {
		Height string `json:"height"`
		Result []struct {
			Content struct {
				Type  string `json:"type"`
				Value struct {
					Title       string      `json:"title"`
					Description string      `json:"description"`
					Changes     interface{} `json:"changes"`
				} `json:"value"`
			} `json:"content"`
			ID               string        `json:"id"`
			ProposalStatus   string        `json:"proposal_status"`
			FinalTallyResult interface{}   `json:"final_tally_result"`
			SubmitTime       string        `json:"submit_time"`
			DepositEndTime   string        `json:"deposit_end_time"`
			TotalDeposit     []interface{} `json:"total_deposit"`
			VotingStartTime  string        `json:"voting_start_time"`
			VotingEndTime    string        `json:"voting_end_time"`
		} `json:"result"`
	}

	ProposalVoters struct {
		Height string `json:"height"`
		Result []struct {
			ProposalID string `json:"proposal_id"`
			Voter      string `json:"voter"`
			Option     string `json:"option"`
		} `json:"result"`
	}

	// Depositors struct which holds the parameters of depositors
	Depositors struct {
		Height string `json:"height"`
		Result []struct {
			ProposalID string `json:"proposal_id"`
			Depositor  string `json:"depositor"`
			Amount     []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"amount"`
		} `json:"result"`
	}

	BorParams struct {
		Height string `json:"height"`
		Result struct {
			SprintDuration int `json:"sprint_duration"`
			SpanDuration   int `json:"span_duration"`
			ProducerCount  int `json:"producer_count"`
		} `json:"result"`
	}

	TotalCheckpoints struct {
		Height string `json:"height"`
		Result struct {
			Result int `json:"result"`
		} `json:"result"`
	}

	BorResult struct {
		Result string `json:"result"`
		ID     int    `json:"id"`
	}

	LatestCheckpoints struct {
		Height string `json:"height"`
		Result struct {
			Proposer   string `json:"proposer"`
			StartBlock int    `json:"start_block"`
			EndBlock   int    `json:"end_block"`
			RootHash   string `json:"root_hash"`
			BorChainID string `json:"bor_chain_id"`
			Timestamp  int    `json:"timestamp"`
		} `json:"result"`
	}

	CheckpointsDuration struct {
		Height string `json:"height"`
		Result struct {
			CheckpointBufferTime    int64 `json:"checkpoint_buffer_time"`
			AvgCheckpointLength     int   `json:"avg_checkpoint_length"`
			MaxCheckpointLength     int   `json:"max_checkpoint_length"`
			ChildChainBlockInterval int   `json:"child_chain_block_interval"`
		} `json:"result"`
	}

	BorBoolResp struct {
		Result bool `json:"result"`
		ID     int  `json:"id"`
	}

	EthResult struct {
		Result string `json:"result"`
		ID     int    `json:"id"`
	}

	BorSignersRes struct {
		Result []string `json:"result"`
		ID     int      `json:"id"`
	}
)
