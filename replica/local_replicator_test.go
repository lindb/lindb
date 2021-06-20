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
	"fmt"
	"testing"

	"github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/tsdb"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestLocalReplicator_Replica(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	shard := tsdb.NewMockShard(ctrl)
	replicator := NewLocalReplicator(shard)
	assert.True(t, replicator.IsReady())
	replicator.Replica(1, []byte{1, 2, 3})

	metricList := field.MetricList{
		Metrics: []*field.Metric{{Name: "test"}},
	}
	data, _ := metricList.Marshal()
	shard.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("errj"))
	replicator.Replica(1, data)
}
