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

package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
)

func TestPBModel(t *testing.T) {
	metric := &protoMetricsV1.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		SimpleFields: []*protoMetricsV1.SimpleField{{
			Name:  "f1",
			Type:  protoMetricsV1.SimpleFieldType_DELTA_SUM,
			Value: 1,
		}},
	}

	data, _ := metric.Marshal()
	metric2 := &protoMetricsV1.Metric{}
	_ = metric2.Unmarshal(data)
	assert.Equal(t, *metric, *metric2)
}
