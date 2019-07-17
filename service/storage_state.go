package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/pathutil"
	"github.com/eleme/lindb/pkg/state"
)

// StorageStateService represents storage cluster state maintain
type StorageStateService interface {
	// Save saves newest storage state for cluster name
	Save(clusterName string, storageState *models.StorageState) error
	// Get gets current storage state for given cluster name, if not exist return ErrNotExist
	Get(clusterName string) (*models.StorageState, error)
}

// storageStateService implements storage state service interface.
// broker need use storage state for write/read operation.
type storageStateService struct {
	repo state.Repository
}

// NewStorageStateService creates storage state service
func NewStorageStateService(repo state.Repository) StorageStateService {
	return &storageStateService{
		repo: repo,
	}
}

// Save saves newest storage state for cluster name
func (s *storageStateService) Save(clusterName string, storageState *models.StorageState) error {
	data, err := json.Marshal(storageState)
	if err != nil {
		return fmt.Errorf("marshal storage state error:%s", err)
	}
	//TODO add timeout????
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = s.repo.Put(ctx, pathutil.GetStorageClusterStatePath(clusterName), data)
	if err != nil {
		return err
	}
	return nil
}

// Get gets current storage state for given cluster name, if not exist return ErrNotExist
func (s *storageStateService) Get(clusterName string) (*models.StorageState, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	data, err := s.repo.Get(ctx, pathutil.GetStorageClusterStatePath(clusterName))
	if err != nil {
		return nil, err
	}
	storageState := &models.StorageState{}
	err = json.Unmarshal(data, storageState)
	if err != nil {
		return nil, err
	}
	return storageState, err
}
