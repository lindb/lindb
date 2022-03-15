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

package sql

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

// metricMetadataStmtParser represents metadata statement parser
type metricMetadataStmtParser struct {
	baseStmtParser
	metadataType stmt.MetricMetadataType
	tagKey       string
	prefix       string
}

// newMetricMetadataStmtParser creates a new metric metadata statement parser
func newMetricMetadataStmtParser(metadataType stmt.MetricMetadataType) *metricMetadataStmtParser {
	return &metricMetadataStmtParser{
		metadataType: metadataType,
		baseStmtParser: baseStmtParser{
			exprStack: collections.NewStack(),
			namespace: constants.DefaultNamespace,
		},
	}
}

// build the metric metadata statement.
func (s *metricMetadataStmtParser) build() (stmt.Statement, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.limit <= 0 {
		s.limit = 100
	}
	return &stmt.MetricMetadata{
		Namespace:  s.namespace,
		MetricName: s.metricName,
		Type:       s.metadataType,
		TagKey:     s.tagKey,
		Prefix:     s.prefix,
		Condition:  s.condition,
		Limit:      s.limit,
	}, nil
}

// visitPrefix visits when production prefix expression is entered
func (s *metricMetadataStmtParser) visitPrefix(ctx *grammar.PrefixContext) {
	s.prefix = strutil.GetStringValue(ctx.Ident().GetText())
}

// visitWithTagKey visits when production with tag key expression is entered
func (s *metricMetadataStmtParser) visitWithTagKey(ctx *grammar.WithTagKeyContext) {
	s.tagKey = strutil.GetStringValue(ctx.Ident().GetText())
}
