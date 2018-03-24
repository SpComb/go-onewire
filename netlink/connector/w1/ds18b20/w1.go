package ds18b20

import (
	"github.com/SpComb/go-onewire/netlink/connector/w1"

	"fmt"
	"io"
)

const Family w1.Family = 0x28

type Cmd uint8

const (
	CmdConvertT        Cmd = 0x44
	CmdWriteScratchpad     = 0x4E
	CmdReadScratchpad      = 0xBE
)

const scratchpadSize = 9

type Scratchpad struct {
	Temperature Temperature
	TempH       uint8
	TempL       uint8
	Config      uint8
	_           uint8
	_           uint8
	_           uint8
	CRC         uint8
}

func (s *Scratchpad) unpack(data []byte) error {
	if len(data) != scratchpadSize {
		return io.EOF
	}

	if !CheckCRC(data) {
		return fmt.Errorf("CRC check failed")
	}

	s.Temperature = unpackTemperature(data[0], data[1])
	s.TempH = data[2]
	s.TempL = data[3]
	s.Config = data[4]
	s.CRC = data[8]

	return nil
}
