package protocol

import (
	"fmt"
	"io"
	"log"
	"net"
)

type HandlerFunc func(request []byte) (response []byte, err error)

type Server struct {
	Addr       string
	BufferSize int
	Handler    HandlerFunc
	Listener   Listener
}

const (
	DEFAULT_BUFFER_SIZE = 64 * 1024
)

func (s *Server) Start() error {

	if s.Handler == nil {
		panic("Server: Server.Handler cannot be nil")
	}

	if s.BufferSize <= 0 {
		s.BufferSize = DEFAULT_BUFFER_SIZE
	}

	if s.Listener == nil {
		s.Listener = &defaultListener{}
	}

	if err := s.Listener.Init(s.Addr); err != nil {
		err = fmt.Errorf("gorpc.Server: [%s]. Cannot listen to: [%s]", s.Addr, err)
		log.Printf("%s", err)
		return err
	}

	defer s.Listener.Close()

	var conn net.Conn
	var err error

	for {
		if conn, err = s.Listener.Accept(); err != nil {
			log.Printf("Server: [%s]. Cannot accept new connection: [%s]", s.Addr, err)
		}

		go handleConnection(s, conn)
	}

	return nil
}

func handleConnection(s *Server, conn net.Conn) {
	transport := NewTransport(conn, s.BufferSize)
	for {
		data, err := transport.Read()
		if err != nil {
			if err != io.EOF {
				log.Printf("read err: %v\n", err)
				return
			}
		}

		if len(data) == 0 {
			return
		}

		resp, err := s.Handler(data)
		if err != nil {
			log.Printf("Error while handle request: %+v", err)
			return
		}

		err = transport.Send(resp)
		if err != nil {
			log.Printf("Error while sending back response to client: %+v", err)
		}
	}
}
