package main

import (
	"flag"
	"log"
	"os"
)

var (
	configFile = flag.String("config", "config.yaml", "Path to configuration file.")
	debug      = flag.Bool("debug", false, "Enable debug logs.")
)

func realMain() error {
	flag.Parse()
	config, err := LoadFromFile(*configFile)
	if err != nil {
		return err
	}
	scanner, err := NewScanner(config, *debug)
	if err != nil {
		return err
	}
	err = scanner.Scan()
	if err != nil {
		return err
	}
	return nil
}

func main() {
	err := realMain()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
