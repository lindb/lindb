package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/pathutil"
	"github.com/eleme/lindb/pkg/state"
)

type ShardAssignService interface {
	Get(databaseName string) (*models.ShardAssignment, error)
	Save(databaseName string, shardAssign *models.ShardAssignment) error
}

type shardAssignService struct {
	repo state.Repository
}

func NewShardAssignService(repo state.Repository) ShardAssignService {
	return &shardAssignService{
		repo: repo,
	}
}

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

func (s *shardAssignService) Save(databaseName string, shardAssign *models.ShardAssignment) error {
	data, err := json.Marshal(shardAssign)
	if err != nil {
		return fmt.Errorf("marshal shard assignment error:%s", err)
	}
	return s.repo.Put(context.TODO(), pathutil.GetDatabaseAssignPath(databaseName), data)
}
