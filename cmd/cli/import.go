package main

import (
	"os"

	"github.com/lindb/lindb/models"
)

// importSQLFile imports and executes SQL statements from the given file path.
func importSQLFile(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	// query and print result in terminal
	executeAndPrint(models.ExecuteParam{SQL: string(data), Database: inputC.db})
	return nil
}
