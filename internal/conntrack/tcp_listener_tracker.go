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
	"errors"
	"net"
	"time"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

//go:generate mockgen  -destination=./listener_mock.go -package=conntrack net Listener

// for testing
var (
	newListenFn = net.Listen
)

// TrackedListener represents net.Listener wrapper for track network statistics.
type TrackedListener struct {
	net.Listener

	statistics *metrics.ConnStatistics
}

// NewTrackedListener returns new tracked TCP listener for the given addr.
func NewTrackedListener(network, addr string, r *linmetric.Registry) (*TrackedListener, error) {
	ln, err := newListenFn(network, addr)
	if err != nil {
		return nil, err
	}

	return &TrackedListener{
		Listener:   ln,
		statistics: metrics.NewConnStatistics(r, addr),
	}, nil
}

// Accept wraps the accept method with listener statistics
func (tl *TrackedListener) Accept() (net.Conn, error) {
	for {
		conn, err := tl.Listener.Accept()
		tl.statistics.Accept.Incr()
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Timeout() {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			tl.statistics.AcceptFailures.Incr()
			return nil, err
		}
		tl.statistics.ActiveConn.Incr()
		tc := &TrackedConn{Conn: conn, statistics: tl.statistics}
		return tc, nil
	}
}
