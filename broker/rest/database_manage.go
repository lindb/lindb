package rest

import (
	"fmt"
	"net/http"

	"github.com/eleme/lindb/service"

	"github.com/eleme/lindb/pkg/option"
)

// Handler for get database config
func GetDatabase(w http.ResponseWriter, r *http.Request) {
	databaseName, err := GetParamsFromRequest("databaseName", r, "", true)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	db := service.New()
	database, err := db.Get(databaseName)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	OKResponse(w, database)
}

// Handler for create or update database config
func CreateOrUpdateDatabase(w http.ResponseWriter, r *http.Request) {
	database := &option.Database{}
	err := checkDatabaseParams(r, database)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	db := service.New()
	err = db.Create(*database)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
	NoContent(w)
}

func checkDatabaseParams(r *http.Request, database *option.Database) error {
	err := GetJSONBodyFromRequest(r, database)
	if err != nil {
		return err
	}
	if database.Name == "" {
		return fmt.Errorf("the database name must not be null")
	}
	if database.NumOfShard < 0 {
		return fmt.Errorf("the  NumOfShard must not less than zero")
	}
	if database.ReplicaFactor < 0 {
		return fmt.Errorf("the ReplicaFactor must not less than zero")
	}
	return nil
}
