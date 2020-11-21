package targets

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetLatency to calculate latency of a peer address
func GetLatency(_ HTTPOptions, cfg *config.Config, c client.Client) error {
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