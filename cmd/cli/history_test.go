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

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetHistory() {
	historyList = nil
	historyIndex = 0
}

func Test_escapeUnescapeHistory(t *testing.T) {
	tests := []struct {
		name    string
		raw     string
		escaped string
	}{
		{
			name:    "plain single-line SQL is unchanged",
			raw:     "SELECT * FROM table;",
			escaped: "SELECT * FROM table;",
		},
		{
			name:    "embedded newline is escaped to literal \\n",
			raw:     "SELECT *\nFROM table;",
			escaped: `SELECT *\nFROM table;`,
		},
		{
			name:    "embedded CRLF is escaped",
			raw:     "SELECT *\r\nFROM table;",
			escaped: `SELECT *\r\nFROM table;`,
		},
		{
			name:    "backslash is doubled",
			raw:     `SELECT 'C:\path' FROM t;`,
			escaped: `SELECT 'C:\\path' FROM t;`,
		},
		{
			name:    "literal \\n in SQL string (not a real newline) is escaped correctly",
			raw:     `SELECT 'line1\nline2' FROM t;`,
			escaped: `SELECT 'line1\\nline2' FROM t;`,
		},
		{
			name:    "mixed: newline + backslash",
			raw:     "SELECT *\nFROM t WHERE p = 'C:\\dir';",
			escaped: `SELECT *\nFROM t WHERE p = 'C:\\dir';`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := escapeHistory(tc.raw)
			assert.Equal(t, tc.escaped, got, "escapeHistory mismatch")

			// Round-trip: unescape(escape(s)) == s
			assert.Equal(t, tc.raw, unescapeHistory(got), "round-trip mismatch")
		})
	}
}

func Test_addToHistory(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantList []string
	}{
		{
			name:     "single-line SQL stored as-is",
			input:    "SELECT * FROM table;",
			wantList: []string{"SELECT * FROM table;"},
		},
		{
			name:     "multi-line SQL with embedded newline preserves formatting",
			input:    "SELECT *\nFROM table;",
			wantList: []string{"SELECT *\nFROM table;"},
		},
		{
			name:     "SQL with tabs and newlines preserves formatting",
			input:    "SELECT *\n\tFROM table\n\tWHERE id = 1;",
			wantList: []string{"SELECT *\n\tFROM table\n\tWHERE id = 1;"},
		},
		{
			name:     "history command is ignored",
			input:    "history;",
			wantList: nil,
		},
		{
			name:     "empty string is ignored",
			input:    "   \n\t  ",
			wantList: nil,
		},
		{
			name:     "single entry added to empty history",
			input:    "SELECT 1;",
			wantList: []string{"SELECT 1;"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetHistory()
			addToHistory(tc.input)
			assert.Equal(t, tc.wantList, []string(historyList))
		})
	}
}

func Test_addToHistory_deduplication(t *testing.T) {
	resetHistory()
	addToHistory("SELECT 1;")
	addToHistory("SELECT 2;")
	// Adding a duplicate of the first entry should move it to the end
	addToHistory("SELECT 1;")
	assert.Equal(t, []string{"SELECT 2;", "SELECT 1;"}, []string(historyList))
}

func Test_addToHistory_maxSize(t *testing.T) {
	resetHistory()
	// Generate maxHistorySize+5 unique entries to verify the FIFO eviction kicks in
	for i := range maxHistorySize + 5 {
		addToHistory("SELECT " + string(rune(0x4E00+i)) + ";") // use CJK codepoints for unique chars
	}
	assert.Len(t, historyList, maxHistorySize)
}

func Test_saveAndLoadHistory(t *testing.T) {
	resetHistory()

	tmp, err := os.CreateTemp("", "lin-history-test-*")
	assert.NoError(t, err)
	tmp.Close()
	defer os.Remove(tmp.Name())

	origPath := historyFilePath
	historyFilePath = tmp.Name()
	defer func() { historyFilePath = origPath }()

	// Multi-line SQL with embedded newlines: formatting must survive save→load round-trip
	addToHistory("SELECT * FROM table;")
	addToHistory("SELECT *\nFROM another_table;")
	addToHistory("SELECT *\n\tFROM foo\n\tWHERE id = 1;")

	saveHistory()

	resetHistory()
	loadHistory()

	assert.Equal(t, []string{
		"SELECT * FROM table;",
		"SELECT *\nFROM another_table;",         // newline preserved
		"SELECT *\n\tFROM foo\n\tWHERE id = 1;", // indentation preserved
	}, []string(historyList))
}

// Test_loadHistory_legacyMultilineFile simulates a ~/.lin-history file written by
// an older CLI version where embedded newlines were NOT escaped.
// loadHistory must reassemble the broken lines back into single entries
// (format cannot be recovered, but entries must not be split).
func Test_loadHistory_legacyMultilineFile(t *testing.T) {
	resetHistory()

	tmp, err := os.CreateTemp("", "lin-history-legacy-*")
	assert.NoError(t, err)
	defer os.Remove(tmp.Name())

	legacy := "SELECT * FROM table;\n" +
		"SELECT field1,\n" +
		"field2\n" +
		"FROM another_table;\n" +
		"use mydb;\n"
	_, err = tmp.WriteString(legacy)
	assert.NoError(t, err)
	tmp.Close()

	origPath := historyFilePath
	historyFilePath = tmp.Name()
	defer func() { historyFilePath = origPath }()

	loadHistory()

	assert.Equal(t, []string{
		"SELECT * FROM table;",
		"SELECT field1, field2 FROM another_table;", // joined with spaces (best-effort)
		"use mydb;",
	}, []string(historyList))
}
