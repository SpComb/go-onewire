package connector

import (
	"github.com/mdlayher/netlink"

	"fmt"
)

func packMessage(msg *Message) (netlink.Message, error) {
	var netlinkMessage = netlink.Message{
		Header: netlink.Header{
			Type:     netlink.HeaderTypeDone,
			Sequence: msg.Seq,
		},
	}

	if data, err := msg.MarshalBinary(); err != nil {
		return netlinkMessage, fmt.Errorf("Msg.MarshalBinary: %v", err)
	} else {
		netlinkMessage.Data = data
	}

	return netlinkMessage, nil
}

func unpackMessages(netlinkMessages []netlink.Message) ([]Message, error) {
	var msgs = make([]Message, len(netlinkMessages))
	for i, netlinkMessage := range netlinkMessages {
		var msg Message

		if err := msg.UnmarshalBinary(netlinkMessage.Data); err != nil {
			return msgs, err
		}

		msgs[i] = msg
	}

	return msgs, nil
}
