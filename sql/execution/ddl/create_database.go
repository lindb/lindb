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

package ddl

import (
	"context"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/tree"
)

type CreateDatabaseTask struct {
	metaMgr   meta.MetadataManager
	statement *tree.CreateDatabase
}

func NewCreateDatabase(metaMgr meta.MetadataManager, statement *tree.CreateDatabase) Task {
	return &CreateDatabaseTask{
		metaMgr:   metaMgr,
		statement: statement,
	}
}

func (task *CreateDatabaseTask) Name() string {
	return "CREATE DATABASE"
}

func (task *CreateDatabaseTask) Execute(ctx context.Context) error {
	// FIXME: check database exist
	engineType := option.Metric
	for _, option := range task.statement.CreateOptions {
		switch createOption := option.(type) {
		case *tree.EngineOption:
			engineType = createOption.Type
		default:
			panic("unknown option type")
		}
	}
	evalCtx := expression.NewEvalContext(ctx)
	// FIXME: need check alive node/shard/replica/engine type
	database, err := task.buildDatabase(evalCtx, engineType)
	if err != nil {
		return err
	}
	// save database config
	// TODO: remove metadata manager
	if err := task.metaMgr.CreateDatabase(ctx, database); err != nil {
		return err
	}

	return nil
}

func (task *CreateDatabaseTask) buildDatabase(
	evalCtx expression.EvalContext,
	engineType option.EngineType,
) (*models.Database, error) {
	options := option.DatabaseOption{
		Engine: engineType,
	}
	if err := task.evalPropsExpression(evalCtx, task.statement.Props, &options); err != nil {
		return nil, err
	}
	if engineType == option.Metric {
		// rollup interval options
		for _, rollup := range task.statement.Rollup {
			rollupOption := option.Interval{}
			if err := task.evalPropsExpression(evalCtx, rollup.Props, &rollupOption); err != nil {
				return nil, err
			}
			options.Intervals = append(options.Intervals, rollupOption)
		}
	}
	database := &models.Database{
		Name:   task.statement.Name,
		Option: &options,
	}
	database.Default()
	if err := database.Validate(); err != nil {
		return nil, err
	}
	return database, nil
}

func (task *CreateDatabaseTask) evalPropsExpression(evalCtx expression.EvalContext, props []*tree.Property, result any) error {
	values := make(map[string]any)
	for _, prop := range props {
		val, err := expression.Eval(evalCtx, prop.Value)
		if err != nil {
			return err
		}
		values[prop.Name.Value] = val
	}
	data := encoding.JSONMarshal(values)
	if err := encoding.JSONUnmarshal(data, result); err != nil {
		return err
	}
	return nil
}
