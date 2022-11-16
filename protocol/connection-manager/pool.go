package connection

import (
	"net"
	"sync"
)

type PoolConn struct {
	net.Conn
	mu       sync.RWMutex
	c        *channel
	unusable bool
}

func (p *PoolConn) Close() error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.unusable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
		return nil
	}

	return p.c.Put(p.Conn)
}

func (p *PoolConn) MarkUnusable() {
	p.mu.Lock()
	p.unusable = true
	p.mu.Unlock()
}
