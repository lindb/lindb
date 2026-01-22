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
	"errors"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/tree"
)

type DropJobTask struct {
	metaMgr   meta.MetadataManager
	statement *tree.DropJob
}

func NewDropJob(metaMgr meta.MetadataManager, statement *tree.DropJob) Task {
	return &DropJobTask{
		metaMgr:   metaMgr,
		statement: statement,
	}
}

// Execute implements [Task].
func (c *DropJobTask) Execute(ctx context.Context) error {
	currentParams := ctx.Value(constants.ContextKeyParams)
	var streaming string
	if currentParams != nil {
		if params, ok := currentParams.(*models.ExecuteParam); ok {
			streaming = params.Database
		}
	}
	if streaming == "" {
		return errors.New("streaming name is empty")
	}
	// TODO: check streaming if exists
	if c.statement.Name == "" {
		// TODO: check job name validity
		return errors.New("job name is empty")
	}
	// TODO: addd check if job exists
	return c.metaMgr.DropJob(ctx, streaming, c.statement.Name)
}

// Name implements [Task].
func (c *DropJobTask) Name() string {
	return "DROP JOB"
}
