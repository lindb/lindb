package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
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
	historyIndex int
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

// addToHistory adds a command to the history, ensuring no duplicates and max size
func addToHistory(cmd string) {
	if strings.HasPrefix(strings.ToLower(cmd), "history") {
		return
	}
	// Remove duplicates
	for i, h := range historyList {
		if h == cmd {
			historyList = append(historyList[:i], historyList[i+1:]...)
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
	for scanner.Scan() {
		historyList = append(historyList, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading history file:", err)
	}
	historyIndex = len(historyList)
}

// saveHistory saves the history to the history file.
func saveHistory() {
	file, err := os.Create(historyFilePath)
	if err != nil {
		fmt.Println("Error creating history file:", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, cmd := range historyList {
		_, err := writer.WriteString(cmd + "\n")
		if err != nil {
			fmt.Println("Error writing to history file:", err)
			return
		}
	}
	writer.Flush()
}
