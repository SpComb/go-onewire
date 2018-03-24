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

type Cmd struct {
	Header struct {
		Cmd CmdType
		_   uint8
		Len uint16
	}
	Data []byte
}

func (cmd *Cmd) MarshalBinary() ([]byte, error) {
	cmd.Header.Len = uint16(len(cmd.Data))

	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, cmd.Header); err != nil {
		return nil, err
	}

	if _, err := buf.Write(cmd.Data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (cmd *Cmd) UnmarshalBinary(data []byte) error {
	var buf = bytes.NewReader(data)

	if err := binary.Read(buf, byteOrder, &cmd.Header); err != nil {
		return err
	}

	cmd.Data = make([]byte, cmd.Header.Len)

	if read, err := buf.Read(cmd.Data); err != nil {
		return err
	} else if read != len(cmd.Data) {
		return io.EOF
	}

	return nil
}
