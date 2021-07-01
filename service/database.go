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

package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./database.go -destination=./database_mock.go -package service

// DatabaseService defines database service interface
type DatabaseService interface {
	// Save saves database config
	Save(database *models.Database) error
	// Get gets database config by name, if not exist return ErrNotExist
	Get(name string) (*models.Database, error)
	// List returns all database configs
	List() ([]*models.Database, error)
}

// databaseService implements DatabaseService interface
type databaseService struct {
	ctx  context.Context
	repo state.Repository
}

// NewDatabaseService creates database service
func NewDatabaseService(ctx context.Context, repo state.Repository) DatabaseService {
	return &databaseService{
		ctx:  ctx,
		repo: repo,
	}
}

// Save saves database config into state's repo
func (db *databaseService) Save(database *models.Database) error {
	if len(database.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}
	if len(database.Cluster) == 0 {
		return fmt.Errorf("cluster name cannot eb empty")
	}

	if database.NumOfShard <= 0 {
		return fmt.Errorf("num. of shard must be > 0")
	}
	if database.ReplicaFactor <= 0 {
		return fmt.Errorf("replica factor must be > 0")
	}
	// validate time series engine option
	if err := database.Option.Validate(); err != nil {
		return err
	}
	data, _ := json.Marshal(database)
	return db.repo.Put(db.ctx, constants.GetDatabaseConfigPath(database.Name), data)
}

// Get returns the database config in the state's repo, if not exist return ErrNotExist
func (db *databaseService) Get(name string) (*models.Database, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("database name must not be null")
	}
	configBytes, err := db.repo.Get(db.ctx, constants.GetDatabaseConfigPath(name))
	if err != nil {
		return nil, err
	}
	database := &models.Database{}
	err = json.Unmarshal(configBytes, database)
	if err != nil {
		return nil, err
	}
	return database, nil
}

// List returns all database configs
func (db *databaseService) List() ([]*models.Database, error) {
	var result []*models.Database
	data, err := db.repo.List(db.ctx, constants.DatabaseConfigPath)
	if err != nil {
		return result, err
	}
	for _, val := range data {
		db := &models.Database{}
		err = json.Unmarshal(val.Value, db)
		if err != nil {
			logger.GetLogger("service", "DatabaseService").
				Warn("unmarshal data error",
					logger.String("data", string(val.Value)))
		} else {
			db.Desc = db.String()
			result = append(result, db)
		}
	}
	return result, nil
}
