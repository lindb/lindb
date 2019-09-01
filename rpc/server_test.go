package rpc

import (
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestNewTCPServer(t *testing.T) {
	s := NewTCPServer(":111111111", nil)
	if s.Start() == nil {
		t.Fatal("should be error")
	}

	ctl := gomock.NewController(t)
	mockTCPHandler := NewMockTCPHandler(ctl)
	mockTCPHandler.EXPECT().Handle(gomock.Any()).Return(errors.New("mock errors"))

	s = NewTCPServer(":9000", mockTCPHandler)

	go func() {
		if err := s.Start(); err != nil {
			fmt.Printf("tcp server start err:%v", err)
		}
	}()

	// wait to server to start
	time.Sleep(time.Millisecond * 20)

	conn, err := net.Dial("tcp", ":9000")
	assert.Nil(t, err)

	// wait for server to handler
	time.Sleep(time.Millisecond * 20)
	err = conn.Close()
	assert.Nil(t, err)

	s.Stop()
}
