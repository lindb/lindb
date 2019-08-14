package rpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTCPServer(t *testing.T) {
	server := NewTCPServer(":9000")
	go func() {
		_ = server.Start()
	}()
	assert.NotNil(t, server.GetServer())
	server1 := NewTCPServer(":9000")

	// wait server start finish
	time.Sleep(10 * time.Millisecond)
	err := server1.Start()
	assert.NotNil(t, err)

	time.Sleep(10 * time.Millisecond)
	server.Stop()

	go func() {
		_ = server1.Start()
	}()
	time.Sleep(10 * time.Millisecond)
	server1.Stop()

	time.Sleep(10 * time.Millisecond)
}
