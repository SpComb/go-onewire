package w1

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type MessageType uint8

const (
	MsgTypeSlaveAdd     MessageType = 0
	MsgTypeSlaveRemove              = 1
	MsgTypeMasterAdd                = 2
	MsgTypeMasterRemove             = 3
	MsgTypeMasterCmd                = 4
	MsgTypeSlaveCmd                 = 5
	MsgTypeListMasters              = 6
)

func (t MessageType) String() string {
	switch t {
	case MsgTypeSlaveAdd:
		return "SlaveAdd"
	case MsgTypeSlaveRemove:
		return "SlaveRemove"
	case MsgTypeMasterAdd:
		return "MasterAdd"
	case MsgTypeMasterRemove:
		return "MasterRemove"
	case MsgTypeMasterCmd:
		return "MasterCmd"
	case MsgTypeSlaveCmd:
		return "SlaveCmd"
	case MsgTypeListMasters:
		return "ListMasters"
	default:
		return fmt.Sprintf("%d", t)
	}
}

type MessageID [8]byte

type Header struct {
	Type   MessageType
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

	if msg.Len == 0 {

	} else if read, err := buf.Read(msg.Data); err != nil {
		return err
	} else if read != len(msg.Data) {
		return io.EOF
	}

	return nil
}
