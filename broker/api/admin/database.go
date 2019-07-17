package admin

import (
	"net/http"

	"github.com/eleme/lindb/broker/api"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/service"
)

// DatabaseAPI represents database admin rest api
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
func (d *DatabaseAPI) GetByName(w http.ResponseWriter, r *http.Request) {
	databaseName, err := api.GetParamsFromRequest("name", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	database, err := d.databaseService.Get(databaseName)
	if err != nil {
		//TODO add not found error?????
		api.Error(w, err)
		return
	}
	api.OK(w, database)
}

// Save creates the database config if there is no database
// config with the name database.Name, otherwise update the config
func (d *DatabaseAPI) Save(w http.ResponseWriter, r *http.Request) {
	database := models.Database{}
	err := api.GetJSONBodyFromRequest(r, &database)
	if err != nil {
		api.Error(w, err)
		return
	}
	err = d.databaseService.Save(database)
	if err != nil {
		api.Error(w, err)
		return
	}
	api.NoContent(w)
}
