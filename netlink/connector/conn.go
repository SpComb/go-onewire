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

	if netlinkMessage, err := packMessage(&msg); err != nil {
		return err
	} else if _, err := c.netlinkConn.Send(netlinkMessage); err != nil {
		return fmt.Errorf("netlink Send: %v", err)
	} else {
		log.Debugf("Send: %#v", msg)
	}

	return nil
}

func (c *Conn) Receive() ([]Message, error) {
	if netlinkMessages, err := c.netlinkConn.Receive(); err != nil {
		return nil, fmt.Errorf("netlink Receive: %v", err)
	} else if msgs, err := unpackMessages(netlinkMessages); err != nil {
		return nil, err
	} else {
		log.Debugf("Recv: %#v", msgs)

		return msgs, nil
	}
}

func (c *Conn) Execute(msg Message) ([]Message, error) {
	if msg.Seq == 0 {
		msg.Seq = c.nextSeq()
	}

	if netlinkMsg, err := packMessage(&msg); err != nil {
		return nil, err
	} else if netlinkMsgs, err := c.netlinkConn.Execute(netlinkMsg); err != nil {
		return nil, fmt.Errorf("netlink Execute %#v: %v", netlinkMsg, err)
	} else if msgs, err := unpackMessages(netlinkMsgs); err != nil {
		return nil, err
	} else {
		log.Debugf("Exchange %#v: %#v", msg, msgs)

		return msgs, nil
	}
}

func (c *Conn) Close() error {
	return c.netlinkConn.Close()
}
