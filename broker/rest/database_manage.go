package rest

import (
	"fmt"
	"net/http"

	"github.com/eleme/lindb/service"

	"github.com/eleme/lindb/pkg/option"
)

// GetDatabase gets a database config by the databaseName.
func GetDatabase(w http.ResponseWriter, r *http.Request) {
	databaseName, err := GetParamsFromRequest("databaseName", r, "", true)
	if err != nil {
		errorResponse(w, err)
		return
	}
	db := service.New()
	database, err := db.Get(databaseName)
	if err != nil {
		errorResponse(w, err)
		return
	}
	okResponse(w, database)
}

// CreateOrUpdateDatabase creates the database config if there is no database
// config with the name database.Name, otherwise update the config
func CreateOrUpdateDatabase(w http.ResponseWriter, r *http.Request) {
	database := &option.Database{}
	err := checkDatabaseParams(r, database)
	if err != nil {
		errorResponse(w, err)
		return
	}
	db := service.New()
	err = db.Create(*database)
	if err != nil {
		errorResponse(w, err)
		return
	}
	noContent(w)
}

//checkDatabaseParams checks whether the database config meets the requirements,
// if false returns en error
func checkDatabaseParams(r *http.Request, database *option.Database) error {
	err := GetJSONBodyFromRequest(r, database)
	if err != nil {
		return err
	}
	if len(database.Name) == 0 {
		return fmt.Errorf("the database name must not be empty")
	}
	if database.NumOfShard <= 0 {
		return fmt.Errorf("the  NumOfShard must lbe greater than 0")
	}
	if database.ReplicaFactor <= 0 {
		return fmt.Errorf("the ReplicaFactor must be greater than 0")
	}
	return nil
}
