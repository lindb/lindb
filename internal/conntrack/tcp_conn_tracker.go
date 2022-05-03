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
	"io"
	"net"

	"github.com/lindb/lindb/metrics"
)

//go:generate mockgen  -destination=./conn_mock.go -package=conntrack net Conn

// TrackedConn tracks a net.Conn with linmetric.
type TrackedConn struct {
	net.Conn
	statistics *metrics.ConnStatistics
}

func (tc *TrackedConn) Read(p []byte) (int, error) {
	n, err := tc.Conn.Read(p)
	tc.statistics.Read.Incr()
	tc.statistics.ReadBytes.Add(float64(n))
	if err != nil && err != io.EOF {
		tc.statistics.ReadFailures.Incr()
	}
	return n, err
}

func (tc *TrackedConn) Write(p []byte) (int, error) {
	n, err := tc.Conn.Write(p)
	tc.statistics.Write.Incr()
	tc.statistics.WriteBytes.Add(float64(n))
	if err != nil {
		tc.statistics.WriteFailures.Incr()
	}
	return n, err
}

func (tc *TrackedConn) Close() error {
	err := tc.Conn.Close()
	tc.statistics.Close.Incr()
	if err != nil {
		tc.statistics.Close.Incr()
	}
	tc.statistics.ActiveConn.Decr()
	return err
}
