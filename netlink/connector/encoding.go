package connector

import (
	"bytes"
	"encoding/binary"
)

var byteOrder = binary.LittleEndian

func marshalBinary(obj interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if err := binary.Write(&buf, byteOrder, obj); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func unmarshalBinary(obj interface{}, data []byte) error {
	var buf = bytes.NewReader(data)

	if err := binary.Read(buf, byteOrder, obj); err != nil {
		return err
	}

	return nil
}
