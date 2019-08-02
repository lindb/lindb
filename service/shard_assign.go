package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/pathutil"
	"github.com/lindb/lindb/pkg/state"
)

// ShardAssignService represents database shard assignment maintain
// Master generates assignment, then storing into related storage cluster's state repo.
// Storage node will create tsdb based on related shard assignment.
type ShardAssignService interface {
	// Get gets shard assignment by given database name, if not exist return ErrNotExist
	Get(databaseName string) (*models.ShardAssignment, error)
	// Save saves shard assignment for given database name, if fail return error
	Save(databaseName string, shardAssign *models.ShardAssignment) error
}

// shardAssignService implements shard assign service interface
type shardAssignService struct {
	repo state.Repository
}

// NewShardAssignService creates shard assign service
func NewShardAssignService(repo state.Repository) ShardAssignService {
	return &shardAssignService{
		repo: repo,
	}
}

// Get gets shard assignment by given database name, if not exist return ErrNotExist
func (s *shardAssignService) Get(databaseName string) (*models.ShardAssignment, error) {
	data, err := s.repo.Get(context.TODO(), pathutil.GetDatabaseAssignPath(databaseName))
	if err != nil {
		return nil, err
	}
	shardAssign := &models.ShardAssignment{}
	if err := json.Unmarshal(data, shardAssign); err != nil {
		return nil, err
	}
	return shardAssign, nil
}

// Save saves shard assignment for given database name, if fail return error
func (s *shardAssignService) Save(databaseName string, shardAssign *models.ShardAssignment) error {
	data, err := json.Marshal(shardAssign)
	if err != nil {
		return fmt.Errorf("marshal shard assignment error:%s", err)
	}
	return s.repo.Put(context.TODO(), pathutil.GetDatabaseAssignPath(databaseName), data)
}
