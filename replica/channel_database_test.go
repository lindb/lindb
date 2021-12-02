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

package replica

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/metric"
)

func TestDatabaseChannel_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newDatabaseChannel(context.TODO(),
		models.Database{
			Name:   "database",
			Option: option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}},
		}, 1, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ch)

	converter := metric.NewProtoConverter()
	batch := metric.NewBrokerBatchRows()
	_ = batch.TryAppend(func(row *metric.BrokerRow) error {
		return converter.ConvertTo(&protoMetricsV1.Metric{
			Name:      "cpu",
			Timestamp: timeutil.Now(),
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
			Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
		}, row)
	})
	err = ch.Write(context.TODO(), batch)
	assert.Equal(t, errChannelNotFound, err)

	shardCh := NewMockChannel(ctrl)
	ch1 := ch.(*databaseChannel)
	ch1.insertShardChannel(models.ShardID(0), shardCh)
	familyChannel := NewMockFamilyChannel(ctrl)
	familyChannel.EXPECT().Write(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	shardCh.EXPECT().GetOrCreateFamilyChannel(gomock.Any()).Return(familyChannel).AnyTimes()

	batch = metric.NewBrokerBatchRows()
	_ = batch.TryAppend(func(row *metric.BrokerRow) error {
		return converter.ConvertTo(&protoMetricsV1.Metric{
			Name:      "cpu",
			Timestamp: timeutil.Now(),
			SimpleFields: []*protoMetricsV1.SimpleField{
				{Name: "f1", Type: protoMetricsV1.SimpleFieldType_DELTA_SUM, Value: 1}},
			Tags: []*protoMetricsV1.KeyValue{{Key: "host", Value: "1.1.1.1"}},
		}, row)
	})
	err = ch.Write(context.TODO(), batch)
	assert.Error(t, err)
}

func TestDatabaseChannel_CreateChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ch, err := newDatabaseChannel(context.TODO(),
		models.Database{
			Name:   "database",
			Option: option.DatabaseOption{Intervals: option.Intervals{{Interval: 10 * 1000}}},
		}, 4, nil)
	assert.NoError(t, err)
	assert.NotNil(t, ch)
	shardCh := NewMockChannel(ctrl)
	ch1 := ch.(*databaseChannel)
	ch1.insertShardChannel(models.ShardID(0), shardCh)
	shardCh2, err := ch.CreateChannel(int32(1), models.ShardID(0))
	assert.NoError(t, err)
	assert.Equal(t, shardCh, shardCh2)

	_, err = ch.CreateChannel(0, 1)
	assert.Equal(t, errInvalidShardID, err)
	_, err = ch.CreateChannel(1, 1)
	assert.Equal(t, errInvalidShardID, err)
	_, err = ch.CreateChannel(2, 1)
	assert.Equal(t, errInvalidShardNum, err)

	_, err = ch.CreateChannel(4, 1)
	assert.NoError(t, err)
}
