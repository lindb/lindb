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

package metric

import (
	"errors"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/replication"
)

func TestWriteAPI_Sum(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replication.NewMockChannelManager(ctrl)
	api := NewWriteAPI(cm)
	// param error
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 500,
	})

	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(errors.New("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum?db=dal&cluster=dal&c=1",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 500,
	})

	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/sum?db=dal&cluster=dal&c=1",
		HandlerFunc:    api.Sum,
		ExpectHTTPCode: 200,
	})

}
