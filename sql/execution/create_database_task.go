package execution

import (
	"context"

	"github.com/lindb/common/pkg/encoding"
	"github.com/mitchellh/mapstructure"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/sql/tree"
)

type CreateDatabaseTask struct {
	deps      *Deps
	statement *tree.CreateDatabase
}

func NewCreateDatabaseTask(deps *Deps, statement *tree.CreateDatabase) DataDefinitionTask {
	return &CreateDatabaseTask{
		deps:      deps,
		statement: statement,
	}
}

func (task *CreateDatabaseTask) Name() string {
	return "CREATE DATABASE"
}

func (task *CreateDatabaseTask) Execute(ctx context.Context) error {
	options := option.DatabaseOption{}
	err := mapstructure.Decode(task.statement.Props, &options)
	if err != nil {
		return err
	}
	// rollup interval options
	for _, rollup := range task.statement.Rollup {
		rollupOption := option.Interval{}
		err = mapstructure.Decode(rollup, &rollupOption)
		if err != nil {
			return err
		}
		options.Intervals = append(options.Intervals, rollupOption)
	}
	engineType := models.Metric
	for _, option := range task.statement.CreateOptions {
		switch createOption := option.(type) {
		case *tree.EngineOption:
			engineType = createOption.Type
		default:
			panic("unknown option type")
		}
	}
	database := &models.Database{
		Name:   task.statement.Name,
		Engine: engineType,
		Option: &options,
	}
	database.Default()
	err = database.Validate()
	if err != nil {
		return err
	}
	// FIXME: need check alive node/shard/replica
	data := encoding.JSONMarshal(database)
	if err := task.deps.Repo.Put(ctx, constants.GetDatabaseConfigPath(database.Name), data); err != nil {
		return err
	}
	return nil
}
