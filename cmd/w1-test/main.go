package main

import (
	"github.com/SpComb/iot-poc/netlink/connector"
	"github.com/SpComb/iot-poc/netlink/connector/w1"
	"github.com/qmsk/go-logging"

	"flag"
	"fmt"
	"os"
)

var log logging.Logging

var options struct {
	Log          logging.Options
	LogW1        logging.Options
	LogConnector logging.Options
}

func init() {
	options.LogW1.Module = "w1"
	options.LogW1.Defaults = &options.Log
	options.LogConnector.Module = "connector"
	options.LogConnector.Defaults = &options.Log

	options.Log.InitFlags()
	options.LogW1.InitFlags()
	options.LogConnector.InitFlags()
}

func run() error {
	w1conn, err := w1.Dial()
	if err != nil {
		return fmt.Errorf("w1.Dial: %v", err)
	} else {
		log.Infof("Connected to w1-netlink: %v", w1conn)
	}

	if masters, err := w1conn.ListMasters(); err != nil {
		return fmt.Errorf("w1.ListMasters: %v", err)
	} else {
		for _, master := range masters {
			log.Infof("Master: %v", master)
		}
	}

	return nil
}

func main() {
	flag.Parse()

	log = options.Log.MakeLogging()

	connector.SetLogging(options.LogConnector.MakeLogging())
	w1.SetLogging(options.LogW1.MakeLogging())

	if err := run(); err != nil {
		log.Errorf("%v", err)
		os.Exit(1)
	}
}
