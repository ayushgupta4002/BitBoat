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
	CmdHas
	CmdJoin
)

type CommandJoin struct{}
type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int32
}
type CommandGet_Del_Has struct {
	Key []byte
	Cmd Command
}

type Status byte

const (
	StatusNone Status = iota
	StatusOK
	StatusErr
	StatusKeyNotFound
	StatusNotLeader
	
)

func (s Status) Normalize() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusErr:
		return "ERR"
	case StatusKeyNotFound:
		return "KEY NOT FOUND"
	case StatusNotLeader:
		return "NOT LEADER"
	default:
		return "INVALID STATUS"
	}
}

type ResponseSet_Has_Delete struct {
	Status Status
}

type ResponseGet struct {
	Status Status
	Value  []byte
}

func (c *ResponseSet_Has_Delete) Bytes() []byte {
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

func ParseResponseSet_Has_Delete(r io.Reader) (*ResponseSet_Has_Delete, error) {
	resp := &ResponseSet_Has_Delete{}
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

func (c *CommandGet_Del_Has) BytesGET() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdGet)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)
	return buf.Bytes()
}

func (c *CommandGet_Del_Has) BytesDEL() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdDel)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)
	return buf.Bytes()
}
func (c *CommandGet_Del_Has) BytesHas() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, CmdHas)
	binary.Write(buf, binary.LittleEndian, int32(len(c.Key)))
	binary.Write(buf, binary.LittleEndian, c.Key)
	return buf.Bytes()
}

func ParseCommand(r io.Reader) (any, error) {
	var cmd Command
	err := binary.Read(r, binary.LittleEndian, &cmd)
	if err != nil {
		// Handle EOF gracefully
		if err == io.EOF {
			return nil, fmt.Errorf("client disconnected: %w", err)
		}
		// Handle any other read error
		return nil, fmt.Errorf("failed to read command: %w", err)
	}

	switch cmd {
	case CmdSet:
		return ParseSet(r), nil
	case CmdGet:
		return ParseGet(r, cmd), nil
	case CmdDel:
		return ParseGet(r, cmd), nil
	case CmdHas:
		return ParseGet(r, cmd), nil
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

func ParseGet(r io.Reader, command Command) *CommandGet_Del_Has {
	cmd := &CommandGet_Del_Has{}
	var keyLen int32
	binary.Read(r, binary.LittleEndian, &keyLen)
	key := make([]byte, keyLen)
	binary.Read(r, binary.LittleEndian, &key)
	cmd.Key = key
	cmd.Cmd = command
	return cmd
}
