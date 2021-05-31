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

package write

import (
	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/replication"

	"io"
	"net/http"
	"testing"
)

func Test_Influx_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cm := replication.NewMockChannelManager(ctrl)
	api := NewInfluxWriter(cm)

	// missing db param
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/influx",
		HandlerFunc:    api.Write,
		ExpectHTTPCode: 500,
	})

	// enrich_tag bad format
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/influx?db=test&ns=ns2&enrich_tag=a",
		HandlerFunc:    api.Write,
		ExpectHTTPCode: 500,
	})

	// bad influx line format
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/influx?db=test&ns=ns3&enrich_tag=a=b",
		HandlerFunc:    api.Write,
		ExpectHTTPCode: 500,
		RequestBody: `
# bad line
a,v=c,d=f a=2 b=3 c=4
`,
	})

	// write error
	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(io.ErrClosedPipe)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/influx?db=test&ns=ns4&enrich_tag=a=b",
		HandlerFunc:    api.Write,
		ExpectHTTPCode: 500,
		RequestBody: `
# good line
measurement,foo=bar value=12 1439587925
measurement value=12 1439587925
`,
	})

	// no content
	cm.EXPECT().Write(gomock.Any(), gomock.Any()).Return(nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPut,
		URL:            "/metric/influx?db=test&ns=ns4&enrich_tag=a=b",
		HandlerFunc:    api.Write,
		ExpectHTTPCode: 204,
		RequestBody: `
# good line
measurement,foo=bar value=12 1439587925
measurement value=12 1439587925
`,
	})
}
