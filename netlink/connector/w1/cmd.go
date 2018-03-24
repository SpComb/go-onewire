package w1

import (
	"bytes"
	"encoding/binary"
	"io"
)

type CmdType uint8

const (
	CmdRead        CmdType = 0
	CmdWrite               = 1
	CmdSearch              = 2
	CmdAlarmSearch         = 3
	CmdTouch               = 4
	CmdReset               = 5
	CmdSlaveAdd            = 6
	CmdSlaveRemove         = 7
	CmdListSlaves          = 8
)

type CmdHeader struct {
	Cmd CmdType
	_   uint8
	Len uint16
}

type Cmd struct {
	CmdHeader
	Data []byte
}

func (cmd *Cmd) marshalBinary(buf *bytes.Buffer) error {
	cmd.Len = uint16(len(cmd.Data))

	if err := binary.Write(buf, byteOrder, cmd.CmdHeader); err != nil {
		return err
	}

	if _, err := buf.Write(cmd.Data); err != nil {
		return err
	}

	return nil
}

func (cmd *Cmd) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer

	if err := cmd.marshalBinary(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (cmd *Cmd) unmarshalBinary(reader io.Reader) error {
	if err := binary.Read(reader, byteOrder, &cmd.CmdHeader); err != nil {
		return err
	}

	cmd.Data = make([]byte, cmd.Len)

	if cmd.Len == 0 {

	} else if read, err := reader.Read(cmd.Data); err != nil {
		return err
	} else if read != len(cmd.Data) {
		return io.EOF
	}

	return nil
}

func (cmd *Cmd) UnmarshalBinary(data []byte) error {
	var reader = bytes.NewReader(data)

	if err := cmd.unmarshalBinary(reader); err != nil {
		return err
	}

	return nil
}

func MarshalCmd(cmdList ...Cmd) ([]byte, error) {
	var buf bytes.Buffer

	for _, cmd := range cmdList {
		if err := cmd.marshalBinary(&buf); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func UnmarshalCmdList(data []byte) ([]Cmd, error) {
	var reader = bytes.NewReader(data)
	var cmdList []Cmd

	for reader.Len() > 0 {
		var cmd Cmd

		if err := cmd.unmarshalBinary(reader); err != nil {
			return cmdList, err
		}

		cmdList = append(cmdList, cmd)
	}

	return cmdList, nil

}
