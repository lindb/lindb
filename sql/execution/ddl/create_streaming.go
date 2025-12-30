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

	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/tree"
)

type CreateStreamingTask struct {
	metaMgr   meta.MetadataManager
	statement *tree.CreateStreaming
}

func NewCreateStreaming(metaMgr meta.MetadataManager, statement *tree.CreateStreaming) Task {
	return &CreateStreamingTask{
		metaMgr:   metaMgr,
		statement: statement,
	}
}

func (task *CreateStreamingTask) Name() string {
	return "CREATE STREAMING"
}

func (task *CreateStreamingTask) Execute(ctx context.Context) error {
	if task.statement.Name == "" {
		return errors.New("streaming name is empty")
	}
	if task.statement.Database == "" {
		return errors.New("database name is empty")
	}
	if task.statement.Observer == "" {
		return errors.New("observer name is empty")
	}
	// save streaming config
	// TODO: remove metadata manager
	if err := task.metaMgr.CreateStreaming(ctx, &models.Streaming{
		Name:     task.statement.Name,
		Database: task.statement.Database,
		Observer: task.statement.Observer,
	}); err != nil {
		return err
	}

	return nil
}
