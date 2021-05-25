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

// Reference: github.com/influxdata/influxdb/pkg/escape
package escape

import (
	"testing"
)

var s string

func BenchmarkStringEscapeNoEscapes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		s = String("no_escapes")
	}
	_ = s
}

func BenchmarkStringUnescapeNoEscapes(b *testing.B) {
	for n := 0; n < b.N; n++ {
		s = UnescapeString("no_escapes")
	}
}

func BenchmarkManyStringEscape(b *testing.B) {
	tests := []string{
		"this is my special string",
		"a field w=i th == tons of escapes",
		"some,commas,here",
	}

	for n := 0; n < b.N; n++ {
		for _, test := range tests {
			s = String(test)
		}
	}
}

func BenchmarkManyStringUnescape(b *testing.B) {
	tests := []string{
		`this\ is\ my\ special\ string`,
		`a\ field\ w\=i\ th\ \=\=\ tons\ of\ escapes`,
		`some\,commas\,here`,
	}

	for n := 0; n < b.N; n++ {
		for _, test := range tests {
			s = UnescapeString(test)
		}
	}
}

func TestStringEscape(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{
			in:       "",
			expected: "",
		},
		{
			in:       "this is my special string",
			expected: `this\ is\ my\ special\ string`,
		},
		{
			in:       "a field w=i th == tons of escapes",
			expected: `a\ field\ w\=i\ th\ \=\=\ tons\ of\ escapes`,
		},
		{
			in:       "no_escapes",
			expected: "no_escapes",
		},
		{
			in:       "some,commas,here",
			expected: `some\,commas\,here`,
		},
	}

	for _, test := range tests {
		if test.expected != String(test.in) {
			t.Errorf("Got %s, expected %s", String(test.in), test.expected)
		}
	}
}

func TestStringUnescape(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{
			in:       "",
			expected: "",
		},
		{
			in:       `this\ is\ my\ special\ string`,
			expected: "this is my special string",
		},
		{
			in:       `a\ field\ w\=i\ th\ \=\=\ tons\ of\ escapes`,
			expected: "a field w=i th == tons of escapes",
		},
		{
			in:       "no_escapes",
			expected: "no_escapes",
		},
		{
			in:       `some\,commas\,here`,
			expected: "some,commas,here",
		},
	}

	for _, test := range tests {
		if test.expected != UnescapeString(test.in) {
			t.Errorf("Got %s, expected %s", UnescapeString(test.in), test.expected)
		}
	}
}
