package main

import "time"

type Command string

const (
	CmdSet Command = "SET"
	CmdGet Command = "GET"
)

type Message struct {
	key   []byte
	value []byte
	cmd   Command
	ttl   time.Duration
}
