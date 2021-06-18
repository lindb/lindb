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

package parallel

import (
	"context"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/stmt"
)

//go:generate mockgen -source=./result_merger.go -destination=./result_merger_mock.go -package=parallel

var mergeLogger = logger.GetLogger("parallel", "merger")

// for testing
var newGroupingAgg = aggregation.NewGroupingAggregator

// ResultMerger represents a merger which merges the task response and aggregates the result
type ResultMerger interface {
	// merge merges the task response and aggregates the result
	merge(resp *pb.TaskResponse)

	// close closes merger
	close()
}

// resultMerger implements ResultMerger interface
type resultMerger struct {
	resultSet chan *series.TimeSeriesEvent
	query     *stmt.Query

	groupAgg aggregation.GroupingAggregator

	events chan *pb.TaskResponse

	closed chan struct{}
	ctx    context.Context

	stats *models.QueryStats
	err   error
}

// newResultMerger create a result merger
func newResultMerger(ctx context.Context, query *stmt.Query, resultSet chan *series.TimeSeriesEvent) ResultMerger {
	merger := &resultMerger{
		resultSet: resultSet,
		query:     query,
		events:    make(chan *pb.TaskResponse),
		closed:    make(chan struct{}),
		ctx:       ctx,
	}
	go func() {
		defer close(merger.closed)
		merger.process()
	}()
	return merger
}

// merge merges and aggregates the result
func (m *resultMerger) merge(resp *pb.TaskResponse) {
	m.events <- resp
}

// close closes merger
func (m *resultMerger) close() {
	close(m.events)
	// waiting process completed
	<-m.closed
	// send result set
	if m.err != nil {
		m.resultSet <- &series.TimeSeriesEvent{Err: m.err, Stats: m.stats}
	} else {
		if m.groupAgg == nil {
			// no data do merge logic
			return
		}
		// send all series data
		resultSet := m.groupAgg.ResultSet()
		if len(resultSet) > 0 {
			m.resultSet <- &series.TimeSeriesEvent{
				SeriesList: resultSet,
				Stats:      m.stats,
			}
		}
	}
}

// process consumes response event, then handles response
func (m *resultMerger) process() {
	for {
		select {
		case event, ok := <-m.events:
			if !ok {
				return
			}
			// if handle event fail, return
			if !m.handleEvent(event) {
				return
			}
		case <-m.ctx.Done():
			return
		}
	}
}

// handleEvent merges the task response
func (m *resultMerger) handleEvent(resp *pb.TaskResponse) bool {
	// handle query stats
	m.handleQueryStats(resp)

	data := resp.Payload
	tsList := &pb.TimeSeriesList{}
	err := tsList.Unmarshal(data)
	if err != nil {
		m.err = err
		return false
	}

	if m.groupAgg == nil {
		AggregatorSpecs := make(aggregation.AggregatorSpecs, len(tsList.FieldAggSpecs))
		for idx, aggSpec := range tsList.FieldAggSpecs {
			AggregatorSpecs[idx] = aggregation.NewAggregatorSpec(
				field.Name(aggSpec.FieldName),
				field.Type(aggSpec.FieldType),
			)
			for _, funcType := range aggSpec.FuncTypeList {
				AggregatorSpecs[idx].AddFunctionType(function.FuncType(funcType))
			}
		}
		// interval ratio is 1 when do merge result.
		m.groupAgg = newGroupingAgg(m.query.Interval, 1, m.query.TimeRange, AggregatorSpecs)
	}

	for _, ts := range tsList.TimeSeriesList {
		// if no field data, ignore this response
		if len(ts.Fields) == 0 {
			return true
		}
		fields := make(map[field.Name][]byte)
		for k, v := range ts.Fields {
			fields[field.Name(k)] = v
		}
		m.groupAgg.Aggregate(series.NewGroupedIterator(ts.Tags, fields))
	}
	return true
}

// handleQueryStats handles query stats if need
func (m *resultMerger) handleQueryStats(resp *pb.TaskResponse) {
	if len(resp.Stats) > 0 {
		// if has query stats, need merge task query stats
		if m.stats == nil {
			m.stats = models.NewQueryStats()
		}
		storageStats := models.NewStorageStats()
		_ = encoding.JSONUnmarshal(resp.Stats, storageStats)
		storageStats.NetCost = timeutil.NowNano() - resp.SendTime
		storageStats.NetPayload = len(resp.Stats) + len(resp.Payload)
		m.stats.MergeStorageTaskStats(resp.TaskID, storageStats)
	}
}

// suggestResultMerger represents the merger which merges the distribution suggest query task's result set
type suggestResultMerger struct {
	resultSet chan []string
}

// newSuggestResultMerger creates the suggest result merger
func newSuggestResultMerger(resultSet chan []string) ResultMerger {
	return &suggestResultMerger{
		resultSet: resultSet,
	}
}

// merge merges the suggest results
func (m *suggestResultMerger) merge(resp *pb.TaskResponse) {
	result := &models.SuggestResult{}
	err := encoding.JSONUnmarshal(resp.Payload, result)
	if err != nil {
		mergeLogger.Error("unmarshal suggest result set", logger.Error(err))
		return
	}
	m.resultSet <- result.Values
}

// close closes the suggest result merge
func (m *suggestResultMerger) close() {
	close(m.resultSet)
}
