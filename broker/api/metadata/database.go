package metadata

import (
	"net/http"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/service"
)

type DatabaseAPI struct {
	databaseService service.DatabaseService
}

// NewDatabaseAPI creates database api instance
func NewDatabaseAPI(databaseService service.DatabaseService) *DatabaseAPI {
	return &DatabaseAPI{
		databaseService: databaseService,
	}
}

// GetByName gets a database config by the name.
func (d *DatabaseAPI) ListDatabaseNames(w http.ResponseWriter, r *http.Request) {
	databases, err := d.databaseService.List()
	if err != nil {
		api.Error(w, err)
		return
	}
	var databaseNames []string
	for _, db := range databases {
		databaseNames = append(databaseNames, db.Name)
	}
	api.OK(w, databaseNames)
}
