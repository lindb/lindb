package execution

import "github.com/lindb/lindb/sql/tree"

type DropDatabaseTask struct {
	statement *tree.DropDatabase
}
