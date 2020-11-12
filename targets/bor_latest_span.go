package targets

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
	"github.com/vitwit/matic-jagar/scraper"
	"github.com/vitwit/matic-jagar/types"
)

// GetBorLatestSpan to get latest span id and also calcualte span validator count
func GetBorLatestSpan(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	latestSpan, err := scraper.BorLatestSpan(ops)
	if err != nil {
		log.Printf("Error in BorLatestSpan: %v", err)
		return
	}

	spanID := latestSpan.Result.SpanID

	// Get previous span id from db
	prevSpanID := GetBorSpanIDFromDb(cfg, c)
	if prevSpanID == "" {
		prevSpanID = "0"
	}
	prevSpan, err := strconv.Atoi(prevSpanID)
	if err != nil {
		log.Printf("Error in conversion from string to int of span ID : %v", err)
		return
	}

	addrExists := false

	for _, val := range latestSpan.Result.ValidatorSet.Validators {
		if strings.EqualFold(val.Signer, cfg.ValDetails.SignerAddress) {
			addrExists = true
		}
	}

	count := GetBorSpanValidatorCountFromDb(cfg, c)
	if count == "" {
		count = "0"
	}
	spanValCount, err := strconv.Atoi(count)
	if err != nil {
		log.Printf("Error in string convertion to int : %v", err)
		return
	}

	if addrExists {
		diff := spanID - prevSpan
		if diff > 0 {
			spanValCount = spanValCount + 1
		}
	}

	_ = writeToInfluxDb(c, bp, "bor_latest_span", map[string]string{}, map[string]interface{}{"span_id": spanID, "span_val_count": spanValCount})
	log.Printf("Bor Latest Span ID: %d and Span Val Count : %d", spanID, spanValCount)
}

// GetBorSpanIDFromDb returns the span ID from db
func GetBorSpanIDFromDb(cfg *config.Config, c client.Client) string {
	var spanID string
	q := client.NewQuery("SELECT last(span_id) FROM bor_latest_span", cfg.InfluxDB.Database, "")
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
	q := client.NewQuery("SELECT last(span_val_count) FROM bor_latest_span", cfg.InfluxDB.Database, "")
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

// GetBlockProducer is to get the proucer counts and checks
//weather the validator is part of block producer or not
func GetBlockProducer(ops types.HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	currentSpan := GetBorSpanIDFromDb(cfg, c)

	ops.Endpoint = ops.Endpoint + currentSpan

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var spanProducers types.BorSpanProducers
	err = json.Unmarshal(resp.Body, &spanProducers)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	addrExists := "No"

	for _, value := range spanProducers.Result.SelectedProducers {
		if strings.EqualFold(value.Signer, cfg.ValDetails.SignerAddress) {
			addrExists = "Yes"
		}
	}

	producerCount := len(spanProducers.Result.SelectedProducers)
	_ = writeToInfluxDb(c, bp, "bor_block_producer", map[string]string{}, map[string]interface{}{"val_part_of_block_producer": addrExists, "producer_count": producerCount})
	log.Printf("Validator is part of block producer : %s\n and Producer Count: %d", addrExists, producerCount)
}
