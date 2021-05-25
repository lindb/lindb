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

package query

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetadataAPI_Handle_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		parseSQLFunc = parseSQL

		ctrl.Finish()
	}()

	api := NewMetadataAPI(nil, nil, nil, nil, nil)

	// case 1: database name not input
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 500,
		RequestBody:    []string{},
	})
	// case 2: parse sql err
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?db=db&sql=show d",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 500,
		RequestBody:    []string{},
	})
	// case 3: wrong type
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?db=db&sql=select f1 from cpu",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 500,
		RequestBody:    []string{},
	})
	// case 4: unknown metadata type
	parseSQLFunc = func(ql string) (*stmt.Metadata, error) {
		return &stmt.Metadata{}, nil
	}
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?db=db&sql=select f1 from cpu",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 500,
		RequestBody:    []string{},
	})
}

func TestMetadataAPI_ShowDatabases(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	databaseService := service.NewMockDatabaseService(ctrl)
	api := NewMetadataAPI(databaseService, nil, nil, nil, nil)

	databaseService.EXPECT().List().Return(nil, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?sql=show databases",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 200,
		RequestBody:    []string{},
	})

	databaseService.EXPECT().List().Return(
		[]*models.Database{
			{Name: "test1"},
			{Name: "test2"},
		},
		nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?sql=show databases",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 200,
		RequestBody: models.Metadata{
			Type:   stmt.Database.String(),
			Values: []string{"test1", "test2"},
		},
	})

	databaseService.EXPECT().List().Return(nil, fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?sql=show databases",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 500,
	})
}

func TestMetadataAPI_SuggestCommon(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := parallel.NewMockExecutorFactory(ctrl)
	exec := parallel.NewMockMetadataExecutor(ctrl)

	factory.EXPECT().NewMetadataBrokerExecutor(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any()).Return(exec).AnyTimes()

	api := NewMetadataAPI(nil, nil, nil, factory, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?sql=show namespaces",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 500,
	})

	exec.EXPECT().Execute().Return(nil, fmt.Errorf("err"))
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?db=db1&sql=show namespaces",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 500,
	})
	exec.EXPECT().Execute().Return([]string{"a", "b"}, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?db=db1&sql=show namespaces",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 200,
		RequestBody: models.Metadata{
			Type:   stmt.Namespace.String(),
			Values: []string{"a", "b"},
		},
	})

	exec.EXPECT().Execute().Return([]string{"ddd"}, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?db=db1&sql=show fields from cpu",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 500,
	})
	exec.EXPECT().Execute().Return([]string{string(encoding.JSONMarshal(&[]field.Meta{{Name: "test", Type: field.SumField}}))}, nil)
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/query/metadata?db=db1&sql=show fields from cpu",
		HandlerFunc:    api.Handle,
		ExpectHTTPCode: 200,
	})
}
