package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/pathutil"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./storage_cluster.go -destination=./storage_service_mock.go -package service

// StorageClusterService defines storage cluster service interface
type StorageClusterService interface {
	// Save saves storage cluster config
	Save(storageCluster *models.StorageCluster) error
	// Delete deletes storage cluster config
	Delete(name string) error
	// Get storage cluster by given name, if not exist return ErrNotExist
	Get(name string) (*models.StorageCluster, error)
	// List lists all storage cluster config
	List() ([]*models.StorageCluster, error)
}

// storageClusterService implements storage cluster service interface
type storageClusterService struct {
	repo state.Repository
}

// NewStorageClusterService creates storage cluster service
func NewStorageClusterService(repo state.Repository) StorageClusterService {
	return &storageClusterService{repo: repo}
}

// Save saves storage cluster config
func (s *storageClusterService) Save(storageCluster *models.StorageCluster) error {
	if storageCluster.Name == "" {
		return fmt.Errorf("storage cluster name cannot be empty")
	}
	data, _ := json.Marshal(storageCluster)
	//TODO add timeout????
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := s.repo.Put(ctx, pathutil.GetStorageClusterConfigPath(storageCluster.Name), data); err != nil {
		return err
	}
	return nil
}

// Delete deletes storage cluster config
func (s *storageClusterService) Delete(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return s.repo.Delete(ctx, pathutil.GetStorageClusterConfigPath(name))
}

// Get storage cluster by given name
func (s *storageClusterService) Get(name string) (*models.StorageCluster, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	data, err := s.repo.Get(ctx, pathutil.GetStorageClusterConfigPath(name))
	if err != nil {
		return nil, err
	}
	storageCluster := &models.StorageCluster{}
	err = json.Unmarshal(data, storageCluster)
	if err != nil {
		return nil, err
	}
	return storageCluster, err
}

// List lists config of all storage clusters
func (s *storageClusterService) List() ([]*models.StorageCluster, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var result []*models.StorageCluster
	data, err := s.repo.List(ctx, constants.StorageClusterConfigPath)
	if err != nil {
		return result, err
	}
	for _, val := range data {
		storageCluster := &models.StorageCluster{}
		err = json.Unmarshal(val, storageCluster)
		if err != nil {
			logger.GetLogger("service/storage/cluster").
				Warn("unmarshal data error",
					logger.String("data", string(val)))
		} else {
			result = append(result, storageCluster)
		}
	}
	return result, nil
}
