package w1

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Family uint8

type SlaveID struct {
	Family Family
	Serial [6]uint8
	CRC    uint8
}

func (id SlaveID) String() string {
	return fmt.Sprintf("%02x-%02x%02x%02x%02x%02x%02x",
		id.Family,
		id.Serial[0],
		id.Serial[1],
		id.Serial[2],
		id.Serial[3],
		id.Serial[4],
		id.Serial[5],
	)
}

func (id SlaveID) Pack() MessageID {
	return MessageID{
		byte(id.Family),
		id.Serial[5],
		id.Serial[4],
		id.Serial[3],
		id.Serial[2],
		id.Serial[1],
		id.Serial[0],
		id.CRC,
	}
}

func (id *SlaveID) Unpack(msgID MessageID) error {
	id.Family = Family(msgID[0])
	id.Serial[0] = msgID[6]
	id.Serial[1] = msgID[5]
	id.Serial[2] = msgID[4]
	id.Serial[3] = msgID[3]
	id.Serial[4] = msgID[2]
	id.Serial[5] = msgID[1]
	id.CRC = msgID[7]

	// TODO: check CRC
	return nil
}

type SlaveList []SlaveID

func (l *SlaveList) UnmarshalBinary(data []byte) error {
	var buf = bytes.NewReader(data)

	for buf.Len() > 0 {
		var slaveID SlaveID
		var messageID MessageID

		if err := binary.Read(buf, byteOrder, &messageID); err != nil {
			return err
		}

		if err := slaveID.Unpack(messageID); err != nil {
			return err
		}

		*l = append(*l, slaveID)
	}

	return nil
}
