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

package stmt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/encoding"
)

func TestMetadataType_String(t *testing.T) {
	assert.Equal(t, "database", Database.String())
	assert.Equal(t, "namespace", Namespace.String())
	assert.Equal(t, "metric", Metric.String())
	assert.Equal(t, "field", Field.String())
	assert.Equal(t, "tagKey", TagKey.String())
	assert.Equal(t, "tagValue", TagValue.String())
	assert.Equal(t, "unknown", MetadataType(0).String())
}

func TestMetadata_MarshalJSON(t *testing.T) {
	query := Metadata{
		Namespace:  "ns",
		MetricName: "test",
		Type:       TagValue,
		Condition: &BinaryExpr{
			Left: &ParenExpr{Expr: &BinaryExpr{
				Left:     &InExpr{Key: "ip", Values: []string{"1.1.1.1", "2.2.2.2"}},
				Operator: AND,
				Right:    &EqualsExpr{Key: "region", Value: "sh"},
			}},
			Operator: AND,
			Right: &ParenExpr{Expr: &BinaryExpr{
				Left:     &EqualsExpr{Key: "path", Value: "/data"},
				Operator: OR,
				Right:    &EqualsExpr{Key: "path", Value: "/home"},
			}},
		},
		TagKey: "tagKey",
		Prefix: "prefix",
		Limit:  100,
	}

	data := encoding.JSONMarshal(&query)
	query1 := Metadata{}
	err := encoding.JSONUnmarshal(data, &query1)
	assert.NoError(t, err)
	assert.Equal(t, query, query1)
}

func TestMetadata_Marshal_Fail(t *testing.T) {
	query := &Metadata{}
	err := query.UnmarshalJSON([]byte{1, 2, 3})
	assert.NotNil(t, err)
	err = query.UnmarshalJSON([]byte("{\"condition\":\"123\"}"))
	assert.NotNil(t, err)
}
