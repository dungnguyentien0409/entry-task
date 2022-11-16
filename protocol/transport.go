package protocol

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"time"
)

type Transport struct {
	conn       net.Conn
	bufferSize int
}

func NewTransport(conn net.Conn, bufferSize int) *Transport {
	return &Transport{
		conn,
		bufferSize,
	}
}

func (t *Transport) Send(data []byte) error {
	buf := make([]byte, len(data)+4)
	binary.BigEndian.PutUint32(buf, uint32(len(data)))
	copy(buf[4:], data)
	_, err := t.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (t *Transport) Read() ([]byte, error) {
	bf := bufio.NewReaderSize(t.conn, t.bufferSize)
	header := make([]byte, 4)
	_, err := io.ReadFull(bf, header)
	if err != nil {
		return nil, err
	}

	len := binary.BigEndian.Uint32(header)
	data := make([]byte, len)
	_, err = io.ReadFull(bf, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type Listener interface {
	Init(addr string) error
	Accept() (conn net.Conn, err error)
	Close() error
	ListenAddr() net.Addr
}

type defaultListener struct {
	L net.Listener
}

func (ln *defaultListener) Init(addr string) (err error) {
	ln.L, err = net.Listen("tcp", addr)
	return
}

func (ln *defaultListener) ListenAddr() net.Addr {
	if ln.L != nil {
		return ln.L.Addr()
	}
	return nil
}

func (ln *defaultListener) Accept() (conn net.Conn, err error) {
	c, err := ln.L.Accept()
	if err != nil {
		return nil, err
	}
	if err = setupKeepalive(c); err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

func (ln *defaultListener) Close() error {
	return ln.L.Close()
}

func setupKeepalive(conn net.Conn) error {
	tcpConn := conn.(*net.TCPConn)
	if err := tcpConn.SetKeepAlive(true); err != nil {
		return err
	}
	if err := tcpConn.SetKeepAlivePeriod(60 * time.Second); err != nil {
		return err
	}
	return nil
}
