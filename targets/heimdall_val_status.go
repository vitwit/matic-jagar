package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// ValidatorStatusAlert to send alerts to telegram and email about validator status
func ValidatorStatusAlert(ops HTTPOptions, cfg *config.Config, c client.Client) {
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

	alertTime1 := cfg.AlertTime1
	alertTime2 := cfg.AlertTime2

	t1, _ := time.Parse(time.Kitchen, alertTime1)
	t2, _ := time.Parse(time.Kitchen, alertTime2)

	now := time.Now().UTC()
	t := now.Format(time.Kitchen)

	a1 := t1.Format(time.Kitchen)
	a2 := t2.Format(time.Kitchen)

	log.Println("a1, a2 and present time : ", a1, a2, t)

	valID := validatorResp.Result.ID
	validatorStatus := validatorResp.Result.Jailed
	log.Println("val status: ", validatorStatus)

	if !validatorStatus {
		if t == a1 || t == a2 {
			_ = SendTelegramAlert(fmt.Sprintf("Your validator %s is currently voting", cfg.ValidatorName), cfg)
			_ = SendEmailAlert(fmt.Sprintf("Your validator %s is currently voting", cfg.ValidatorName), cfg)
			log.Println("Sent validator status alert")
		}
		_ = writeToInfluxDb(c, bp, "matic_val_status", map[string]string{}, map[string]interface{}{"status": 1, "val_id": valID})
	} else {
		_ = SendTelegramAlert(fmt.Sprintf("Your validator %s is in jailed status", cfg.ValidatorName), cfg)
		_ = SendEmailAlert(fmt.Sprintf("Your validator %s is in jailed status", cfg.ValidatorName), cfg)
		log.Println("Sent validator status alert")

		_ = writeToInfluxDb(c, bp, "matic_val_status", map[string]string{}, map[string]interface{}{"status": 0, "val_id": valID})
	}
	return
}

// GetValID returns ID of the validator from db
func GetValID(cfg *config.Config, c client.Client) string {
	var ID string
	q := client.NewQuery("SELECT last(val_id) FROM matic_val_status", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						value := r.Series[0].Values[0][idx]
						ID = fmt.Sprintf("%v", value)
						break
					}
				}
			}
		}
	}
	return ID
}

// GetValStatusFromDB returns latest current height from db
func GetValStatusFromDB(cfg *config.Config, c client.Client) string {
	var valStatus string
	q := client.NewQuery("SELECT last(status) FROM matic_val_status", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						status := r.Series[0].Values[0][idx]
						valStatus = fmt.Sprintf("%v", status)
						break
					}
				}
			}
		}
	}
	return valStatus
}
