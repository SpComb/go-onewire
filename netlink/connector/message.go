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

func (m *Message) MarshalBinary() ([]byte, error) {
	m.Header.Len = uint16(len(m.Data))

	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, m.Header); err != nil {
		return nil, err
	}

	if _, err := buf.Write(m.Data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (m *Message) UnmarshalBinary(data []byte) error {
	var buf = bytes.NewReader(data)

	if err := binary.Read(buf, byteOrder, &m.Header); err != nil {
		return err
	}

	m.Data = make([]byte, m.Header.Len)

	if read, err := buf.Read(m.Data); err != nil {
		return err
	} else if read != len(m.Data) {
		return io.EOF
	}

	return nil
}
