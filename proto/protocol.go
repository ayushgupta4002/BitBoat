package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type Command byte

const (
	Cmdnonce Command = iota
	CmdSet
	CmdGet
	CmdDel
	CmdJoin
)

type CommandJoin struct{}
type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int32
}
type CommandGet struct {
	Key []byte
}

type Status byte

const (
	StatusNone Status = iota
	StatusOK
	StatusErr
	StatusKeyNotFound
)

func (s Status) Normalize() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusErr:
		return "ERR"
	case StatusKeyNotFound:
		return "KEY NOT FOUND"
	default:
		return "INVALID STATUS"
	}
}

type ResponseSet struct {
	Status Status
}

type ResponseGet struct {
	Status Status
	Value  []byte
}

func (c *ResponseSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, c.Status)
	return buf.Bytes()
}

func (c *ResponseGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, c.Status)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)
	return buf.Bytes()
}

func ParseResponseSet(r io.Reader) (*ResponseSet, error) {
	resp := &ResponseSet{}
	binary.Read(r, binary.LittleEndian, &resp.Status)
	return resp, nil
}

func ParseResponseGet(r io.Reader) (*ResponseGet, error) {
	resp := &ResponseGet{}
	binary.Read(r, binary.LittleEndian, &resp.Status)
	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	value := make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &value)
	resp.Value = value
	return resp, nil
}

func (c *CommandSet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdSet)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Value)))
	binary.Write(buf, binary.LittleEndian, c.Value)
	binary.Write(buf, binary.LittleEndian, int32(c.TTL))
	return buf.Bytes()
}

func (c *CommandGet) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGet)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)
	return buf.Bytes()
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	err := binary.Read(r, binary.LittleEndian, &cmd)
	if err != nil {
		return nil, err
	}
	switch cmd {
	case CmdSet:
		return ParseSet(r), nil
	case CmdGet:
		return ParseGet(r), nil
	case CmdJoin:
		return &CommandJoin{}, nil
	default:
		return nil, fmt.Errorf("invalid command: %d", cmd)
	}
}

func ParseSet(r io.Reader) *CommandSet {
	cmd := &CommandSet{}
	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	key := make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &key)
	var valueLen int32
	binary.Read(r, binary.LittleEndian, &valueLen)
	value := make([]byte, valueLen)
	binary.Read(r, binary.LittleEndian, &value)
	var ttl int32
	binary.Read(r, binary.LittleEndian, &ttl)
	cmd.Key = key
	cmd.Value = value
	cmd.TTL = ttl
	return cmd

}

func ParseGet(r io.Reader) *CommandGet {
	cmd := &CommandGet{}
	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	key := make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &key)
	cmd.Key = key
	return cmd
}
