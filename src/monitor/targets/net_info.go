package targets

import (
	"encoding/json"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/src/monitor/types"
)

// GetNetInfo to get no.of peers and addresses
func GetNetInfo(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error getting node_info: %v", err)
		return
	}
	var ni types.NetInfo
	err = json.Unmarshal(resp.Body, &ni)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	numPeers, err := strconv.Atoi(ni.Result.NumPeers)
	log.Printf("No. of peers: %d \n Peer Addresses: %v", numPeers)
}
