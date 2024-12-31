package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ayushgupta4002/bitboat/cache"
	"github.com/ayushgupta4002/bitboat/client"
	"github.com/ayushgupta4002/bitboat/proto"
)

type ServerOpts struct {
	listenAddr string
	isAdmin    bool
	adminAddr  string
}

type Server struct {
	ServerOpts
	subs       map[*client.Client]struct{}
	cache      cache.Cacher
	leaderConn net.Conn
}

func NEWServer(opts ServerOpts, c cache.Cacher) *Server {
	return &Server{
		ServerOpts: opts,
		cache:      c,
		subs:       make(map[*client.Client]struct{}),
		leaderConn: nil,
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer l.Close()

	if !s.isAdmin && len(s.adminAddr) > 0 {
		go func() {
			err = s.dialAdmin()
			if err != nil {
				log.Println(err)
			}
		}()
	}

	fmt.Printf("server listening on %s\n", s.listenAddr)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("failed to accept: %v", err)
			continue
		}
		fmt.Printf("New client connected: %s\n", conn.RemoteAddr()) // Log new connection

		go s.handleConn(conn)
	}
}

func (s *Server) dialAdmin() error {
	conn, err := net.Dial("tcp", s.adminAddr)
	if err != nil {
		return fmt.Errorf("failed to connect to admin server", err)
	}
	log.Println("connected to leader:", s.adminAddr)
	s.leaderConn = conn
	binary.Write(conn, binary.LittleEndian, proto.CmdJoin)
	s.handleConn(conn)
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	for {
		cmdStruct, err := proto.ParseCommand(conn)
		if err != nil {
			if err.Error() == "EOF" {
				// Gracefully handle the client closing the connection
				fmt.Println("Client closed the connection.")
				break
			}
			fmt.Printf("failed to read: %v\n", err)
			break
		}
		s.handleCommand(conn, cmdStruct)
	}
	fmt.Println("connection closed:", conn.RemoteAddr())

}

func (s *Server) handleCommand(conn net.Conn, cmdStr any) {
	switch cmd := cmdStr.(type) {
	case *proto.CommandSet:
		if len(cmd.Key) == 0 {
			return
		}
		s.handleSet(conn, cmd)
	case *proto.CommandGet_Del_Has:
		if len(cmd.Key) == 0 {
			return
		}
		if cmd.Cmd == proto.CmdDel {
			s.handleDel(conn, cmd)
		} else if cmd.Cmd == proto.CmdHas {
			s.handleHas(conn, cmd)
		} else {
			s.handleGet(conn, cmd)
		}
	case *proto.CommandJoin:
		s.handleJoin(conn, cmd)
	}
}

func (s *Server) handleHas(conn net.Conn, msg *proto.CommandGet_Del_Has) error {
	resp := &proto.ResponseSet_Has_Delete{}
	has, err := s.cache.Has(msg.Key)
	if err != nil {
		return err
	}
	if has {
		resp.Status = proto.StatusOK
	} else {
		resp.Status = proto.StatusKeyNotFound
	}
	_, err = conn.Write(resp.Bytes())
	if err != nil {
		fmt.Println("Error writing response to client:", err)
	}
	return nil
}

func (s *Server) handleDel(conn net.Conn, msg *proto.CommandGet_Del_Has) error {

	resp := &proto.ResponseSet_Has_Delete{}

	go func() {
		for subs := range s.subs {
			err := subs.Delete(context.TODO(), msg.Key)
			if err != nil {
				fmt.Println("Error setting key on follower:", err)
			}
		}
	}()

	err := s.cache.Delete(msg.Key)
	if err != nil {
		return err
	}
	resp.Status = proto.StatusOK

	_, err = conn.Write(resp.Bytes())
	if err != nil {
		fmt.Println("Error writing delete response to client:", err)
	}
	return nil
}

func (s *Server) handleGet(conn net.Conn, msg *proto.CommandGet_Del_Has) error {
	resp := &proto.ResponseGet{}

	value, err := s.cache.Get(msg.Key)
	if err != nil {
		return err
	}
	resp.Status = proto.StatusOK
	resp.Value = value
	_, err = conn.Write(resp.Bytes())
	if err != nil {
		fmt.Println("Error writing response to client:", err)
	}
	return nil
}
func (s *Server) handleSet(conn net.Conn, msg *proto.CommandSet) error {
	resp := &proto.ResponseSet_Has_Delete{}

	if s.isAdmin {
		// Leader logic
		log.Printf("SET %s to %s", msg.Key, msg.Value)

		// Propagate to followers
		go func() {
			for subs := range s.subs {
				err := subs.Set(context.TODO(), msg.Key, msg.Value, msg.TTL)
				if err != nil {
					fmt.Println("Error setting key on follower:", err)
				}
			}
		}()

		// Set in leader's cache
		err := s.cache.Set(msg.Key, msg.Value, time.Duration(msg.TTL))
		if err != nil {
			return err
		}

		resp.Status = proto.StatusOK
	} else {
		// Follower logic
		if !s.isLeaderConn(conn) {
			// Reject if not from leader
			resp.Status = proto.StatusNotLeader

			log.Println("Rejecting SET command: only leader can initiate key updates")
		} else {
			// Accept if from leader
			log.Printf("SET %s to %s follower", msg.Key, msg.Value)

			err := s.cache.Set(msg.Key, msg.Value, time.Duration(msg.TTL))
			if err != nil {
				return err
			}
			resp.Status = proto.StatusOK
		}
	}

	_, err := conn.Write(resp.Bytes())
	if err != nil {
		fmt.Println("Error writing response to client:", err)
	}
	return nil
}

func (s *Server) handleJoin(conn net.Conn, msg *proto.CommandJoin) error {
	fmt.Println("subscriber just joined the cluster:", conn.RemoteAddr())
	s.subs[client.NewConn(conn)] = struct{}{}
	return nil
}

// func (s *Server) commandParser(msg []byte) (*Message, error) {
// 	var msgStr = string(msg)
// 	var parts = strings.Split(msgStr, " ")
// 	fmt.Println("parts", parts)

// 	if len(parts) < 2 {
// 		return nil, errors.New("invalid command")
// 	}
// 	msgStruct := &Message{
// 		cmd: Command(parts[0]),
// 		key: []byte(parts[1]),
// 	}
// 	if msgStruct.cmd == CmdSet {
// 		if len(parts) < 4 {
// 			return nil, errors.New("invalid SET command")
// 		}
// 		ttl, err := strconv.Atoi(parts[3])
// 		if err != nil {
// 			return nil, errors.New("invalid TTL")
// 		}
// 		msgStruct.value = []byte(parts[2])
// 		msgStruct.ttl = time.Duration(ttl)
// 	}

// 	return msgStruct, nil
// }

// func (s *Server) sendFollowers(ctx context.Context, msg *Message) error {
// 	fmt.Println("Sending message to followers", s.followers)
// 	for conn := range s.followers {

// 		rawMsg := msg.toBytes()
// 		_, err := conn.Write(rawMsg)
// 		if err != nil {
// 			fmt.Println("Error sending message to follower:", err)
// 			continue
// 		}
// 	}
// 	return nil
// }

func (s *Server) isLeaderConn(conn net.Conn) bool {
	if s.leaderConn == nil {
		return false
	}
	return conn.RemoteAddr().String() == s.leaderConn.RemoteAddr().String()
}
