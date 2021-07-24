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

package linmetric

import (
	"testing"

	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/tag"

	"github.com/stretchr/testify/assert"
)

func Test_Gather(t *testing.T) {
	gather := NewGather(
		WithGlobalKeyValueOption(tag.KeyValuesFromMap(map[string]string{
			"host": "alpha",
			"ip":   "1.1.1.1",
		})),
		WithReadRuntimeOption(),
		WithNamespaceOption("default-ns"),
	)
	_ = gather.Gather()
}

func Test_Gather_appendKeyValuesToFront(t *testing.T) {
	gather1 := &gather{
		keyValues: tag.KeyValues{
			{Key: "1", Value: "a"},
			{Key: "1", Value: "b"},
		}}
	var m = &protoMetricsV1.Metric{
		Tags: []*protoMetricsV1.KeyValue{
			{Key: "1", Value: "0"},
			{Key: "b", Value: "2"},
			{Key: "c", Value: "2"},
		},
	}
	gather1.appendKeyValuesToFront(m)
	assert.Equal(t, []*protoMetricsV1.KeyValue{
		{Key: "1", Value: "a"},
		{Key: "1", Value: "b"},
		{Key: "1", Value: "0"},
		{Key: "b", Value: "2"},
		{Key: "c", Value: "2"},
	}, m.Tags)
}
