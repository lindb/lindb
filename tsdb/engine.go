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

package tsdb

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
)

//go:generate mockgen -source=./engine.go -destination=./engine_mock.go -package=tsdb

// for testing
var (
	mkDirIfNotExist = fileutil.MkDirIfNotExist
	listDir         = fileutil.ListDir
	decodeToml      = ltoml.DecodeToml
	newDatabaseFunc = newDatabase
)

var engineLogger = logger.GetLogger("tsdb", "Engine")

// Engine represents a time series engine
type Engine interface {
	// CreateDatabase creates database instance by database's name
	// return success when creating database's path successfully
	CreateDatabase(databaseName string) (Database, error)
	// GetDatabase returns the time series database by given name
	GetDatabase(databaseName string) (Database, bool)
	// FlushDatabase produces a signal to workers for flushing memory database by name
	FlushDatabase(ctx context.Context, databaseName string) bool
	// Close closes the cached time series databases
	Close()

	//FIXME stone1100
	// databaseMetaFlusher flushes database meta periodically
	//databaseMetaFlusher(ctx context.Context)
}

// engine implements Engine
type engine struct {
	cfg              config.TSDB        // the common cfg of time series database
	databases        sync.Map           // databaseName -> Database
	ctx              context.Context    // context
	cancel           context.CancelFunc // cancel function of flusher
	dataFlushChecker DataFlushChecker
}

// NewEngine creates an engine for manipulating the databases
func NewEngine(cfg config.TSDB) (Engine, error) {
	engine, err := newEngine(cfg)
	if err != nil {
		return nil, err
	}
	return engine, nil
}

// newEngine creates an engine
func newEngine(cfg config.TSDB) (*engine, error) {
	// create time series storage path
	if err := mkDirIfNotExist(cfg.Dir); err != nil {
		return nil, fmt.Errorf("create time sereis storage path[%s] erorr: %s", cfg.Dir, err)
	}
	e := &engine{
		cfg: cfg,
	}
	e.ctx, e.cancel = context.WithCancel(context.Background())
	e.dataFlushChecker = newDataFlushChecker(e.ctx)
	e.dataFlushChecker.Start()

	if err := e.load(); err != nil {
		engineLogger.Error("load engine data error when create a new engine", logger.Error(err))
		// close opened engine
		e.Close()
		return nil, err
	}
	return e, nil
}

// CreateDatabase creates database instance by database's name
// return success when creating database's path successfully
func (e *engine) CreateDatabase(databaseName string) (Database, error) {
	dbPath := filepath.Join(e.cfg.Dir, databaseName)
	if err := mkDirIfNotExist(dbPath); err != nil {
		return nil, fmt.Errorf("create database[%s]'s path with error: %s", databaseName, err)
	}
	cfgPath := optionsPath(dbPath)
	cfg := &databaseConfig{}
	if fileutil.Exist(cfgPath) {
		if err := decodeToml(cfgPath, cfg); err != nil {
			return nil, fmt.Errorf("load database[%s] config from file[%s] with error: %s",
				databaseName, cfgPath, err)
		}
	}
	db, err := newDatabaseFunc(databaseName, dbPath, cfg, e.dataFlushChecker)
	if err != nil {
		return nil, err
	}
	e.databases.Store(databaseName, db)
	return db, nil
}

// GetDatabase returns the time series database by given name
func (e *engine) GetDatabase(databaseName string) (Database, bool) {
	item, _ := e.databases.Load(databaseName)
	db, ok := item.(Database)
	return db, ok
}

// Close closes the cached time series databases
func (e *engine) Close() {
	if e.dataFlushChecker != nil {
		e.dataFlushChecker.Stop()
	}

	e.databases.Range(func(key, value interface{}) bool {
		db := value.(Database)
		if err := db.Close(); err != nil {
			engineLogger.Error("close database", logger.Error(err))
		}
		return true
	})
}

// FlushDatabase produces a signal to workers for flushing memory database by name
func (e *engine) FlushDatabase(ctx context.Context, name string) bool {
	item, ok := e.databases.Load(name)
	if !ok {
		return false
	}
	db := item.(Database)
	if err := db.Flush(); err != nil {
		//TODO add log and metric
		return false
	}
	return true
}

// load loads the time series engines if exist
func (e *engine) load() error {
	databaseNames, err := listDir(e.cfg.Dir)
	if err != nil {
		return err
	}
	for _, databaseName := range databaseNames {
		_, err := e.CreateDatabase(databaseName)
		if err != nil {
			return err
		}
	}
	return nil
}

//FIXME stone1100
//func (e *engine) databaseMetaFlusher(ctx context.Context) {
//	ticker := time.NewTicker(flushMetaInterval.Load())
//	defer ticker.Stop()
//
//	select {
//	case <-ctx.Done():
//		return
//	case <-ticker.C:
//	}
//
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case <-ticker.C:
//			//e.flushAllDatabases(ctx)
//		}
//	}
//}
