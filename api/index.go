package api

import (
	"time"
)

type Index struct {
	Sensors []Sensor
}

type SensorID string

type Temperature float32

type SensorConfig struct {
	Bus    string
	Type   string
	Serial string
}

type SensorStatus struct {
	At          time.Time
	Temperature *Temperature `json:",omitempty"`
	Error       *Error       `json:",omitempty"`
}

type Sensor struct {
	ID     SensorID
	Config SensorConfig
	Status SensorStatus
}
