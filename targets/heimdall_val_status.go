package targets

import (
	"fmt"
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/alerter"
	"github.com/vitwit/matic-jagar/config"
	db "github.com/vitwit/matic-jagar/influxdb"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// ValidatorStatusAlert will checks whether the validator is voting or jailed
// Alerter will send alerts according to the timings of regualr status alerting configured in config.toml
func ValidatorStatusAlert(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := db.CreateBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	validatorResp, err := scraper.GetValStatus(ops)
	if err != nil {
		log.Printf("Error in validator status : %v", err)
		return
	}

	if &validatorResp.Result == nil {
		log.Printf("Got an empty response of validator status : %v", err)
		return
	}

	now := time.Now().UTC()
	currentTime := now.Format(time.Kitchen)

	var alertsArray []string

	for _, value := range cfg.RegularStatusAlerts.AlertTimings {
		t, _ := time.Parse(time.Kitchen, value)
		alertTime := t.Format(time.Kitchen)

		alertsArray = append(alertsArray, alertTime)
	}

	log.Printf("current time : %s, and Status Alert Timings : %v", currentTime, alertsArray)

	valID := validatorResp.Result.ID
	validatorStatus := validatorResp.Result.Jailed
	log.Println("val status: ", validatorStatus)

	if !validatorStatus {
		for _, statusAlertTime := range alertsArray {
			if currentTime == statusAlertTime {
				_ = alerter.SendTelegramAlert(fmt.Sprintf("Your validator %s is currently voting", cfg.ValDetails.ValidatorName), cfg)
				_ = alerter.SendEmailAlert(fmt.Sprintf("Your validator %s is currently voting", cfg.ValDetails.ValidatorName), cfg)
				log.Println("Sent validator status alert")
			}
		}
		_ = db.WriteToInfluxDb(c, bp, "heimdall_val_status", map[string]string{}, map[string]interface{}{"status": 1, "val_id": valID})
	} else {
		for _, statusAlertTime := range alertsArray {
			if currentTime == statusAlertTime {
				_ = alerter.SendTelegramAlert(fmt.Sprintf("Your validator %s is in jailed status", cfg.ValDetails.ValidatorName), cfg)
				_ = alerter.SendEmailAlert(fmt.Sprintf("Your validator %s is in jailed status", cfg.ValDetails.ValidatorName), cfg)
				log.Println("Sent validator status alert")
			}
		}

		_ = db.WriteToInfluxDb(c, bp, "heimdall_val_status", map[string]string{}, map[string]interface{}{"status": 0, "val_id": valID})
	}
	return
}

// GetValID returns ID of the validator from db
func GetValID(cfg *config.Config, c client.Client) string {
	var ID string
	q := client.NewQuery("SELECT last(val_id) FROM heimdall_val_status", cfg.InfluxDB.Database, "")
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
	q := client.NewQuery("SELECT last(status) FROM heimdall_val_status", cfg.InfluxDB.Database, "")
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
