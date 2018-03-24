package w1

import (
	"bytes"
	"encoding/binary"
)

type MasterID uint32

type IDMaster struct {
	ID uint32
	_  uint32
}

func (id *IDMaster) MarshalBinary() ([]byte, error) {
	return marshalBinary(*id)
}

func (id *IDMaster) UnmarshalBinary(data []byte) error {
	return unmarshalBinary(id, data)
}

type MasterList []MasterID

func (l *MasterList) UnmarshalBinary(data []byte) error {
	var buf = bytes.NewReader(data)

	for buf.Len() > 0 {
		var masterID MasterID

		if err := binary.Read(buf, byteOrder, &masterID); err != nil {
			return err
		}

		*l = append(*l, masterID)
	}

	return nil
}
