package protocol

import (
	"entrytask/protocol/connection-manager"
	"log"
	"net"
)

type Client struct {
	connPool connection.IChannel
}

func NewClient(minCap, maxCap int, address string) (*Client, error) {
	p, err := connection.NewChannelPool(minCap, maxCap, func() (net.Conn, error) {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			log.Printf("Error while creating connection %+v", err)
			return nil, err
		}

		err = setupKeepalive(conn)
		if err != nil {
			log.Printf("Error while setup connection keep alive %+v", err)
			conn.Close()
			return nil, err
		}
		return conn, nil
	})

	if err != nil {
		panic(err)
	}

	return &Client{
		connPool: p}, nil
}

func (c *Client) Call(data []byte) ([]byte, error) {
	conn, err := c.connPool.Get()
	if err != nil {
		log.Printf("Error while getting connection from pool %+v", err)
		return nil, err
	}

	transport := NewTransport(conn, DEFAULT_BUFFER_SIZE)
	transport.Send(data)
	rsp, err := transport.Read()

	if err != nil {
		log.Printf("Error while reading response from server %+v", err)
		conn.Close()
		return nil, err
	}

	err = c.connPool.Put(conn)
	if err != nil {
		log.Printf("Error while putting back the connection %+v", err)
		return nil, err
	}

	return rsp, nil
}
