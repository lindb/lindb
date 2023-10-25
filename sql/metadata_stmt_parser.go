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
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

// metadataStmtParser represents metadata statement parser.
type metadataStmtParser struct {
	metadata *stmt.Metadata
}

// newMetadataStmtParser creates a metadata statement parser.
func newMetadataStmtParser(metadataType stmt.MetadataType) *metadataStmtParser {
	return &metadataStmtParser{
		metadata: &stmt.Metadata{MetadataType: metadataType},
	}
}

// visitStorageFilter visits storage filter.
func (m *metadataStmtParser) visitStorageFilter(ctx *grammar.StorageFilterContext) {
	m.metadata.ClusterName = strutil.GetStringValue(ctx.Ident().GetText())
}

// visitBrokerFilter visits broker filter.
func (m *metadataStmtParser) visitBrokerFilter(ctx *grammar.BrokerFilterContext) {
	m.metadata.ClusterName = strutil.GetStringValue(ctx.Ident().GetText())
}

// visitTypeFilter visits the type filter.
func (m *metadataStmtParser) visitTypeFilter(ctx *grammar.TypeFilterContext) {
	m.metadata.Type = strutil.GetStringValue(ctx.Ident().GetText())
}

// visitSource visits source form.
func (m *metadataStmtParser) visitSource(ctx *grammar.SourceContext) {
	switch {
	case ctx.T_STATE_MACHINE() != nil:
		m.metadata.Source = stmt.StateMachineSource
	case ctx.T_STATE_REPO() != nil:
		m.metadata.Source = stmt.StateRepoSource
	}
}

// build the metadata statement.
func (m *metadataStmtParser) build() (stmt.Statement, error) {
	return m.metadata, nil
}
