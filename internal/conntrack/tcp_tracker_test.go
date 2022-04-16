// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package conntrack

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/linmetric"
)

type testListenerTracker struct {
	httpServer     http.Server
	serverListener net.Listener
}

func (tracker *testListenerTracker) Prepare(t *testing.T) {
	var err error
	tracker.serverListener, err = NewTrackedListener("tcp", ":23424", linmetric.StorageRegistry)
	assert.NoErrorf(t, err, "failed to listen on 23424")
	tracker.httpServer = http.Server{
		Addr: ":23424",
		Handler: http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
			resp.WriteHeader(http.StatusOK)
		}),
	}
	up := make(chan struct{})
	go func() {
		up <- struct{}{}
		_ = tracker.httpServer.Serve(tracker.serverListener)
	}()
	<-up
	time.Sleep(100 * time.Millisecond)
}

func (tracker *testListenerTracker) shutdown() {
	if tracker.serverListener != nil {
		_ = tracker.serverListener.Close()
	}
}

func Test_TrackedListenerTracker(t *testing.T) {
	tracker := &testListenerTracker{}
	defer tracker.shutdown()
	tracker.Prepare(t)

	conn, err := (&net.Dialer{}).DialContext(context.TODO(), "tcp", tracker.serverListener.Addr().String())
	assert.NoError(t, err)
	for i := 0; i < 10; i++ {
		_, err = conn.Write([]byte("hello"))
		assert.NoError(t, err)
	}
	assert.Nil(t, conn.Close())

	time.Sleep(200 * time.Millisecond)
}
