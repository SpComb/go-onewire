package w1

import (
	"bytes"
	"encoding/binary"
	"io"
)

type MsgType uint8

const (
	MsgTypeSlaveAdd     MsgType = 0
	MsgTypeSlaveRemove          = 1
	MsgTypeMasterAdd            = 2
	MsgTypeMasterRemove         = 3
	MsgTypeMasterCmd            = 4
	MsgTypeSlaveCmd             = 5
	MsgTypeListMasters          = 6
)

type MessageID [8]byte

type Header struct {
	Type   MsgType
	Status ErrorStatus
	Len    uint16
	ID     MessageID
}

type Message struct {
	Header
	Data []byte
}

func (msg *Message) MarshalBinary() ([]byte, error) {
	msg.Len = uint16(len(msg.Data))

	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, msg.Header); err != nil {
		return nil, err
	}

	if _, err := buf.Write(msg.Data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (msg *Message) UnmarshalBinary(data []byte) error {
	var buf = bytes.NewReader(data)

	if err := binary.Read(buf, byteOrder, &msg.Header); err != nil {
		return err
	}

	msg.Data = make([]byte, msg.Len)

	if read, err := buf.Read(msg.Data); err != nil {
		return err
	} else if read != len(msg.Data) {
		return io.EOF
	}

	return nil
}
