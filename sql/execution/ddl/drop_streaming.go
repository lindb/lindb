package ddl

import (
	"context"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/sql/tree"
)

type DropStreamingTask struct {
	metaMgr   meta.MetadataManager
	statement *tree.DropStreaming
}

func NewDropStreaming(metaMgr meta.MetadataManager, statement *tree.DropStreaming) Task {
	return &DropStreamingTask{
		metaMgr:   metaMgr,
		statement: statement,
	}
}

// Execute implements [Task].
func (d *DropStreamingTask) Execute(ctx context.Context) error {
	if d.statement.Exists {
		if _, exist := d.metaMgr.GetStreaming(d.statement.Name); !exist {
			return constants.ErrStreamingNotExist
		}
	}
	return d.metaMgr.DropStreaming(ctx, d.statement.Name)
}

// Name implements [Task].
func (d *DropStreamingTask) Name() string {
	return "DROP STREAMING"
}
