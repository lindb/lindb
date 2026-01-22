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

package cep

import (
	"fmt"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/streaming/cep/runtime"
	"github.com/lindb/lindb/streaming/cep/source"
)

type Engine struct {
	db      *models.Database
	source  source.Source
	runtime runtime.Runtime
}

func NewEngine(stream *models.Streaming, db *models.Database) *Engine {
	return &Engine{
		db:      db,
		runtime: runtime.NewRuntime(stream.Name),
	}
}

func (e *Engine) Send(event models.Event) {
	e.source.Receive(event)
}

func (e *Engine) DeployJob(script string) error {
	return e.runtime.Query(script)
}

func (e *Engine) Start() error {
	s, err := source.GetSource(source.SourceType(e.db.Option.Engine), e.runtime)
	if err != nil {
		return fmt.Errorf("source not found for type %s", e.db.Option.Engine)
	}
	e.source = s

	return nil
}

func (e *Engine) Stop() error {
	return nil
}
