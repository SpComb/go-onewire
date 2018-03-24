package connector

import (
	"bytes"
	"encoding/binary"
	"io"
)

type Header struct {
	ID    ID
	Seq   uint32
	Ack   uint32
	Len   uint16
	Flags uint16
}

type Message struct {
	Header
	Data []byte
}

func (msg *Message) MarshalBinary() ([]byte, error) {
	msg.Header.Len = uint16(len(msg.Data))

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
