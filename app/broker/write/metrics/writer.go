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

package metrics

import (
	"context"
	"sort"

	"github.com/apache/arrow-go/v18/arrow/memory"
	larrow "github.com/lindb/arrow/pkg/arrow"
	lmetrics "github.com/lindb/arrow/pkg/metrics"
	"github.com/lindb/arrow/pkg/model"

	writerpkg "github.com/lindb/lindb/app/broker/write/writer"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/common/field"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

// writer implements the Writer interface for Arrow-based metric ingestion.
// It routes each metric row to a shard based on AttrHash() so that all
// data points for the same attribute set (label set) always land on the
// same shard, improving series-index locality on the storage side.
type writer struct {
	ctx          context.Context
	database     writerpkg.DatabaseAccessor[*model.Metric]
	intervalCalc timeutil.IntervalCalculator
}

// NewWriter creates a new metric writer backed by the given database accessor.
// RetentionPolicies are sorted ascending so that index 0 holds the finest
// (smallest) interval, which is used to compute the segment family time.
func NewWriter(
	ctx context.Context,
	databaseCfg models.Database,
	database writerpkg.DatabaseAccessor[*model.Metric],
) *writer {
	sort.Sort(databaseCfg.Option.RetentionPolicies)
	return &writer{
		ctx:          ctx,
		intervalCalc: databaseCfg.Option.RetentionPolicies[0].Interval.Calculator(),
		database:     database,
	}
}

// Write decodes an Arrow IPC payload produced by MetricBuilder.Bytes(),
// then routes each metric row to the appropriate shard segment.
// The encoding parameter is accepted for interface compatibility but unused
// because Arrow IPC is self-describing.
func (w *writer) Write(ctx context.Context, data []byte, _ constants.EncodingType) error {
	shards := w.database.GetShards()
	numOfShards := uint64(len(shards))
	if numOfShards == 0 {
		return constants.ErrNoAvailableStorageNode
	}

	// deserialize the Arrow IPC bytes written by MetricBuilder.Bytes()
	reader, err := lmetrics.NewReader(data)
	if err != nil {
		return err
	}
	defer reader.Release()

	numRows := reader.NumRows()
	for row := range numRows {
		// route by AttrHash so the same label-set always goes to the same shard
		attrHash := reader.AttrHash(row)
		shard := shards[attrHash%numOfShards]

		// timestamp is stored in nanoseconds; divide by 1_000_000 to get ms
		timestampMs := reader.Timestamp(row) / 1_000_000
		segmentTime := w.intervalCalc.CalcFamilyTime(timestampMs)

		segment := shard.GetOrCreateSegment(segmentTime, func() larrow.EntryBuilder[*model.Metric] {
			return lmetrics.NewMetricBuilder(memory.NewGoAllocator())
		})

		// reconstruct a model.Metric so the generic segment builder can call Append()
		m := &model.Metric{
			Namespace: reader.Namespace(row),
			Name:      reader.Name(row),
			Timestamp: reader.Timestamp(row),
		}
		reader.Fields(row, func(name string, kind field.Type, value float64) {
			m.Fields = append(m.Fields, model.Field{Name: name, Kind: kind, Value: value})
		})
		m.Attributes = &model.Attributes{}
		reader.Attributes(row, func(key, value string) {
			m.Attributes.Append(key, value)
		})
		reader.Exemplars(row, func(e *model.Exemplar) {
			m.Exemplars = append(m.Exemplars, e)
		})

		segment.Write(ctx, m)
	}
	return nil
}
