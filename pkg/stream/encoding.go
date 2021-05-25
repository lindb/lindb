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

package stream

// UvariantSize returns the bytes-size of a uint64 uvariant encoded number.
func UvariantSize(value uint64) int {
	i := 0
	for value >= 0x80 {
		value >>= 7
		i++
	}
	return i + 1
}

// VariantSize returns the bytes-size of a int64 variant encoded number.
func VariantSize(value int64) int {
	ux := uint64(value) << 1
	if value < 0 {
		ux = ^ux
	}
	return UvariantSize(ux)
}
