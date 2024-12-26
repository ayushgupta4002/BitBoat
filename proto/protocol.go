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
)

type CommandSet struct {
	Key   []byte
	Value []byte
	TTL   int32
}
type CommandGet struct {
	Key []byte
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
