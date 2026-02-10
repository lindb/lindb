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

package logs

import (
	"context"
	"sort"

	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	"github.com/lindb/arrow/pkg/logs"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"

	writerpkg "github.com/lindb/lindb/app/broker/write/writer"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

type writer struct {
	ctx          context.Context
	database     writerpkg.DatabaseAccessor[*logs.Log]
	intervalCalc timeutil.IntervalCalculator
}

func NewWriter(
	ctx context.Context,
	databaseCfg models.Database,
	database writerpkg.DatabaseAccessor[*logs.Log],
) *writer {
	// TODO: need validation
	sort.Sort(databaseCfg.Option.Intervals)
	interval := databaseCfg.Option.Intervals[0].Interval
	return &writer{
		ctx:          ctx,
		intervalCalc: interval.Calculator(),
		database:     database,
	}
}

func (c *writer) Write(ctx context.Context, data []byte, encoding constants.EncodingType) error {
	shards := c.database.GetShards()
	numOfShards := uint64(len(shards))
	if numOfShards == 0 {
		return constants.ErrNoAvailableStorageNode
	}

	req := plogotlp.NewExportRequest()
	if err := req.UnmarshalProto(data); err != nil {
		return err
	}

	reqLogs := req.Logs()
	rLogs := reqLogs.ResourceLogs()
	for i := range rLogs.Len() {
		log := rLogs.At(i)
		attr := log.Resource().Attributes()
		scopeLogs := log.ScopeLogs()
		for j := range scopeLogs.Len() {
			sl := scopeLogs.At(j)
			lrs := sl.LogRecords()

			// iterate log records
			for k := range lrs.Len() {
				lr := lrs.At(k)
				log := logs.GetLog()

				log.Timestamp = int64(lr.Timestamp()) // in nanoseconds
				// TODO: check in writable time range

				log.Message = lr.Body().AsString()
				log.Level = lr.SeverityText()
				log.EventName = lr.EventName()

				// set trace id and span id if trace id is not empty
				traceID := lr.TraceID()
				if !traceID.IsEmpty() {
					log.TraceID = append(log.TraceID[:0], traceID[:]...)
					spanID := lr.SpanID()
					log.SpanID = append(log.SpanID[:0], spanID[:]...)
				}

				// merge attributes(scope resource/log record)
				lr.Attributes().Range(func(k string, v pcommon.Value) bool {
					vStr := v.AsString()
					if k != "" && vStr != "" {
						log.Attributes.Append(k, vStr)
					}
					return true
				})
				attr.Range(func(k string, v pcommon.Value) bool {
					vStr := v.AsString()
					if k != "" && vStr != "" {
						log.Attributes.Append(k, vStr)
					}
					return true
				})

				// route log to shard by log hash
				shard := shards[log.Hash()%numOfShards]
				segmentTime := c.intervalCalc.CalcFamilyTime(log.Timestamp / 1000_000)
				segment := shard.GetOrCreateSegment(segmentTime, func() larrow.EntryBuilder[*logs.Log] {
					return logs.NewLogsBuilder(memory.NewGoAllocator())
				})
				// write log to segment
				segment.Write(ctx, log)
			}
		}
	}
	return nil
}
