package connpool

import (
	"net"
)

type ConnPool struct {
	connections []*net.Conn
	messageFrom *net.Conn
}

func New() *ConnPool {
	return &ConnPool{connections: make([]*net.Conn, 0)}
}

func (c *ConnPool) AddConnection(conn *net.Conn) {
	c.connections = append(c.connections, conn)
}

func (c *ConnPool) MessageFrom(conn *net.Conn) {
	c.messageFrom = conn
}

func (c *ConnPool) Write(message []byte) (int, error) {
	message = append(message, byte('\n'))
	for _, conn := range c.connections {
		if c.messageFrom == conn {
			continue
		}
		(*conn).Write(message)
	}
	return 0, nil
}
