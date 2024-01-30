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

package strutil

import (
	"strings"
	"unsafe"
)

// GetStringValue aggregation format function name
func GetStringValue(rawString string) string {
	if len(rawString) > 0 {
		if (strings.HasPrefix(rawString, "'") && strings.HasSuffix(rawString, "'")) ||
			(strings.HasPrefix(rawString, "\"") && strings.HasSuffix(rawString, "\"")) {
			return rawString[1 : len(rawString)-1]
		}
		return rawString
	}
	return ""
}

func ByteSlice2String(bytes []byte) string {
	return unsafe.String(&bytes[0], len(bytes))
}

func String2ByteSlice(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

// DeDupStringSlice removes the duplicated string in a list
func DeDupStringSlice(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	var m = make(map[string]struct{})
	for _, item := range items {
		m[item] = struct{}{}
	}
	var dst = make([]string, len(m))
	idx := 0
	for k := range m {
		dst[idx] = k
		idx++
	}
	return dst
}
