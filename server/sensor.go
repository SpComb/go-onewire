package server

import (
	"github.com/SpComb/go-onewire/api"
	"github.com/SpComb/go-onewire/netlink/connector/w1"
	"github.com/SpComb/go-onewire/netlink/connector/w1/ds18b20"

	"time"
)

func newSensor(s *Server, id w1.SlaveID) *Sensor {
	var device = ds18b20.MakeDevice(s.queryConn, id)

	return &Sensor{
		id:     api.SensorID(id.String()),
		device: &device,
	}
}

type Sensor struct {
	id     api.SensorID
	device *ds18b20.Device
	temp   ds18b20.Temperature
	err    error
	at     time.Time
}

func (sensor *Sensor) refresh() {
	sensor.at = time.Now()

	if s, err := sensor.device.Read(); err != nil {
		log.Warnf("Sensor %v: Read: %v", sensor, err)

		sensor.err = err
	} else if err := sensor.device.ConvertT(); err != nil {
		log.Warnf("Sensor %v: ConvertT: %v", sensor, err)

		sensor.err = err
	} else {
		sensor.err = nil
		sensor.temp = s.Temperature

		log.Infof("Sensor %v: Temperature=%v", sensor, sensor.temp)
	}
}

func (sensor *Sensor) MakeConfig() api.SensorConfig {
	return api.SensorConfig{
		Bus:    "onewire",
		Type:   "ds18b20",
		Serial: sensor.device.String(),
	}
}

func (sensor *Sensor) MakeStatus() api.SensorStatus {
	var status = api.SensorStatus{
		At: sensor.at,
	}

	if sensor.err != nil {
		var apiError = api.Error{sensor.err}

		status.Error = &apiError
	} else {
		var apiTemperature = api.Temperature(sensor.temp.Float32())

		status.Temperature = &apiTemperature
	}

	return status
}

func (sensor *Sensor) MakeAPI() api.Sensor {
	return api.Sensor{
		ID:     sensor.id,
		Config: sensor.MakeConfig(),
		Status: sensor.MakeStatus(),
	}
}
