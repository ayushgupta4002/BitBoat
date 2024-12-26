package proto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCommand(t *testing.T) {

	cmd := &CommandGet{
		Key: []byte("foo"),
	}
	r := bytes.NewReader(cmd.Bytes())
	CmdGet, err := ParseCommand(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, CmdGet, cmd)
}

func TestSetCommand(t *testing.T) {

	cmd := &CommandSet{
		Key:   []byte("foo"),
		Value: []byte("bar"),
		TTL:   42,
	}
	r := bytes.NewReader(cmd.Bytes())
	CmdSet, err := ParseCommand(r)
	if err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, CmdSet, cmd)

}
