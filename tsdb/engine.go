package tsdb

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

//go:generate mockgen -source=./engine.go -destination=./engine_mock.go -package=tsdb

var engineLogger = logger.GetLogger("tsdb", "Engine")

// Engine represents a time series engine
type Engine interface {
	// CreateDatabase creates database instance by database's name
	// return success when creating database's path successfully
	CreateDatabase(databaseName string) (Database, error)
	// GetDatabase returns the time series database by given name
	GetDatabase(databaseName string) (Database, bool)
	// Close closes the cached time series databases
	Close()
}

// engine implements Engine
type engine struct {
	cfg       config.Engine // the common cfg of time series database
	databases sync.Map      // databaseName -> Database
}

// NewEngine creates an engine for manipulating the database
func NewEngine(cfg config.Engine) (Engine, error) {
	// create time series storage path
	if err := fileutil.MkDirIfNotExist(cfg.Dir); err != nil {
		return nil, fmt.Errorf("create time sereis storage path[%s] erorr: %s", cfg.Dir, err)
	}
	f := &engine{cfg: cfg}
	if err := f.load(); err != nil {
		// close opened engine
		f.Close()
		return nil, err
	}
	return f, nil
}

func (f *engine) CreateDatabase(databaseName string) (Database, error) {
	dbPath := filepath.Join(f.cfg.Dir, databaseName)
	if err := fileutil.MkDirIfNotExist(dbPath); err != nil {
		return nil, fmt.Errorf("create database[%s]'s path with error: %s", databaseName, err)
	}
	cfgPath := optionsPath(dbPath)
	cfg := &databaseConfig{}
	if fileutil.Exist(cfgPath) {
		if err := fileutil.DecodeToml(cfgPath, cfg); err != nil {
			return nil, fmt.Errorf("load database[%s] config from file[%s] with error: %s",
				databaseName, cfgPath, err)
		}
	}
	db, err := newDatabase(databaseName, dbPath, cfg)
	if err != nil {
		return nil, err
	}
	f.databases.Store(databaseName, db)
	return db, nil
}

func (f *engine) GetDatabase(databaseName string) (Database, bool) {
	item, _ := f.databases.Load(databaseName)
	db, ok := item.(Database)
	return db, ok
}

func (f *engine) Close() {
	f.databases.Range(func(key, value interface{}) bool {
		db, ok := value.(Database)
		if ok {
			if err := db.Close(); err != nil {
				engineLogger.Error("close database", logger.Error(err))
			}
		}
		return true
	})
}

// load loads the time series engines if exist
func (f *engine) load() error {
	databaseNames, err := fileutil.ListDir(f.cfg.Dir)
	if err != nil {
		return err
	}
	for _, databaseName := range databaseNames {
		_, err := f.CreateDatabase(databaseName)
		if err != nil {
			return err
		}
	}
	return nil
}
