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

package point

import "bytes"

type escapeSet struct {
	k   [1]byte
	esc [2]byte
}

var (
	metricNameEscapeCodes = [...]escapeSet{
		{k: [1]byte{','}, esc: [2]byte{'\\', ','}},
		{k: [1]byte{' '}, esc: [2]byte{'\\', ' '}},
	}
)

func UnescapeMetricName(in []byte) []byte {
	if bytes.IndexByte(in, '\\') == -1 {
		return in
	}

	for i := range metricNameEscapeCodes {
		c := &metricNameEscapeCodes[i]
		if bytes.IndexByte(in, c.k[0]) != -1 {
			in = bytes.Replace(in, c.esc[:], c.k[:], -1)
		}
	}
	return in
}

func EscapeMetricName(in []byte) []byte {
	for _, c := range metricNameEscapeCodes {
		if bytes.IndexByte(in, c.k[0]) != -1 {
			in = bytes.Replace(in, c.k[:], c.esc[:], -1)
		}
	}
	return in
}
