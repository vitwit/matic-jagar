package influxdb

import (
	"log"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

// CreateDataPoint to create a data point
func CreateDataPoint(name string, tags map[string]string, fields map[string]interface{}) (*client.Point, error) {
	p, err := client.NewPoint(name, tags, fields, time.Now())
	if err != nil {
		log.Printf("Error creating data point: %v", err)
		return nil, err
	}
	return p, nil
}

// CreateBatchPoints to create batch points
func CreateBatchPoints(db string) (client.BatchPoints, error) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db,
		Precision: "s",
	})
	if err != nil {
		log.Printf("err creating batch points: %v", err)
		return nil, err
	}
	return bp, nil
}

// WriteBatchPoints is to write data into influx db
func WriteBatchPoints(c client.Client, bp client.BatchPoints) error {
	if err := c.Write(bp); err != nil {
		log.Printf("err writing batch points to client: %v", err)
		return err
	}
	return nil
}

// WriteToInfluxDb is to create data points and writes the data into influxdb by calling WriteBatchPoints
func WriteToInfluxDb(c client.Client, bp client.BatchPoints, name string, tags map[string]string,
	fields map[string]interface{}) error {
	p, err := CreateDataPoint(name, tags, fields)
	if err != nil {
		return err
	}
	bp.AddPoint(p)
	err = WriteBatchPoints(c, bp)
	if err != nil {
		return err
	}
	return nil
}
