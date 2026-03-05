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
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

// Create a replacer to handle various SQL escape sequences.
// This handles both Standard SQL (doubling) and MySQL/Postgres (backslash) styles.
var sqlUnescaper = strings.NewReplacer(
	"''", "'", // Standard: 'It''s' -> It's
	`\'`, "'", // MySQL: 'It\'s' -> It's
	"``", "`", // Identifiers: `my``table` -> my`table
	`""`, `"`, // Double Quotes: "" -> "
	`\"`, `"`, // MySQL Double Quote: \" -> "
	`\\`, `\`, // Escaped Backslash: \\ -> \
)

// UnescapeString handles the removal of surrounding quotes and restores escaped characters.
// It follows SQL standards where a double quote of the same type (e.g., ”)
// represents a single literal character.
func UnescapeString(raw string) string {
	// A valid quoted string must have at least 2 characters (the quotes).
	if len(raw) < 2 {
		return raw
	}

	// Identify the quote character used (', ", or `).
	startQuote := raw[0]
	endQuote := raw[len(raw)-1]

	// Ensure the string is correctly bounded by matching quotes.
	if startQuote != endQuote || (startQuote != '\'' && startQuote != '"' && startQuote != '`') {
		return raw
	}

	// Extract the inner content.
	content := raw[1 : len(raw)-1]

	// Apply replacements based on the quote type.
	// We use the replacer to ensure backslashes are consumed correctly.
	return sqlUnescaper.Replace(content)
}

// ByteSlice2String returns the string of the byte slice.
func ByteSlice2String(bytes []byte) string {
	return unsafe.String(&bytes[0], len(bytes))
}

// String2ByteSlice returns the byte slice of the string.
func String2ByteSlice(str string) []byte {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

// DeDupStringSlice removes the duplicated string in a list
func DeDupStringSlice(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	m := make(map[string]struct{})
	for _, item := range items {
		m[item] = struct{}{}
	}
	dst := make([]string, len(m))
	idx := 0
	for k := range m {
		dst[idx] = k
		idx++
	}
	return dst
}

// SliceToTypedString returns the string of the slice.
func SliceToTypedString(slice []any) string {
	if len(slice) == 0 {
		return ""
	}
	var sb strings.Builder
	for _, v := range slice {
		switch val := v.(type) {
		case int:
			sb.WriteString(strconv.Itoa(val))
		case string:
			sb.WriteString(val)
		case bool:
			if val {
				sb.WriteString("1")
			} else {
				sb.WriteString("0")
			}
		default:
			sb.WriteString(fmt.Sprintf("%v", val))
		}
	}
	return sb.String()
}
