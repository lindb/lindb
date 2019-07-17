package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/pathutil"
	"github.com/eleme/lindb/pkg/state"
)

// DatabaseService defines database service interface
type DatabaseService interface {
	// Save saves database config
	Save(database *models.Database) error
	// Get gets database config by name, if not exist return ErrNotExist
	Get(name string) (*models.Database, error)
}

// databaseService implements DatabaseService interface
type databaseService struct {
	repo state.Repository
}

// NewDatabaseService creates database service
func NewDatabaseService(repo state.Repository) DatabaseService {
	return &databaseService{
		repo: repo,
	}
}

// Save saves database config into state's repo
func (db *databaseService) Save(database *models.Database) error {
	if len(database.Name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}
	if len(database.Clusters) == 0 {
		return fmt.Errorf("cluster is empty")
	}
	for _, cluster := range database.Clusters {
		if len(cluster.Name) == 0 {
			return fmt.Errorf("cluster name is empty")
		}
		if cluster.NumOfShard <= 0 {
			return fmt.Errorf("num. of shard must be > 0")
		}
		if cluster.ReplicaFactor <= 0 {
			return fmt.Errorf("replica factor must be > 0")
		}
	}
	data, err := json.Marshal(database)
	if err != nil {
		return fmt.Errorf("marshal database config error:%s", err)
	}
	return db.repo.Put(context.TODO(), pathutil.GetDatabaseConfigPath(database.Name), data)
}

// Get returns the database config in the state's repo, if not exist return ErrNotExist
func (db *databaseService) Get(name string) (*models.Database, error) {
	if name == "" {
		return nil, fmt.Errorf("database name must not be null")
	}
	configBytes, err := db.repo.Get(context.TODO(), pathutil.GetDatabaseConfigPath(name))
	if err != nil {
		return nil, err
	}
	database := &models.Database{}
	err = json.Unmarshal(configBytes, database)
	if err != nil {
		return database, err
	}
	return database, nil
}
