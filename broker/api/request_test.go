package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
)

func TestGetJSONBodyFromRequest(t *testing.T) {
	req, _ := http.NewRequest("GET", "/database", nil)
	_, err := GetParamsFromRequest("test", req, "defaultValue", true)
	assert.Error(t, err)

	value, _ := GetParamsFromRequest("test", req, "defaultValue", false)
	assert.Equal(t, "defaultValue", value)

	req2, _ := http.NewRequest("POST", "/database", nil)
	_, err = GetParamsFromRequest("test", req2, "defaultValue", true)
	assert.Error(t, err)

	value2, _ := GetParamsFromRequest("test", req2, "defaultValue", false)
	assert.Equal(t, "defaultValue", value2)

	_, err = GetParamsFromRequest("", req2, "defaultValue", false)
	assert.Error(t, err)

	req, _ = http.NewRequest(http.MethodOptions, "/database", nil)
	_, err = GetParamsFromRequest("test", req, "defaultValue", false)
	assert.Error(t, err)

	req, _ = http.NewRequest(http.MethodGet, "/database?key=value", nil)
	value, err = GetParamsFromRequest("key", req, "defaultValue", true)
	assert.NoError(t, err)
	assert.Equal(t, "value", value)

	req2, _ = http.NewRequest("POST", "/database", nil)
	req2.PostForm = url.Values{"key": []string{"value"}}
	value, err = GetParamsFromRequest("key", req2, "defaultValue", true)
	assert.NoError(t, err)
	assert.Equal(t, "value", value)
}

func TestGetParamsFromRequest(t *testing.T) {
	database := models.Database{Name: "test"}
	databaseByte, _ := json.Marshal(database)
	req, _ := http.NewRequest("POST", "/database", bytes.NewReader(databaseByte))
	newDatabase := &models.Database{}
	_ = GetJSONBodyFromRequest(req, newDatabase)
	assert.Equal(t, newDatabase.Name, database.Name)
}
