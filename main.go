package main

import (
	"io/ioutil"
	"os"

	"time"

	"fmt"

	"github.com/BurntSushi/toml"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s config.toml\n", os.Args[0])
		os.Exit(1)
	}

	configPath := os.Args[1]
	tomlData, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("unable to read %s: %s\n", configPath, err)
		os.Exit(1)
	}

	var config Config
	if _, err := toml.Decode(string(tomlData), &config); err != nil {
		fmt.Printf("unable to parse %s: %s\n", configPath, err)
		os.Exit(1)
	}

	setupInflux(config.Influx)
	go statsWriter(time.Duration(config.Influx.Interval) * time.Second)
	capture(config.Capture)
}
