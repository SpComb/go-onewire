package w1

import (
	"bytes"
	"encoding/binary"
)

type MasterID uint32

type IDMaster struct {
	ID MasterID
	_  uint32
}

func (id IDMaster) Pack() MessageID {
	var messageID MessageID

	byteOrder.PutUint32(messageID[0:4], uint32(id.ID))
	byteOrder.PutUint32(messageID[4:8], 0)

	return messageID
}

func (id *IDMaster) Unpack(data MessageID) {
	id.ID = MasterID(byteOrder.Uint32(data[0:4]))
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
