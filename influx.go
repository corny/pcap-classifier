package main

import "github.com/influxdata/influxdb/client/v2"

var (
	influxClient      client.Client
	influxMeasurement string
	influxBatchConfig = client.BatchPointsConfig{
		Precision: "m",
	}
)

func setupInflux(config InfluxConfig) {
	// Make client
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.Host,
		Username: config.Username,
		Password: config.Password,
	})

	if err != nil {
		panic(err)
	}

	influxClient = c
	influxMeasurement = config.Measurement
	influxBatchConfig.Database = config.Database
}
