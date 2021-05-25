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
	"time"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/series"
)

func TestMetricAPI_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executorFactory := parallel.NewMockExecutorFactory(ctrl)
	brokerExecutor := parallel.NewMockBrokerExecutor(ctrl)
	executeCtx := parallel.NewMockBrokerExecuteContext(ctrl)
	brokerExecutor.EXPECT().ExecuteContext().Return(executeCtx)
	brokerExecutor.EXPECT().Execute()

	executorFactory.EXPECT().NewBrokerExecutor(gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any()).Return(brokerExecutor)

	api := NewMetricAPI(nil, nil, nil, executorFactory, nil)

	ch := make(chan *series.TimeSeriesEvent)

	executeCtx.EXPECT().ResultCh().Return(ch)
	executeCtx.EXPECT().Emit(gomock.Any())
	executeCtx.EXPECT().ResultSet().Return(&models.ResultSet{}, nil)

	time.AfterFunc(100*time.Millisecond, func() {
		ch <- nil
		close(ch)
	})

	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?db=test&sql=select f from cpu",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 200,
	})
}

func TestNewMetricAPI_Search_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	executorFactory := parallel.NewMockExecutorFactory(ctrl)
	api := NewMetricAPI(nil, nil, nil, executorFactory, nil)

	// param error
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 500,
	})

	// param error
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?db=test",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 500,
	})

	brokerExecutor := parallel.NewMockBrokerExecutor(ctrl)
	executeCtx := parallel.NewMockBrokerExecuteContext(ctrl)
	brokerExecutor.EXPECT().ExecuteContext().Return(executeCtx)
	brokerExecutor.EXPECT().Execute()

	executorFactory.EXPECT().NewBrokerExecutor(gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any(), gomock.Any(),
		gomock.Any(), gomock.Any()).Return(brokerExecutor)

	ch := make(chan *series.TimeSeriesEvent)

	executeCtx.EXPECT().ResultCh().Return(ch)
	executeCtx.EXPECT().ResultSet().Return(&models.ResultSet{}, fmt.Errorf("err"))

	time.AfterFunc(100*time.Millisecond, func() {
		close(ch)
	})
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?db=test&sql=select f from cpu",
		HandlerFunc:    api.Search,
		ExpectHTTPCode: 500,
	})
}
