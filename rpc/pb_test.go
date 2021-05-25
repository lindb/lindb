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
	pb "github.com/lindb/lindb/rpc/proto/field"
)

func TestPBModel(t *testing.T) {
	metric := &pb.Metric{
		Name:      "test",
		Timestamp: timeutil.Now(),
		Fields: []*pb.Field{{
			Name:  "f1",
			Type:  pb.FieldType_Sum,
			Value: 1.0,
		}},
	}

	data, _ := metric.Marshal()
	metric2 := &pb.Metric{}
	_ = metric2.Unmarshal(data)
	assert.Equal(t, *metric, *metric2)
}
