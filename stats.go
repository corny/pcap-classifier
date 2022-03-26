package main

import (
	"strings"
	"sync"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
)

var (
	stats = make(Stats)
	mtx   = sync.Mutex{}
)

// Counter counts bytes and packets per packet type
type Counter struct {
	Bytes   int64
	Packets uint
}

// Stats is a map from packet type to counter
type Stats map[string]*Counter

// Adds a new packet to the statistics
func (stats Stats) add(key string, size uint16) {
	mtx.Lock()
	defer mtx.Unlock()

	counter, _ := stats[key]

	if counter == nil {
		counter = &Counter{}
		stats[key] = counter
	}

	counter.Bytes += int64(size)
	counter.Packets++
}

// Writes the data periodically
func statsWriter(interval time.Duration) {
	for range time.NewTicker(interval).C {
		writeStats()
	}
}

// Write the data into the database
func writeStats() {
	mtx.Lock()
	current := stats
	stats = make(Stats)
	mtx.Unlock()

	now := time.Now()

	bp, err := client.NewBatchPoints(influxBatchConfig)
	if err != nil {
		panic(bp)
	}

	for key, counter := range current {
		tags := make(map[string]string)

		if i := strings.IndexByte(key, '-'); i != -1 {
			tags["proto"] = key[:i]
			tags["type"] = key[i+1:]
		} else {
			tags["proto"] = key
		}

		point, err := client.NewPoint(
			influxMeasurement,
			tags,
			models.Fields{
				"bytes":   counter.Bytes,
				"packets": counter.Packets,
			},
			now,
		)

		if err != nil {
			panic(err)
		}
		bp.AddPoint(point)
	}

	if err = influxClient.Write(bp); err != nil {
		panic(err)
	}
}
