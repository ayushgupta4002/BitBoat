package client

import (
	"context"
	"net"

	"github.com/ayushgupta4002/bitboat/proto"
)

type Client struct {
	conn net.Conn
}
type ClientOpts struct {
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

	return err
}

func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &proto.CommandGet{
		Key: key,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 2048)
	_, err = c.conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (c *Client) Close() {
	c.conn.Close()
}
