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

package brokerquery

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
)

func Test_MetricQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)
	stateMgr := broker.NewMockStateManager(ctrl)
	stateMgr.EXPECT().GetCurrentNode().Return(currentNode).AnyTimes()
	taskManager := NewMockTaskManager(ctrl)

	queryFactory := &queryFactory{
		stateMgr:    stateMgr,
		taskManager: taskManager,
	}
	brokerNodes := []models.StatelessNode{
		generateBrokerActiveNode("1.1.1.1", 8000),
		generateBrokerActiveNode("1.1.1.2", 8000),
		currentNode,
		generateBrokerActiveNode("1.1.1.4", 8000),
	}
	stateMgr.EXPECT().GetLiveNodes().Return(brokerNodes).AnyTimes()

	q, err := sql.Parse("select f from cpu")
	assert.NoError(t, err)
	// case 1: database not found
	qry := newMetricQuery(context.Background(),
		"test_db",
		q.(*stmt.Query),
		queryFactory)
	stateMgr.EXPECT().GetDatabaseCfg("test_db").Return(models.Database{}, false)
	_, err = qry.WaitResponse()
	assert.Error(t, err)

	// case 2: storage nodes not exist
	opt := &option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}}
	stateMgr.EXPECT().GetDatabaseCfg("test_db").
		Return(models.Database{Option: opt}, true).
		AnyTimes()
	qry = newMetricQuery(context.Background(),
		"test_db",
		q.(*stmt.Query),
		queryFactory)
	stateMgr.EXPECT().GetQueryableReplicas("test_db").Return(nil, nil)
	_, err = qry.WaitResponse()
	assert.Error(t, err)

	storageNodes := map[string][]models.ShardID{
		"1.1.1.1:9000": {1, 2, 4},
		"1.1.1.2:9000": {3, 6, 9},
		"1.1.1.3:9000": {5, 7, 8},
		"1.1.1.4:9000": {10, 13, 15},
		"1.1.1.5:9000": {11, 12, 14},
	}
	stateMgr.EXPECT().GetQueryableReplicas("test_db").
		Return(storageNodes, nil).AnyTimes()
	stateMgr.EXPECT().GetLiveNodes().
		Return(brokerNodes).AnyTimes()

	// timeout
	eventCh1 := make(chan *series.TimeSeriesEvent)
	taskManager.EXPECT().SubmitMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(eventCh1, nil)
	ctx, cancel := context.WithCancel(context.Background())
	qry = newMetricQuery(ctx,
		"test_db", q.(*stmt.Query),
		queryFactory)
	time.AfterFunc(time.Millisecond*200, cancel)
	_, err = qry.WaitResponse()
	assert.Error(t, err)

	qry = newMetricQuery(context.Background(),
		"test_db", q.(*stmt.Query),
		queryFactory)
	// has error
	eventCh2 := make(chan *series.TimeSeriesEvent)
	taskManager.EXPECT().SubmitMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventCh2, nil)
	time.AfterFunc(time.Millisecond*200, func() {
		eventCh2 <- &series.TimeSeriesEvent{Err: io.ErrClosedPipe}
	})
	_, err = qry.WaitResponse()
	assert.Error(t, err)

	// closed channel
	eventCh3 := make(chan *series.TimeSeriesEvent)
	taskManager.EXPECT().SubmitMetricTask(gomock.Any(), gomock.Any(), gomock.Any()).Return(eventCh3, nil)
	time.AfterFunc(time.Millisecond*200, func() { close(eventCh3) })
	_, err = qry.WaitResponse()
	assert.Error(t, err)
}

// mockSingleIterator returns mock an iterator of single field
func mockSingleIterator(ctrl *gomock.Controller, aggType field.AggType) series.FieldIterator {
	it := series.NewMockFieldIterator(ctrl)
	primitiveIt := series.NewMockPrimitiveIterator(ctrl)
	it.EXPECT().HasNext().Return(true)
	it.EXPECT().Next().Return(primitiveIt)
	primitiveIt.EXPECT().AggType().Return(aggType)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(4, 4.0)
	primitiveIt.EXPECT().HasNext().Return(true)
	primitiveIt.EXPECT().Next().Return(50, 50.0)
	primitiveIt.EXPECT().HasNext().Return(false)
	it.EXPECT().HasNext().Return(false)
	return it
}

func mockTimeSeries(ctrl *gomock.Controller, startTime int64,
	fieldName field.Name, fieldType field.Type,
	aggType field.AggType,
) series.Iterator {
	timeSeries := series.NewMockIterator(ctrl)
	timeSeries.EXPECT().FieldType().Return(fieldType)
	timeSeries.EXPECT().FieldName().Return(fieldName)
	it := mockSingleIterator(ctrl, aggType)
	timeSeries.EXPECT().HasNext().Return(true)
	timeSeries.EXPECT().Next().Return(startTime, it)
	timeSeries.EXPECT().HasNext().Return(false)
	return timeSeries
}

func Test_MetricQuery_makeResultSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var familyTime, _ = timeutil.ParseTimestamp("20190702 19:00:00", "20060102 15:04:05")
	var now, _ = timeutil.ParseTimestamp("20190702 19:10:00", "20060102 15:04:05")

	series1 := mockTimeSeries(ctrl, familyTime, "f1", field.SumField, field.Sum)
	series2 := mockTimeSeries(ctrl, familyTime, "f2", field.MinField, field.Min)
	timeSeries := series.NewMockGroupedIterator(ctrl)

	q, err := sql.Parse("select (f1+f2)*100 as f from cpu group by node")
	assert.NoError(t, err)
	query := q.(*stmt.Query)
	expression := aggregation.NewExpression(timeutil.TimeRange{
		Start: now,
		End:   now + timeutil.OneHour*2,
	}, timeutil.OneMinute, query.SelectItems)
	gomock.InOrder(
		timeSeries.EXPECT().Tags().Return("node2"),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series1),
		timeSeries.EXPECT().HasNext().Return(true),
		timeSeries.EXPECT().Next().Return(series2),
		timeSeries.EXPECT().HasNext().Return(false),
	)
	qry := &metricQuery{
		expression: expression,
		stmtQuery: &stmt.Query{
			MetricName: "1",
			TimeRange:  timeutil.TimeRange{End: 2, Start: 1},
			GroupBy:    []string{"node"},
		},
	}
	_ = qry.makeResultSet(&series.TimeSeriesEvent{
		SeriesList: []series.GroupedIterator{timeSeries},
		Stats: &models.QueryStats{
			TotalCost:   100,
			ExpressCost: 200,
		},
	})
}
