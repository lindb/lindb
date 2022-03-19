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
	"bytes"
	"fmt"
	"testing"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/linmetrics"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/tsdb"

	"github.com/golang/mock/gomock"
	"github.com/klauspost/compress/snappy"
	"github.com/stretchr/testify/assert"
)

func TestLocalReplicator_Replica(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	database := tsdb.NewMockDatabase(ctrl)
	database.EXPECT().Name().Return("test-database").AnyTimes()
	shard := tsdb.NewMockShard(ctrl)
	var interval timeutil.Interval
	_ = interval.ValueOf("10s")
	shard.EXPECT().CurrentInterval().Return(interval).AnyTimes()
	shard.EXPECT().Database().Return(database).AnyTimes()
	shard.EXPECT().ShardID().Return(models.ShardID(1)).AnyTimes()
	family := tsdb.NewMockDataFamily(ctrl)
	family.EXPECT().CommitSequence(gomock.Any(), gomock.Any()).AnyTimes()
	family.EXPECT().AckSequence(gomock.Any(), gomock.Any()).AnyTimes()

	replicator := NewLocalReplicator(&ReplicatorChannel{State: &models.ReplicaState{Leader: 1}}, shard, family)
	assert.True(t, replicator.IsReady())
	// bad sequence
	family.EXPECT().ValidateSequence(gomock.Any(), gomock.Any()).Return(false)
	replicator.Replica(1, []byte{1, 2, 3})

	family.EXPECT().ValidateSequence(gomock.Any(), gomock.Any()).Return(true).AnyTimes()

	// bad compressed data
	replicator.Replica(1, []byte{1, 2, 3})
	// data ok
	buf := &bytes.Buffer{}
	converter := metric.NewProtoConverter()
	var row metric.BrokerRow
	_ = converter.ConvertTo(&protoMetricsV1.Metric{
		Namespace: "test",
		Name:      "test",
		Timestamp: fasttime.UnixMilliseconds(),
		TagsHash:  0,
		SimpleFields: []*protoMetricsV1.SimpleField{
			{Name: "f1", Type: protoMetricsV1.SimpleFieldType_Min, Value: 1},
		},
	}, &row)
	_, _ = row.WriteTo(buf)
	var dst []byte
	dst = snappy.Encode(dst, buf.Bytes())
	shard.EXPECT().LookupRowMetricMeta(gomock.Any()).Return(fmt.Errorf("err"))
	replicator.Replica(1, dst)

	shard.EXPECT().LookupRowMetricMeta(gomock.Any()).Return(nil)
	family.EXPECT().WriteRows(gomock.Any()).Return(fmt.Errorf("err"))
	replicator.Replica(1, dst)
	// bad data
	dst = snappy.Encode(dst, []byte("bad-data"))
	replicator.Replica(1, dst)
}
