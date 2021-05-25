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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/models"
)

func TestInitialize(t *testing.T) {
	defer func() {
		newRequest = http.NewRequest
		doRequest = http.DefaultClient.Do
	}()
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = ioutil.ReadAll(r.Body)
		}))
	defer ts.Close()

	newRequest = func(method, url string, body io.Reader) (*http.Request, error) {
		return nil, fmt.Errorf("err")
	}
	init := newInitialize(ts.URL)

	init.initInternalDatabase(models.Database{})
	init.initStorageCluster(config.StorageCluster{})
	newRequest = http.NewRequest

	doRequest = func(req *http.Request) (*http.Response, error) {
		return nil, fmt.Errorf("err")
	}
	init.initInternalDatabase(models.Database{})
	init.initStorageCluster(config.StorageCluster{})

	doRequest = http.DefaultClient.Do
	init.initInternalDatabase(models.Database{})
	init.initStorageCluster(config.StorageCluster{})
}
