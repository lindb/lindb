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

package fasttime

import (
	"sync/atomic"
	"time"
)

var (
	unixNano         = time.Now().UnixNano()
	unixMicroSeconds = time.Now().UnixNano() / 1e3
	unixMilliseconds = time.Now().UnixNano() / 1e6
	unixTimestamp    = time.Now().UnixNano() / 1e9
)

func init() {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 5)
		defer ticker.Stop()
		for t := range ticker.C {
			n := t.UnixNano()
			atomic.StoreInt64(&unixNano, n)
			atomic.StoreInt64(&unixMicroSeconds, n/1e3)
			atomic.StoreInt64(&unixMilliseconds, n/1e6)
			atomic.StoreInt64(&unixTimestamp, n/1e9)
		}
	}()
}

// UnixNano returns approximate nanoseconds,
func UnixNano() int64 {
	return atomic.LoadInt64(&unixNano)
}

// UnixMicroseconds returns approximate microseconds
func UnixMicroseconds() int64 {
	return atomic.LoadInt64(&unixMicroSeconds)
}

// UnixMilliseconds returns approximate milliseconds
func UnixMilliseconds() int64 {
	return atomic.LoadInt64(&unixMilliseconds)
}

// UnixTimestamp returns approximate timestamp in second
func UnixTimestamp() int64 {
	return atomic.LoadInt64(&unixTimestamp)
}
