package main

import (
	"fmt"
	"net"

	"github.com/ayushgupta4002/bitboat/cache"
)

type ServerOpts struct {
	listenAddr string
	isAdmin    bool
}

type Server struct {
	ServerOpts
	cache cache.Cacher
}

func NEWServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer l.Close()
	fmt.Printf("server listening on %s\n", s.listenAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Errorf("failed to accept: %v", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Errorf("failed to read: %v", err)
			return
		}
		if n == 0 {
			return
		}
		fmt.Printf("received: %s\n", string(buf[:n]))
	}
}
