package ds18b20

import (
	"github.com/snksoft/crc"
)

var crcParameters = crc.Parameters{
	Width:      8,
	Polynomial: 0x31,
	ReflectIn:  true,
	ReflectOut: true,
}

func CheckCRC(data []byte) bool {
	crc8 := crc.CalculateCRC(&crcParameters, data)

	return crc8 == 0
}
