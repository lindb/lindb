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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
)

func TestMetadataStatement(t *testing.T) {
	q, err := Parse("show metadata types")
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Metadata{MetadataType: stmt.MetadataTypes}, q)

	q, err = Parse("show root metadata from state_repo where type='/a/b'")
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Metadata{MetadataType: stmt.RootMetadata, Type: "/a/b", Source: stmt.StateRepoSource}, q)

	q, err = Parse("show broker metadata from state_repo where type='/a/b'")
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Metadata{MetadataType: stmt.BrokerMetadata, Type: "/a/b", Source: stmt.StateRepoSource}, q)

	q, err = Parse("show master metadata from state_machine where type='/a/b'")
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Metadata{MetadataType: stmt.MasterMetadata, Type: "/a/b", Source: stmt.StateMachineSource}, q)

	q, err = Parse("show storage metadata from state_repo where type='/a/b'")
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Metadata{MetadataType: stmt.StorageMetadata,
		Type: "/a/b", Source: stmt.StateRepoSource}, q)

	q, err = Parse("show broker metadata from state_machine where type='/a/b' and broker='test'")
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Metadata{MetadataType: stmt.BrokerMetadata, ClusterName: "test", Type: "/a/b", Source: stmt.StateMachineSource}, q)
}
