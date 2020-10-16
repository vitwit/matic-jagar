package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// GetBorLatestSpan to get latest span id and also calcualte span validator count
func GetBorLatestSpan(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var latestSpan BorLatestSpan
	err = json.Unmarshal(resp.Body, &latestSpan)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	spanID := latestSpan.Result.SpanID

	// Get previous span id from db
	prevSpanID := GetBorSpanIDFromDb(cfg, c)
	prevSpan, _ := strconv.Atoi(prevSpanID)

	addrExists := false

	for _, val := range latestSpan.Result.ValidatorSet.Validators {
		if val.Signer == cfg.ValDetails.SignerAddress {
			addrExists = true
		}
	}

	count := GetBorSpanValidatorCountFromDb(cfg, c)
	spanValCount, _ := strconv.Atoi(count)

	if addrExists {
		diff := spanID - prevSpan
		if diff > 0 {
			spanValCount = spanValCount + 1
		}
	}

	_ = writeToInfluxDb(c, bp, "matic_bor_latest_span", map[string]string{}, map[string]interface{}{"span_id": spanID, "span_val_count": spanValCount})
	log.Printf("Bor Latest Span ID: %d and Span Val Count : %d", spanID, spanValCount)
}

// GetBorSpanIDFromDb returns the span ID from db
func GetBorSpanIDFromDb(cfg *config.Config, c client.Client) string {
	var spanID string
	q := client.NewQuery("SELECT last(span_id) FROM matic_bor_latest_span", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						value := r.Series[0].Values[0][idx]
						spanID = fmt.Sprintf("%v", value)
						break
					}
				}
			}
		}
	}
	return spanID
}

// GetBorSpanValidatorCountFromDb returns the span val count from the db
func GetBorSpanValidatorCountFromDb(cfg *config.Config, c client.Client) string {
	var count string
	q := client.NewQuery("SELECT last(span_val_count) FROM matic_bor_latest_span", cfg.InfluxDB.Database, "")
	if response, err := c.Query(q); err == nil && response.Error() == nil {
		for _, r := range response.Results {
			if len(r.Series) != 0 {
				for idx, col := range r.Series[0].Columns {
					if col == "last" {
						value := r.Series[0].Values[0][idx]
						count = fmt.Sprintf("%v", value)
						break
					}
				}
			}
		}
	}
	return count
}
