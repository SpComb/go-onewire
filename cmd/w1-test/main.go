package main

import (
	"github.com/SpComb/go-onewire/netlink/connector"
	"github.com/SpComb/go-onewire/netlink/connector/w1"
	"github.com/SpComb/go-onewire/netlink/connector/w1/ds18b20"
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

func scan(w1conn *w1.Conn) error {
	if masters, err := w1conn.ListMasters(); err != nil {
		return fmt.Errorf("w1.ListMasters: %v", err)
	} else {
		for _, masterID := range masters {
			log.Infof("Master: %v", masterID)

			if slaves, err := w1conn.ListSlaves(masterID); err != nil {
				return fmt.Errorf("w1.ListSlaves %v: %v", masterID, err)
			} else {
				for _, slaveID := range slaves {
					switch slaveID.Family {
					case ds18b20.Family:
						var device = ds18b20.MakeDevice(w1conn, slaveID)

						if s, err := device.Read(); err != nil {
							log.Errorf("Slave %v: Read: %v", slaveID, err)
						} else {
							log.Infof("DS18B20 %v: %#v", slaveID, s)
						}

						if err := device.ConvertT(); err != nil {
							log.Errorf("Slave %v: ConvertT: %v", slaveID, err)
						} else if temp, err := device.ReadTemperature(); err != nil {
							log.Errorf("Slave %v: ReadTemperature: %v", slaveID, err)
						} else {
							log.Infof("DS18B20 %v: Temperature=%v", slaveID, temp)
						}

					default:
						log.Warnf("Slave %v: unknown family=%02x", slaveID, slaveID.Family)
					}
				}
			}
		}
	}

	return nil
}

func onEvent(event w1.Event) error {
	switch event.Type {
	case w1.MsgTypeSlaveAdd:
		log.Infof("Add Slave: %v", event.SlaveID())
	case w1.MsgTypeSlaveRemove:
		log.Infof("Remove Slave: %v", event.SlaveID())
	case w1.MsgTypeMasterAdd:
		log.Infof("Add Master: %v", event.MasterID())
	case w1.MsgTypeMasterRemove:
		log.Infof("Remove Master: %v", event.MasterID())
	}

	return nil
}

func listen(w1conn *w1.Conn) error {
	if err := w1conn.Listen(onEvent); err != nil {
		return fmt.Errorf("w1 Listen: %v", err)
	}

	return nil
}

func run() error {
	w1conn, err := w1.Dial()
	if err != nil {
		return fmt.Errorf("w1.Dial: %v", err)
	} else {
		log.Infof("Connected to w1-netlink: %v", w1conn)
	}

	if err := scan(w1conn); err != nil {
		return err
	}

	if err := listen(w1conn); err != nil {
		return err
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
