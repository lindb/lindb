package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/pathutil"
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

	clusters := database.Clusters
	for i := range clusters {
		cluster := clusters[i]

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
	data, _ := json.Marshal(database)
	return db.repo.Put(context.TODO(), pathutil.GetDatabaseConfigPath(database.Name), data)
}

// Get returns the database config in the state's repo, if not exist return ErrNotExist
func (db *databaseService) Get(name string) (*models.Database, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("database name must not be null")
	}
	configBytes, err := db.repo.Get(context.TODO(), pathutil.GetDatabaseConfigPath(name))
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
	data, err := db.repo.List(context.TODO(), constants.DatabaseConfigPath)
	if err != nil {
		return result, err
	}
	for _, val := range data {
		db := &models.Database{}
		err = json.Unmarshal(val, db)
		if err != nil {
			logger.GetLogger("service/db").
				Warn("unmarshal data error",
					logger.String("data", string(val)))
		} else {
			result = append(result, db)
		}
	}
	return result, nil
}
