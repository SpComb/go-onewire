package connector

import (
	"github.com/mdlayher/netlink"
	"golang.org/x/sys/unix"

	"fmt"
	"math/rand"
)

type Conn struct {
	netlinkConn *netlink.Conn
	seq         uint32
}

func Dial() (*Conn, error) {
	var conn Conn

	if err := conn.init(); err != nil {
		return nil, err
	}

	return &conn, nil
}

func (c *Conn) init() error {
	var netlinkConfig = netlink.Config{}

	c.seq = rand.Uint32()

	if netlinkConn, err := netlink.Dial(unix.NETLINK_CONNECTOR, &netlinkConfig); err != nil {
		return err
	} else {
		c.netlinkConn = netlinkConn
	}

	return nil
}

func (c *Conn) nextSeq() uint32 {
	c.seq++

	return c.seq
}

func (c *Conn) Send(msg Message) error {
	if msg.Seq == 0 {
		msg.Seq = c.nextSeq()
	}

	log.Debugf("Send: %#v", msg)

	var netlinkMessage = netlink.Message{
		Header: netlink.Header{
			Type:     netlink.HeaderTypeDone,
			Sequence: msg.Seq,
		},
	}

	if data, err := msg.MarshalBinary(); err != nil {
		return fmt.Errorf("Msg.MarshalBinary: %v", err)
	} else {
		netlinkMessage.Data = data
	}

	if _, err := c.netlinkConn.Send(netlinkMessage); err != nil {
		return fmt.Errorf("netlink Send: %v", err)
	}

	return nil
}

func (c *Conn) Receive() ([]Message, error) {
	netlinkMessages, err := c.netlinkConn.Receive()
	if err != nil {
		return nil, fmt.Errorf("netlink Receive: %v", err)
	}

	var msgs = make([]Message, len(netlinkMessages))
	for i, netlinkMessage := range netlinkMessages {
		log.Debugf("netlink Recv: %#v", netlinkMessage)

		var msg Message

		if err := msg.UnmarshalBinary(netlinkMessage.Data); err != nil {
			return nil, err
		}

		log.Debugf("Recv: %#v", msg)

		msgs[i] = msg
	}

	return msgs, nil
}

func (c *Conn) Close() error {
	return c.netlinkConn.Close()
}
