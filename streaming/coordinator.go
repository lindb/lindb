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

package streaming

import (
	"reflect"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/streaming/cep"
)

type Coordinator struct {
	streamings map[string]*models.Streaming
	databases  map[string]*models.Database
	logger     logger.Logger
}

func NewCoordinator() meta.Watcher {
	return &Coordinator{
		streamings: make(map[string]*models.Streaming),
		databases:  make(map[string]*models.Database),
		logger:     logger.GetLogger("Streaming", "Coordinator"),
	}
}

// OnEvent implements [meta.Watcher].
func (c *Coordinator) OnEvent(e meta.Event) {
	c.logger.Info("receive event", logger.Any("type", reflect.TypeOf(e)), logger.Any("event", e))
	switch event := e.(type) {
	case *models.Database:
		c.databases[event.Name] = event

		// trigger scheduling streaming
		c.scheduleStreaming()
	case *models.Streaming:
		if event.Observer != meta.CurrentObserver() {
			c.logger.Warn("streaming job observer not match, skip processing",
				logger.String("streaming", event.Name),
				logger.String("expect", meta.CurrentObserver()),
				logger.String("actual", event.Observer))
			return
		}
		c.streamings[event.Name] = event
		// trigger scheduling streaming
		c.scheduleStreaming()
	case *models.DeleteDatabase:
		delete(c.databases, event.Database)

		// trigger scheduling streaming
		c.scheduleStreaming()
	case *models.DeleteStreaming:
		delete(c.streamings, event.Streaming)

		// trigger scheduling streaming
		c.scheduleStreaming()
	case *models.ModifyStreamingJob:
		engine, ok := c.getCEPEngine(event.Streaming)
		if !ok {
			c.logger.Warn("streaming job deploy failed, engine not found",
				logger.String("streaming", event.Streaming), logger.String("job", event.JobName))
			return
		}
		err := engine.DeployJob(event.JobName, event.Script)
		if err != nil {
			c.logger.Error("undeploy streaming job failed",
				logger.String("job", event.JobName), logger.Error(err))
		}
	case *models.DeleteStreamingJob:
		engine, ok := c.getCEPEngine(event.Streaming)
		if !ok {
			c.logger.Warn("streaming job undeploy failed, engine not found",
				logger.String("streaming", event.Streaming), logger.String("job", event.JobName))
			return
		}
		err := engine.UndeployJob(event.JobName)
		if err != nil {
			c.logger.Error("undeploy streaming job failed",
				logger.String("job", event.JobName), logger.Error(err))
		}
	}
}

func (c *Coordinator) scheduleStreaming() {
	activeDS := make(map[string]struct{})
	for _, streaming := range c.streamings {
		dsName := streaming.Database
		dbCfg, ok := c.databases[dsName]
		if !ok {
			c.logger.Warn("database config not found for streaming scheduling", logger.String("database", dsName),
				logger.String("streaming", streaming.Name))
			continue
		}
		activeDS[dsName] = struct{}{}
		ds, ok := GetManager().GetDataSource(dsName)
		if !ok {
			ds = NewDataSource(dbCfg)
			GetManager().AddDataSource(ds)
		}
		if err := ds.ScheduleStream(streaming); err != nil {
			c.logger.Error("schedule streaming failed", logger.String("streaming", streaming.Name), logger.Error(err))
		}
	}

	// shutdown inactive data sources
	datasources := GetManager().GetDataSources()
	for _, ds := range datasources {
		_, ok := activeDS[ds.Name()]
		if ok {
			continue
		}
		c.logger.Info("remove inactive data source", logger.String("ds", ds.Name()), logger.Any("activeDS", activeDS))
		GetManager().RemoteDataSource(ds.Name())
	}
}

func (c *Coordinator) getCEPEngine(streaming string) (*cep.Engine, bool) {
	streamingObj, ok := c.streamings[streaming]
	if !ok {
		c.logger.Warn("streaming not found", logger.String("streaming", streaming))
		return nil, false
	}
	database := streamingObj.Database
	ds, ok := GetManager().GetDataSource(database)
	if !ok {
		return nil, false
	}
	engine, ok := ds.GetEngine(streaming)
	if !ok {
		c.logger.Warn("cep engine not found", logger.String("ds", database), logger.String("streaming", streaming))
		return nil, false
	}
	// only cep engine supported(cep job deploy)
	cepEngine, ok := engine.(*cep.Engine)
	if !ok {
		c.logger.Warn("streaming job deploy failed, engine type invalid", logger.String("ds", database),
			logger.String("streaming", streaming))
		return nil, false
	}
	return cepEngine, true
}

// Subscribe implements [meta.Watcher].
func (c *Coordinator) Subscribe(sub meta.Subscriber) {
}

// Unsubscribe implements [meta.Watcher].
func (c *Coordinator) Unsubscribe(sub meta.Subscriber) {
}

// Close implements [meta.Watcher].
func (c *Coordinator) Close() {
}
