package main

// Config is the top configuration object
type Config struct {
	Capture CaptureConfig
	Influx  InfluxConfig
}

// CaptureConfig is the configuration for the pcap capturing
type CaptureConfig struct {
	Interface string
	Filter    string
}

// InfluxConfig is the configuration for InfluxDB
type InfluxConfig struct {
	Host        string
	Username    string
	Password    string
	Database    string
	Measurement string
	Interval    uint
}
