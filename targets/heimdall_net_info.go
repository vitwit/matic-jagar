package targets

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// NetInfo is to get no.of peers, addresses and also calculates it's alatency and stores it in db
func NetInfo(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}
	var pts []*client.Point

	ni, err := scraper.GetNetInfo(ops)
	if err != nil {
		log.Printf("Error in net nfo: %v", err)
		return
	}

	numPeers, err := strconv.Atoi(ni.Result.NPeers)
	if err != nil {
		log.Printf("Error converting num_peers to int: %v", err)
		numPeers = 0
	} else if int64(numPeers) < cfg.AlertingThresholds.NumPeersThreshold && strings.ToUpper(cfg.AlerterPreferences.NumPeersAlerts) == "YES" {
		_ = SendTelegramAlert(fmt.Sprintf("Number of peers connected to your validator has fallen below %d", cfg.AlertingThresholds.NumPeersThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Number of peers connected to your validator has fallen below %d", cfg.AlertingThresholds.NumPeersThreshold), cfg)
	}
	p1, err := createDataPoint("heimdall_num_peers", map[string]string{}, map[string]interface{}{"count": numPeers})
	if err == nil {
		pts = append(pts, p1)
	}

	peerAddrs := make([]string, len(ni.Result.Peers))
	for i, peer := range ni.Result.Peers {
		peerAddrs[i] = peer.RemoteIP + " - " + peer.NodeInfo.Moniker
	}

	addrs := strings.Join(peerAddrs[:], ",  ")
	p2, err := createDataPoint("heimdall_peer_addresses", map[string]string{"addresses_count": strconv.Itoa(numPeers)}, map[string]interface{}{"addresses": addrs})
	if err == nil {
		pts = append(pts, p2)
	}

	bp.AddPoints(pts)
	_ = writeBatchPoints(c, bp)
	log.Printf("No. of peers: %d \n", numPeers)

	// Calling funtion to get peer latency
	err = PeerLatency(ops, cfg, c)
	if err != nil {
		log.Printf("Error while calculating peer latency : %v", err)
		return
	}
}

// PeerLatency is to calculate latency of a peer address and stores it in db
func PeerLatency(_ types.HTTPOptions, cfg *config.Config, c client.Client) error {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return err
	}

	q := client.NewQuery(fmt.Sprintf("SELECT * FROM heimdall_peer_addresses"), cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		var addresses []string
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				noOfValues := len(r.Series[0].Values)
				if noOfValues != 0 {
					n := noOfValues - 1
					addressValues := fmt.Sprintf("%v", r.Series[0].Values[n][1])
					addresses = strings.Split(addressValues, ", ")
				}
			}
		}
		for _, addr := range addresses {
			log.Printf("peer address %s", addr)
			cmd := exec.Command("ping", "-c", "5", addr)
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("Error while running ping command %v", err)
				return err
			}
			pingResp := string(out)
			rtt := pingResp[len(pingResp)-35 : len(pingResp)-1]
			splitString := strings.Split(rtt, "/")
			avgRtt := splitString[1]
			log.Println("Writing address latency in db ", addr, avgRtt)
			err = writeToInfluxDb(c, bp, "heimdall_validator_latency", map[string]string{"peer_address": addr}, map[string]interface{}{"address": addr, "avg_rtt": avgRtt})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// GetPeersCount returns count of peer addresses from db
func GetPeersCount(cfg *config.Config, c client.Client) string {
	var count string
	q := client.NewQuery("SELECT last(count) FROM heimdall_num_peers", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						c := r.Series[0].Values[0][idx]
						count = fmt.Sprintf("%v", c)
						break
					}
				}
			}
		}
	}

	return count
}
