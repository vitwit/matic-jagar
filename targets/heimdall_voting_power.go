package targets

import (
	"encoding/json"
	"fmt"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetValidatorVotingPower to get voting power of a validator
func GetValidatorVotingPower(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var validatorResp ValStatusResp
	err = json.Unmarshal(resp.Body, &validatorResp)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	vp := validatorResp.Result.Power
	_ = writeToInfluxDb(c, bp, "matic_voting_power", map[string]string{}, map[string]interface{}{"power": vp})
	log.Println("Voting Power \n", vp)

	if int64(vp) <= cfg.VotingPowerThreshold {
		_ = SendTelegramAlert(fmt.Sprintf("Your validator %s voting power has dropped below %d", cfg.ValidatorName, cfg.VotingPowerThreshold), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Your validator %s voting power has dropped below %d", cfg.ValidatorName, cfg.VotingPowerThreshold), cfg)
	}
}
