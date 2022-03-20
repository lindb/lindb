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
)

type connStatistics struct {
	acceptCounter *linmetric.BoundCounter
	acceptErrors  *linmetric.BoundCounter
	connNum       *linmetric.BoundGauge
	readCounter   *linmetric.BoundCounter
	readBytes     *linmetric.BoundCounter
	readErrors    *linmetric.BoundCounter
	writeCounter  *linmetric.BoundCounter
	writeBytes    *linmetric.BoundCounter
	writeErrors   *linmetric.BoundCounter
	closeCounter  *linmetric.BoundCounter
	closeErrors   *linmetric.BoundCounter
}

type TrackedListener struct {
	net.Listener
	connStatistics connStatistics
}

// NewTrackedListener returns new tracked TCP listener for the given addr.
func NewTrackedListener(network, addr string, r *linmetric.Registry) (*TrackedListener, error) {
	ln, err := net.Listen(network, addr)
	if err != nil {
		return nil, err
	}

	tcpScope := r.NewScope("lindb.traffic.tcp", "addr", addr)
	return &TrackedListener{
		Listener: ln,
		connStatistics: connStatistics{
			acceptCounter: tcpScope.NewCounter("accept_conns"),
			acceptErrors:  tcpScope.NewCounter("accept_errors"),
			connNum:       tcpScope.NewGauge("conns_num"),
			readCounter:   tcpScope.NewCounter("read_count"),
			readBytes:     tcpScope.NewCounter("read_bytes"),
			readErrors:    tcpScope.NewCounter("read_errors"),
			writeCounter:  tcpScope.NewCounter("write_count"),
			writeBytes:    tcpScope.NewCounter("write_bytes"),
			writeErrors:   tcpScope.NewCounter("write_errors"),
			closeCounter:  tcpScope.NewCounter("close_conns"),
			closeErrors:   tcpScope.NewCounter("close_errors"),
		},
	}, nil
}

// Accept wraps the accept method with listener statistics
func (tl *TrackedListener) Accept() (net.Conn, error) {
	for {
		conn, err := tl.Listener.Accept()
		tl.connStatistics.acceptCounter.Incr()
		if err != nil {
			var ne net.Error
			if errors.As(err, &ne) && ne.Temporary() {
				time.Sleep(time.Millisecond * 100)
				continue
			}
			tl.connStatistics.acceptErrors.Incr()
			return nil, err
		}
		tl.connStatistics.connNum.Incr()
		tc := &TrackedConn{Conn: conn, statistics: &tl.connStatistics}
		return tc, nil
	}
}
