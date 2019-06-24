package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/state"
)

const databaseConfigNode = "/lindb/database/config"

// DatabaseService defines database service interface
type DatabaseService interface {
	// Create creates database based on config
	Create(database option.Database) error
}

// databaseService implements DatabaseService interface
type databaseService struct {
	repo state.Repository
}

// New creates database service
func New() DatabaseService {
	return &databaseService{
		repo: state.GetRepo(),
	}
}

// Create creates database, saving config into state's repo
func (db *databaseService) Create(database option.Database) error {
	data, err := json.Marshal(database)
	if err != nil {
		return fmt.Errorf("marshal database config error:%s", err)
	}
	return db.repo.Put(context.TODO(), databaseConfigNode, data)
}
