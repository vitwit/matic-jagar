package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetNetInfo to get no.of peers and addresses
func GetNetInfo(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}
	var pts []*client.Point

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error getting node_info: %v", err)
		return
	}
	var ni NetInfo
	err = json.Unmarshal(resp.Body, &ni)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	numPeers, err := strconv.Atoi(ni.Result.NPeers)
	if err != nil {
		log.Printf("Error converting num_peers to int: %v", err)
		numPeers = 0
	} else if int64(numPeers) < cfg.NumPeersThreshold {
		_ = SendTelegramAlert(fmt.Sprintf("Number of peers connected to your validator has fallen below %d", cfg.NumPeersThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Number of peers connected to your validator has fallen below %d", cfg.NumPeersThreshold), cfg)
	}
	p1, err := createDataPoint("matic_num_peers", map[string]string{}, map[string]interface{}{"count": numPeers})
	if err == nil {
		pts = append(pts, p1)
	}

	peerAddrs := make([]string, len(ni.Result.Peers))
	for i, peer := range ni.Result.Peers {
		peerAddrs[i] = peer.RemoteIP + " - " + peer.NodeInfo.Moniker
	}

	addrs := strings.Join(peerAddrs[:], ",  ")
	p2, err := createDataPoint("matic_peer_addresses", map[string]string{"addresses_count": strconv.Itoa(numPeers)}, map[string]interface{}{"addresses": addrs})
	if err == nil {
		pts = append(pts, p2)
	}

	bp.AddPoints(pts)
	_ = writeBatchPoints(c, bp)
	log.Printf("No. of peers: %d \n", numPeers)
}

func BorPeersCount(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error getting node_info: %v", err)
		return
	}

	var peerInfo BorPeersInfo
	err = json.Unmarshal(resp.Body, &peerInfo)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	count := HexToIntConversion(peerInfo.Result)

	_ = writeToInfluxDb(c, bp, "matic_bor_peers", map[string]string{}, map[string]interface{}{"address_count": count})
	log.Printf("Bor Peers Count: %d", count)
}

func HexToIntConversion(hex string) int {
	val := hex[2:]

	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		panic(err)
	}
	n2 := int(n)

	return n2
}
