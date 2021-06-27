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

package bootstrap

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
)

// for testing
var (
	newRequest    = http.NewRequest
	defaultClient = http.Client{Timeout: time.Second * 10}
	doRequest     = defaultClient.Do
)

const brokerAPIPrefix = "/api/v1"

// ClusterInitializer initializes cluster(storage/internal database)
type ClusterInitializer struct {
	endpoint string
}

// NewClusterInitializer creates a initializer
func NewClusterInitializer(endpoint string) *ClusterInitializer {
	u, _ := url.Parse(endpoint)
	u.Path = path.Join(u.Path, brokerAPIPrefix)
	return &ClusterInitializer{endpoint: u.String()}
}

// InitStorageCluster initializes the storage cluster
func (i *ClusterInitializer) InitStorageCluster(storageCfg config.StorageCluster) error {
	reader := bytes.NewReader(encoding.JSONMarshal(&storageCfg))
	req, err := newRequest("POST", fmt.Sprintf("%s/storage/cluster", i.endpoint), reader)
	if err != nil {
		return err
	}
	return doPost(req)
}

// InitInternalDatabase initializes internal database
func (i *ClusterInitializer) InitInternalDatabase(database models.Database) error {
	reader := bytes.NewReader(encoding.JSONMarshal(&database))
	req, err := newRequest("POST", fmt.Sprintf("%s/database", i.endpoint), reader)
	if err != nil {
		return err
	}
	return doPost(req)
}

// doPost does http post request
func doPost(req *http.Request) error {
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	writeResp, err := doRequest(req)
	if err != nil {
		return err
	}
	return writeResp.Body.Close()
}
