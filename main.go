package main

import (
	"flag"

	"github.com/sirupsen/logrus"
)

var (
	configFile = flag.String("config", "config.yaml", "Path to configuration file.")
	debug      = flag.Bool("debug", false, "Enable debug logs.")
)

func main() {
	flag.Parse()
	config, err := LoadFromFile(*configFile)
	if err != nil {
		logrus.Fatal(err)
	}
	scanner, err := NewScanner(config, *debug)
	if err != nil {
		logrus.Fatal(err)
	}
	err = scanner.Scan()
	if err != nil {
		logrus.Fatal(err)
	}
}
