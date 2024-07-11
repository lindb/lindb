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

import (
	"unsafe"
)

func U32SliceToBytes(u []uint32) []byte {
	if len(u) == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(u))), len(u)*4)
}

func BytesToU32Slice(b []byte) []uint32 {
	if len(b) == 0 {
		return nil
	}
	return unsafe.Slice((*uint32)(unsafe.Pointer(unsafe.SliceData(b))), len(b)/4)
}

func U64SliceToBytes(u []uint64) []byte {
	if len(u) == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(unsafe.SliceData(u))), len(u)*8)
}

func BytesToU64Slice(b []byte) []uint64 {
	if len(b) == 0 {
		return nil
	}
	return unsafe.Slice((*uint64)(unsafe.Pointer(unsafe.SliceData(b))), len(b)/8)
}

func Float64ToBytes(f64 float64) []byte {
	return unsafe.Slice((*byte)(unsafe.Pointer(&f64)), 8)
}

func BytesToFloat64(b []byte) float64 {
	return unsafe.Slice((*float64)(unsafe.Pointer(unsafe.SliceData(b))), 1)[0]
}
