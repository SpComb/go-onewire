package w1

import (
	"github.com/SpComb/iot-poc/netlink/connector"

	"fmt"
)

func Dial() (*Conn, error) {
	var conn Conn

	if connectorConn, err := connector.Dial(); err != nil {
		return nil, err
	} else {
		conn.connectorConn = connectorConn
	}

	return &conn, nil
}

type Conn struct {
	connectorConn *connector.Conn
}

func (c *Conn) Send(msg Message) error {
	log.Debugf("Send: %#v", msg)

	var connectorMsg = connector.Message{
		Header: connector.Header{
			ID: W1_ID,
		},
	}

	if data, err := msg.MarshalBinary(); err != nil {
		return err
	} else {
		connectorMsg.Data = data
	}

	return c.connectorConn.Send(connectorMsg)
}

func (c *Conn) Receive() ([]Message, error) {
	connectorMessages, err := c.connectorConn.Receive()
	if err != nil {
		return nil, err
	}

	var msgs = make([]Message, len(connectorMessages))
	for i, connectorMessage := range connectorMessages {
		var msg Message

		if err := msg.UnmarshalBinary(connectorMessage.Data); err != nil {
			return nil, err
		}

		log.Debugf("Receive: %#v", msg)

		msgs[i] = msg
	}

	return msgs, nil
}

func (c *Conn) Request(msg Message) ([]Message, error) {
	if err := c.Send(msg); err != nil {
		return nil, fmt.Errorf("Send: %v", err)
	}

	if msgs, err := c.Receive(); err != nil {
		return nil, fmt.Errorf("Receive: %v", err)
	} else {
		return msgs, err
	}
}

func (c *Conn) ListMasters() (MasterList, error) {
	var msg = Message{
		Header: Header{
			Type: MsgTypeListMasters,
		},
	}

	msgs, err := c.Request(msg)
	if err != nil {
		return nil, err
	}

	var masterList MasterList

	for _, msg := range msgs {
		if err := masterList.UnmarshalBinary(msg.Data); err != nil {
			return masterList, fmt.Errorf("Unmarshal %T: %v", masterList, err)
		}

		log.Debugf("ListMasters: %#v", masterList)
	}

	return masterList, nil
}

func (c *Conn) Close() error {
	return c.connectorConn.Close()
}
