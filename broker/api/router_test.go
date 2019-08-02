package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/lindb/lindb/models"

	"github.com/stretchr/testify/assert"
)

func TestGetJSONBodyFromRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/database", nil)
	_, err := GetParamsFromRequest("test", req, "defaultValue", true)
	assert.True(t, err != nil)

	value, _ := GetParamsFromRequest("test", req, "defaultValue", false)
	assert.Equal(t, "defaultValue", value)

	req2, _ := http.NewRequest("POST", "/database", nil)
	_, err2 := GetParamsFromRequest("test", req2, "defaultValue", true)
	assert.True(t, err2 != nil)

	value2, _ := GetParamsFromRequest("test", req2, "defaultValue", false)
	assert.Equal(t, "defaultValue", value2)

}

func TestGetParamsFromRequest(t *testing.T) {
	database := models.Database{Name: "test"}
	databaseByte, _ := json.Marshal(database)
	req, _ := http.NewRequest("POST", "/database", bytes.NewReader(databaseByte))
	newDatabase := &models.Database{}
	_ = GetJSONBodyFromRequest(req, newDatabase)
	assert.Equal(t, newDatabase.Name, database.Name)
}
