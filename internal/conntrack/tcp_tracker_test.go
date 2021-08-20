package conntrack

import (
	"context"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testListenerTracker struct {
	httpServer     http.Server
	serverListener net.Listener
}

func (tracker *testListenerTracker) Prepare(t *testing.T) {
	var err error
	tracker.serverListener, err = NewTrackedListener("tcp", ":23424")
	assert.NoErrorf(t, err, "failed to listen on 23424")
	tracker.httpServer = http.Server{
		Addr: ":23424",
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			resp.WriteHeader(http.StatusOK)
		}),
	}
	go func() {
		_ = tracker.httpServer.Serve(tracker.serverListener)
	}()
}

func (tracker *testListenerTracker) shutdown() {
	if tracker.serverListener != nil {
		_ = tracker.serverListener.Close()
	}
}

func Test_TrackedListenerTracker(t *testing.T) {
	tracker := &testListenerTracker{}
	tracker.Prepare(t)

	conn, err := (&net.Dialer{}).DialContext(context.TODO(), "tcp", tracker.serverListener.Addr().String())
	assert.NoError(t, err)
	_, err = conn.Write([]byte("hello"))
	assert.NoError(t, err)
	assert.Nil(t, conn.Close())
}
