package targets

import (
	"encoding/json"
	"log"

	client "github.com/influxdata/influxdb1-client/v2"

	"github.com/vitwit/matic-jagar/config"
)

// NodeVersion to get application version from the LCD
func NodeVersion(ops HTTPOptions, cfg *config.Config, c client.Client) {
	bp, err := createBatchPoints(cfg.InfluxDB.Database)
	if err != nil {
		return
	}

	resp, err := HitHTTPTarget(ops)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	var applicationInfo ApplicationInfo
	err = json.Unmarshal(resp.Body, &applicationInfo)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	appVersion := applicationInfo.ApplicationVersion.Version

	_ = writeToInfluxDb(c, bp, "matic_version", map[string]string{}, map[string]interface{}{"v": appVersion})
	log.Printf("Version: %s", appVersion)
}
