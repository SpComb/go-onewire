package ds18b20

import (
	"fmt"
)

func unpackTemperature(lsb byte, msb byte) Temperature {
	return Temperature((int16(msb) << 8) | int16(lsb))
}

type Temperature int16

func (t Temperature) Float32() float32 {
	return float32(t) / 16.0
}

func (t Temperature) String() string {
	return fmt.Sprintf("%.3f", t.Float32())
}
