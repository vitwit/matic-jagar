package targets

import (
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/alerter"
	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// ValidatorVotingPower is to get voting power of a validator and stores it in db
// Alerter will notify if there is any change in voting power
func ValidatorVotingPower(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	validatorResp, err := scraper.GetValStatus(ops)
	if err != nil {
		log.Printf("Error in validator voting power : %v", err)
		return
	}

	vp := validatorResp.Result.Power
	prevVotingPower := GetVotingPowerFromDb(cfg, c)
	previousVP, _ := strconv.Atoi(prevVotingPower)

	if previousVP != vp {
		_ = alerter.SendTelegramAlert(fmt.Sprintf("Voting Power Alert : Your validator voting power has changed from %d to %d", previousVP, vp), cfg)
		_ = alerter.SendEmailAlert(fmt.Sprintf("Voting Power Alert : Your validator voting power has changed from %d to %d", previousVP, vp), cfg)
	}

	_ = db.WriteToInfluxDb(c, bp, "heimdall_voting_power", map[string]string{}, map[string]interface{}{"power": vp})
	log.Println("Voting Power \n", vp)
}

// GetVotingPowerFromDb returns voting power of a validator from db
func GetVotingPowerFromDb(cfg *config.Config, c client.Client) string {
	var vp string
	q := client.NewQuery("SELECT last(power) FROM heimdall_voting_power", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						v := r.Series[0].Values[0][idx]
						vp = fmt.Sprintf("%v", v)
						break
					}
				}
			}
		}
	}
	return vp
}
