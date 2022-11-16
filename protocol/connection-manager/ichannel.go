package connection

import (
	"errors"
	"net"
)

var (
	ErrClosed = errors.New("pool is closed")
)

type IChannel interface {
	Get() (net.Conn, error)
	Put(conn net.Conn) error
	Close()
	Len() int
}
