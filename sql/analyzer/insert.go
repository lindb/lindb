package analyzer

import "github.com/lindb/lindb/sql/tree"

type Insert struct {
	Table   *tree.Table
	Columns []string
}
