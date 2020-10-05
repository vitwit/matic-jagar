package targets

import (
	"net/http"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/src/monitor/targets"
	"github.com/vitwit/matic-jagar/src/monitor/types"
)

type targetRunner struct{}

// // NewRunner returns targetRunner
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
				Endpoint: cfg.RPCEndpoint + "/net_info?",
				Method:   http.MethodGet,
			},
			Func:        targets.GetNetInfo,
			ScraperRate: "2s",
		},
	}}
}
