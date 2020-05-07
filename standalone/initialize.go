package standalone

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
)

// for testing
var (
	newRequest = http.NewRequest
	doRequest  = http.DefaultClient.Do
)

var initializeLogger = logger.GetLogger("standalone", "Initialize")

// initialize represents initialize standalone cluster(storage/internal database)
type initialize struct {
	endpoint string
}

// newInitialize creates a initialize
func newInitialize(endpoint string) *initialize {
	return &initialize{endpoint: endpoint}
}

// initStorageCluster initializes the storage cluster
func (i *initialize) initStorageCluster(storageCfg config.StorageCluster) {
	reader := bytes.NewReader(encoding.JSONMarshal(&storageCfg))
	req, err := newRequest("POST", fmt.Sprintf("%s/storage/cluster", i.endpoint), reader)
	if err != nil {
		initializeLogger.Error("new create storage cluster request error", logger.Error(err))
		return
	}
	doPost(req)
}

// initInternalDatabase initializes internal database
func (i *initialize) initInternalDatabase(database models.Database) {
	reader := bytes.NewReader(encoding.JSONMarshal(&database))
	req, err := newRequest("POST", fmt.Sprintf("%s/database", i.endpoint), reader)
	if err != nil {
		initializeLogger.Error("new create init request error", logger.Error(err))
		return
	}
	doPost(req)
}

// doPost does http post request
func doPost(req *http.Request) {
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	writeResp, err := doRequest(req)
	if err != nil {
		initializeLogger.Error("do init request error", logger.Error(err))
		return
	}
	_ = writeResp.Body.Close()
}
