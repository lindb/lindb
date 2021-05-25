// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

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
