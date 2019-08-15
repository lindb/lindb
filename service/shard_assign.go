package service

import (
	"context"
	"encoding/json"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/pathutil"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./shard_assign.go -destination=./shard_assign_mock.go -package service

// ShardAssignService represents database shard assignment maintain
// Master generates assignment, then storing into related storage cluster's state repo.
// Storage node will create the time series engine based on related shard assignment.
type ShardAssignService interface {
	// List returns all database's shard assignments
	List() ([]*models.ShardAssignment, error)
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

// List returns all database's shard assignments
func (s *shardAssignService) List() ([]*models.ShardAssignment, error) {
	data, err := s.repo.List(context.TODO(), constants.DatabaseAssignPath)
	if err != nil {
		return nil, err
	}

	var result []*models.ShardAssignment
	for _, val := range data {
		shardAssign := &models.ShardAssignment{}
		err = encoding.JSONUnmarshal(val, shardAssign)
		if err != nil {
			logger.GetLogger("service/shard/assign").
				Warn("unmarshal data error",
					logger.String("data", string(val)))
		} else {
			result = append(result, shardAssign)
		}
	}
	return result, nil
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
	data, _ := json.Marshal(shardAssign)
	return s.repo.Put(context.TODO(), pathutil.GetDatabaseAssignPath(databaseName), data)
}
