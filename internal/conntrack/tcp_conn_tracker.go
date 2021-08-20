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
)

// TrackedConn tracks a net.Conn with linmetric
type TrackedConn struct {
	net.Conn
	statistics *connStatistics
}

func (tc *TrackedConn) Read(p []byte) (int, error) {
	n, err := tc.Conn.Read(p)
	tc.statistics.readCounter.Incr()
	tc.statistics.readBytes.Add(float64(n))
	if err != nil && err != io.EOF {
		tc.statistics.readErrors.Incr()
	}
	return n, err
}

func (tc *TrackedConn) Write(p []byte) (int, error) {
	n, err := tc.Conn.Write(p)
	tc.statistics.writeCounter.Incr()
	tc.statistics.writeBytes.Add(float64(n))
	if err != nil {
		tc.statistics.writeErrors.Incr()
	}
	return n, err
}

func (tc *TrackedConn) Close() error {
	tc.statistics.closeCounter.Incr()
	err := tc.Conn.Close()
	if err != nil {
		tc.statistics.closeErrors.Incr()
	}
	tc.statistics.connNum.Decr()
	return err
}
