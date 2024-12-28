package client

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ayushgupta4002/bitboat/proto"
)

type Client struct {
	conn net.Conn
}
type ClientOpts struct {
}

func NewConn(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}

}

func NewClient(addr string, opts ClientOpts) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn: conn,
	}, nil
}

func (c *Client) Set(ctx context.Context, key []byte, value []byte, ttl int32) error {
	cmd := &proto.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}
	resp, err := proto.ParseResponseSet(c.conn)
	if err != nil {
		return err
	}
	if resp.Status != proto.StatusOK {
		return fmt.Errorf("server responsed with non OK status [%s]", resp.Status.Normalize())
	}

	return nil
}

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &proto.CommandGet{
		Key: key,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}
	resp, err := proto.ParseResponseGet(c.conn)
	if resp.Status != proto.StatusOK {
		return nil, fmt.Errorf("server responsed with non OK status [%s]", resp.Status.Normalize())
	}

	log.Printf("GET Status", resp.Status)
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
