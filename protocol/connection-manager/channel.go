package connection

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type channel struct {
	mu      sync.RWMutex
	conns   chan net.Conn
	factory Factory
}

func (c *channel) Get() (net.Conn, error) {
	conns, factory := c.getConnAndFactory()
	if conns == nil {
		return nil, ErrClosed
	}

	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}

		return c.WrapConnection(conn), nil
	default:
		conn, err := factory()
		if err != nil {
			return nil, err
		}

		return c.WrapConnection(conn), nil
	}
}

func (c *channel) Close() {
	c.mu.Lock()
	conns := c.conns
	c.conns = nil
	c.factory = nil
	c.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

func (c *channel) Put(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.conns == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}

	select {
	case c.conns <- conn:
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

func (c *channel) Len() int {
	conns, _ := c.getConnAndFactory()
	return len(conns)
}

// Factory is a function to create new connections.
type Factory func() (net.Conn, error)

func NewChannelPool(initialCap, maxCap int, factory Factory) (IChannel, error) {
	if initialCap < 0 || maxCap <= 0 || initialCap > maxCap {
		return nil, errors.New("invalid capacity settings")
	}

	c := &channel{
		conns:   make(chan net.Conn, maxCap),
		factory: factory,
	}

	for i := 0; i < initialCap; i++ {
		conn, err := factory()
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}

		c.conns <- conn
	}

	return c, nil
}

func (c *channel) getConnAndFactory() (chan net.Conn, Factory) {
	c.mu.RLock()
	conns := c.conns
	factory := c.factory
	c.mu.RLock()
	return conns, factory
}

func (c *channel) WrapConnection(conn net.Conn) net.Conn {
	p := &PoolConn{c: c}
	p.Conn = conn
	return p
}
