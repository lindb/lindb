package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"
)

const storageClusterPath = "/storage/cluster"

// StorageClusterService defines storage cluster service interface
type StorageClusterService interface {
	// Save saves storage cluster config
	Save(storageCluster models.StorageCluster) error
	// Delete deletes storage cluster config
	Delete(name string) error
	// Get storage cluster by given name
	Get(name string) (models.StorageCluster, error)
	// List lists all storage cluster config
	List() ([]models.StorageCluster, error)
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
func (s *storageClusterService) Save(storageCluster models.StorageCluster) error {
	if storageCluster.Name == "" {
		return fmt.Errorf("storage cluster name cannot be empty")
	}
	data, err := json.Marshal(storageCluster)
	if err != nil {
		return fmt.Errorf("marshal storage cluster error:%s", err)
	}
	//TODO add timeout????
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err = s.repo.Put(ctx, s.getClusterPath(storageCluster.Name), data)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes storage cluster config
func (s *storageClusterService) Delete(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return s.repo.Delete(ctx, s.getClusterPath(name))
}

// Get storage cluster by given name
func (s *storageClusterService) Get(name string) (models.StorageCluster, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// err = s.repo.Put(ctx, s.getClusterPath(storageCluster.Name), data)
	data, err := s.repo.Get(ctx, s.getClusterPath(name))
	storageCluster := models.StorageCluster{}
	if err != nil {
		return storageCluster, err
	}
	err = json.Unmarshal(data, &storageCluster)
	if err != nil {
		return storageCluster, err
	}
	return storageCluster, err
}

// List lists all storage cluster config
func (s *storageClusterService) List() ([]models.StorageCluster, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var result []models.StorageCluster
	data, err := s.repo.List(ctx, storageClusterPath)
	if err != nil {
		return result, err
	}
	for _, val := range data {
		storageCluster := models.StorageCluster{}
		err = json.Unmarshal(val, &storageCluster)
		if err != nil {
			logger.GetLogger().Warn("unmarshal storage cluster data error", zap.String("data", string(val)))
		} else {
			result = append(result, storageCluster)
		}
	}
	return result, err
}

// getClusterPath return cluster storage path
func (s *storageClusterService) getClusterPath(name string) string {
	return fmt.Sprintf("%s/%s", storageClusterPath, name)
}
