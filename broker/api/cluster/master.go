package cluser

import (
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/coordinator"
)

// MasterAPI represents query cluster master state
type MasterAPI struct {
	master coordinator.Master
}

// NewMasterAPI creates the master api
func NewMasterAPI(master coordinator.Master) *MasterAPI {
	return &MasterAPI{
		master: master,
	}
}

// GetMaster returns the current cluster's master
func (m *MasterAPI) GetMaster(w http.ResponseWriter, r *http.Request) {
	master := m.master.GetMaster()
	if master == nil {
		api.NotFound(w)
	} else {
		api.OK(w, master)
	}
}
