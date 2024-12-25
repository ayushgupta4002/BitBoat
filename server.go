package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/ayushgupta4002/bitboat/cache"
)

type ServerOpts struct {
	listenAddr string
	isAdmin    bool
	adminAddr  string
}

type Server struct {
	ServerOpts
	cache     cache.Cacher
	followers map[net.Conn]struct{}
}

func NEWServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		followers:  make(map[net.Conn]struct{}),
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer l.Close()
	if !s.isAdmin {
		go func() {
			connection, err := net.Dial("tcp", s.adminAddr)
			if err != nil {
				fmt.Printf("failed to listen: %v", err)
				return
			}
			fmt.Println("Connected to admin server", s.adminAddr)
			go s.handleConn(connection)

		}()
	}
	fmt.Printf("server listening on %s\n", s.listenAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("failed to accept: %v", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	if s.isAdmin {
		s.followers[conn] = struct{}{}
	}

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				// Gracefully handle the client closing the connection
				fmt.Println("Client closed the connection.")
				break
			}
			fmt.Printf("failed to read: %v\n", err)
			break
		}
		if n == 0 {
			return
		}
		fmt.Printf("received: %s\n", string(buf[:n]))

		s.handleCommand(conn, buf[:n])
	}
}

func (s *Server) handleCommand(conn net.Conn, rawMsg []byte) {
	msg, err := s.commandParser(rawMsg)
	if err != nil {
		fmt.Println("Error parsing command:", err)
		conn.Write([]byte(err.Error()))
		return
	}
	switch msg.cmd {
	case CmdSet:
		fmt.Println("SET command")
		err = s.handleSet(conn, msg)

	case CmdGet:
		fmt.Println("GET command")
		err = s.handleGet(conn, msg)
	}
	if err != nil {
		conn.Write([]byte(err.Error()))
	}
}
func (s *Server) handleGet(conn net.Conn, msg *Message) error {
	value, err := s.cache.Get(msg.key)
	if err != nil {
		return err
	}
	_, err = conn.Write([]byte(value))

	return err
}
func (s *Server) handleSet(conn net.Conn, msg *Message) error {
	err := s.cache.Set(msg.key, msg.value, msg.ttl)
	if err != nil {
		return err
	}
	s.sendFollowers(context.TODO(), msg)
	return nil
}

func (s *Server) sendFollowers(ctx context.Context, msg *Message) error {
	fmt.Println("Sending message to followers", s.followers)
	for conn := range s.followers {

		rawMsg := msg.toBytes()
		_, err := conn.Write(rawMsg)
		if err != nil {
			fmt.Println("Error sending message to follower:", err)
			continue
		}
	}
	return nil
}

func (s *Server) commandParser(msg []byte) (*Message, error) {
	var msgStr = string(msg)
	var parts = strings.Split(msgStr, " ")
	fmt.Println("parts", parts)

	if len(parts) < 2 {
		return nil, errors.New("invalid command")
	}
	msgStruct := &Message{
		cmd: Command(parts[0]),
		key: []byte(parts[1]),
	}
	if msgStruct.cmd == CmdSet {
		if len(parts) < 4 {
			return nil, errors.New("invalid SET command")
		}
		ttl, err := strconv.Atoi(parts[3])
		if err != nil {
			return nil, errors.New("invalid TTL")
		}
		msgStruct.value = []byte(parts[2])
		msgStruct.ttl = time.Duration(ttl)
	}

	return msgStruct, nil
}
