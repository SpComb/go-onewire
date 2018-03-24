package w1

import (
	"github.com/SpComb/go-onewire/netlink/connector"
)

var ConnectorID = connector.ID{Idx: 0x3, Val: 0x1}

func packMessage(msg Message) (connector.Message, error) {
	var connectorMsg = connector.Message{
		Header: connector.Header{
			ID: ConnectorID,
		},
	}

	if data, err := msg.MarshalBinary(); err != nil {
		return connectorMsg, err
	} else {
		connectorMsg.Data = data
	}

	return connectorMsg, nil
}

func unpackMessages(connectorMessages []connector.Message) ([]Message, error) {
	var msgs = make([]Message, len(connectorMessages))
	for i, connectorMessage := range connectorMessages {
		var msg Message

		if err := msg.UnmarshalBinary(connectorMessage.Data); err != nil {
			return nil, err
		}

		msgs[i] = msg
	}

	return msgs, nil
}

func validateMessages(msg Message, msgs []Message) error {
	for _, msg := range msgs {
		if msg.Status != 0 {
			return msg.Status
		}
	}

	return nil
}
