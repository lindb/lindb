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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/lindb/common/models"

	"github.com/lindb/lindb/pkg/terminal"
)

type history []string

var (
	// historyList list
	historyList history
	// history index
	historyIndex int //nolint
	// history file path
	historyFilePath string
)

const maxHistorySize = 200

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}
	historyFilePath = filepath.Join(homeDir, ".lin-history")
}

func (h history) String() string {
	writer := models.NewTableFormatter()
	writer.SetStyle(terminal.TableSylte())
	for _, hist := range h {
		writer.AppendRow(table.Row{hist})
	}
	writer.SetColumnConfigs([]table.ColumnConfig{{Number: 1, WidthMax: terminal.GetTerminalWidth() - 3, WidthMaxEnforcer: text.WrapSoft}})
	return writer.Render()
}

// escapeHistory encodes a history entry so it occupies exactly one physical line
// in the history file. Backslash, carriage-return, and newline are escaped:
//
//	'\' → '\\'
//	'\r' → '\r'
//	'\n' → '\n'
func escapeHistory(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	for i := range len(s) {
		switch s[i] {
		case '\\':
			b.WriteString(`\\`)
		case '\r':
			b.WriteString(`\r`)
		case '\n':
			b.WriteString(`\n`)
		default:
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// unescapeHistory is the inverse of escapeHistory. It restores the original
// content (including embedded newlines) from its escaped representation.
func unescapeHistory(s string) string {
	if !strings.ContainsRune(s, '\\') {
		return s // fast path: nothing to unescape
	}
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case '\\':
				b.WriteByte('\\')
				i += 2
				continue
			case 'r':
				b.WriteByte('\r')
				i += 2
				continue
			case 'n':
				b.WriteByte('\n')
				i += 2
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
}

// addToHistory adds a command to the history, ensuring no duplicates and max size.
// Internal whitespace and newlines are preserved so the original SQL formatting
// can be recalled later. Escaping happens in saveHistory.
func addToHistory(cmd string) {
	// Only strip surrounding whitespace; keep internal formatting (newlines, indentation).
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return
	}
	if strings.HasPrefix(strings.ToLower(cmd), "history") {
		return
	}
	// Remove duplicates
	for i, h := range historyList {
		if h == cmd {
			historyList = slices.Delete(historyList, i, i+1)
			break
		}
	}
	// Add to history
	historyList = append(historyList, cmd)
	// Ensure max size
	if len(historyList) > maxHistorySize {
		historyList = historyList[1:]
	}
	historyIndex = len(historyList)
}

// loadHistory loads the history from the history file.
// It handles three file formats transparently:
//  1. New format  – each entry is one escaped line (escape sequences restored by unescapeHistory).
//  2. Legacy broken format – a single SQL was written across multiple lines due to embedded \n;
//     lines are reassembled by joining until a ";" suffix is found.
//  3. Legacy correct format – each entry is already a single line; treated as format 1 (unescape is no-op).
func loadHistory() {
	file, err := os.Open(historyFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Error opening history file:", err)
		}
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var parts []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue // skip blank separator lines
		}
		parts = append(parts, line)
		if strings.HasSuffix(line, ";") {
			// Either a complete new-format entry (single escaped line) or the last
			// fragment of a legacy broken multi-line entry.
			entry := unescapeHistory(strings.Join(parts, " "))
			historyList = append(historyList, entry)
			parts = nil
		}
	}
	// Flush any trailing entry without a ";" (e.g. "use <db>").
	if len(parts) > 0 {
		historyList = append(historyList, unescapeHistory(strings.Join(parts, " ")))
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading history file:", err)
	}
	historyIndex = len(historyList)
}

// saveHistory saves the history to the history file.
// Each entry is escaped so that embedded newlines do not corrupt the line-oriented format.
func saveHistory() {
	file, err := os.Create(historyFilePath)
	if err != nil {
		fmt.Println("Error creating history file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, cmd := range historyList {
		_, err := writer.WriteString(escapeHistory(cmd) + "\n")
		if err != nil {
			fmt.Println("Error writing to history file:", err)
			return
		}
	}
	writer.Flush()
}
