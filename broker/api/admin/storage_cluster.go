package admin

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/eleme/lindb/broker/api"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/service"
)

// StorageClusterAPI represents storage cluster admin rest api
type StorageClusterAPI struct {
	storageClusterService service.StorageClusterService
	log                   *zap.Logger
}

// NewStorageClusterAPI create storage cluster api
func NewStorageClusterAPI(storageClusterService service.StorageClusterService) *StorageClusterAPI {
	return &StorageClusterAPI{
		storageClusterService: storageClusterService,
		log:                   logger.GetLogger(),
	}
}

// Create creates config of storage cluster
func (s *StorageClusterAPI) Create(w http.ResponseWriter, r *http.Request) {
	storage := models.StorageCluster{}
	err := api.GetJSONBodyFromRequest(r, &storage)
	if err != nil {
		api.Error(w, err)
		return
	}
	err = s.storageClusterService.Save(storage)
	if err != nil {
		api.Error(w, err)
		return
	}
	api.NoContent(w)
}

// GetByName gets storage cluster by name
func (s *StorageClusterAPI) GetByName(w http.ResponseWriter, r *http.Request) {
	name, err := api.GetParamsFromRequest("name", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	cluster, err := s.storageClusterService.Get(name)
	if err != nil {
		api.NotFound(w)
		return
	}
	api.OK(w, cluster)
}

// DeleteByName deletes storage cluster by name
func (s *StorageClusterAPI) DeleteByName(w http.ResponseWriter, r *http.Request) {
	name, err := api.GetParamsFromRequest("name", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	err = s.storageClusterService.Delete(name)
	if err != nil {
		api.Error(w, err)
		return
	}
	api.NoContent(w)
}

// List lists all storage clusters
func (s *StorageClusterAPI) List(w http.ResponseWriter, r *http.Request) {
	clusters, err := s.storageClusterService.List()
	if err != nil {
		api.Error(w, err)
		return
	}
	api.OK(w, clusters)
}
