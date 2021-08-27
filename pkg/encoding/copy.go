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

package encoding

import "math"

// MustCopy makes sure that data in src are copied into dst
// If cap of dst < src, a new slice will be returned
// If cap of dst >= src, the slice only needs to be resliced
func MustCopy(dst []byte, src []byte) []byte {
	if cap(dst) < len(src) {
		dst = growSlice(dst, len(src))
	}
	dst = dst[:len(src)]
	copy(dst, src)
	return dst
}

// growSlice grows the byte slice with specified length
// If cap < 1024, cap = cap * 2
// If cap >= 1024, cap >= cap * 1.25
func growSlice(dst []byte, n int) []byte {
	if cap(dst) >= n {
		return dst
	}
	if cap(dst) < 1024 {
		if cap(dst)*2 >= n {
			return make([]byte, cap(dst)*2)
		}
		return make([]byte, n)
	}
	// capacity >= 1024
	targetSize := int(math.Ceil(float64(cap(dst)) * 1.25))
	if targetSize >= n {
		return make([]byte, targetSize)
	}
	return make([]byte, n)
}
