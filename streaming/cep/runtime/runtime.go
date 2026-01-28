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

package runtime

import (
	"context"
	"fmt"
	"sync"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/streaming/cep/stream"
	"github.com/lindb/lindb/streaming/cep/stream/input"
)

func init() {
	spi.RegisterSourceConnectorProvider(&stream.TableHandle{}, stream.NewSourceConnectorProvider())
}

type Runtime interface {
	AddEventType(eventType any)
	RegisterStreamByType(eventType any) error
	RegisterStreamBySchema(name string, schema *types.TableSchema) error

	DeployJob(name, statement string) error
	UndeployJob(name string) error

	GetInputHandler(stream string) input.InputHandler

	Shutdown()
}

type runtime struct {
	ctx        context.Context
	cacellFunc context.CancelFunc

	database string
	jobs     map[string]JobRuntime

	lock sync.Mutex

	logger logger.Logger
}

func NewRuntime(database string) Runtime {
	ctx, cancel := context.WithCancel(context.Background())
	return &runtime{
		ctx:        ctx,
		cacellFunc: cancel,
		database:   database,
		jobs:       make(map[string]JobRuntime),
		logger:     logger.GetLogger("CEP", "Runtime"),
	}
}

func (r *runtime) RegisterStreamByType(eventType any) error {
	return stream.GetManager().GetStreamManager(r.database).RegisterStreamByType(eventType)
}

func (r *runtime) RegisterStreamBySchema(name string, schema *types.TableSchema) error {
	return stream.GetManager().GetStreamManager(r.database).RegisterStreamBySchema(name, schema)
}

func (r *runtime) GetInputHandler(stream string) input.InputHandler {
	return input.GetManager().GetInputHandler(r.database, stream)
}

func (r *runtime) DeployJob(name, statement string) error {
	jobRuntime, ok := r.getJobRuntime(name)
	if ok && jobRuntime.Statement() == statement {
		// job already deployed
		r.logger.Info("job already deployed", logger.String("job", name))
		return nil
	}

	// create and startup job runtime
	jobRuntime = NewJobRuntime(r.ctx, r.database, name, statement)
	if err := jobRuntime.Startup(); err != nil {
		return err
	}

	// store job runtime
	r.lock.Lock()
	r.jobs[name] = jobRuntime
	r.lock.Unlock()
	r.logger.Info("job deployed", logger.String("job", name))
	return nil
}

func (r *runtime) UndeployJob(name string) error {
	jobRuntime, ok := r.getJobRuntime(name)
	if !ok {
		return fmt.Errorf("job '%s' not found", name)
	}

	// shutdown job runtime
	jobRuntime.Shutdown()

	// remove job runtime
	r.lock.Lock()
	delete(r.jobs, name)
	r.lock.Unlock()
	r.logger.Info("job undeployed", logger.String("job", name))
	return nil
}

// AddEventType implements [Runtime].
func (r *runtime) AddEventType(eventType any) {
	panic("unimplemented")
}

// Shutdown implements [Runtime].
func (r *runtime) Shutdown() {
	r.cacellFunc()
	// remove all inputs for database
	input.GetManager().RemoveInpputByDatabase(r.database)
}

func (r *runtime) getJobRuntime(name string) (JobRuntime, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()

	jobRuntime, ok := r.jobs[name]
	return jobRuntime, ok
}
