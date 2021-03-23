package core

import (
	"net"
	"time"
)

type conn struct {
	net.Conn

	IdleTimeout time.Duration
}

func (c *conn) Write(p []byte) (n int, err error) {
	c.updateDeadline()
	n, err = c.Conn.Write(p)
	return n, err
}

func (c *conn) Read(b []byte) (n int, err error) {
	c.updateDeadline()

	n, err = c.Conn.Read(b)
	return n, err
}

func (c *conn) Close() (err error) {
	err = c.Conn.Close()
	return err
}

func (c *conn) updateDeadline() {
	idleDeadline := time.Now().Add(c.IdleTimeout)
	c.Conn.SetDeadline(idleDeadline)
}
