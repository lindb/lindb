package state

import (
	"context"
	"fmt"
)

// global repository for state storage
var repository Repository

// Repository stores state data, such as metadata/config/status/task etc.
type Repository interface {
	// Get retrieves value for given key from repository
	Get(ctx context.Context, key string) ([]byte, error)
	// Put puts a key-value pair into repository
	Put(ctx context.Context, key string, val []byte) error
	// Delete deletes value for given key from repository
	Delete(ctx context.Context, key string) error
	// Close closes repository and release resources
	Close() error
}

// New creates global state reposistory
func New(repoType string, config interface{}) error {
	if repoType == "etcd" {
		repo, err := newETCDRepository(config)
		if err != nil {
			return err
		}
		repository = repo
		return nil
	}
	return fmt.Errorf("repo type not define, type is:%s", repoType)
}

// GetRepo returns state repository
func GetRepo() Repository {
	return repository
}
