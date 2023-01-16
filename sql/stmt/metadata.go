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

// SourceType represents metadata source.
type SourceType int

const (
	// StateRepoSource represents from state persist repo.
	StateRepoSource = iota + 1
	// StateMachineSource represents from state machine in current memory.
	StateMachineSource
)

// MetadataType represents metadata type.
type MetadataType int

const (
	// MetadataTypes represents all metadata types.
	MetadataTypes MetadataType = iota + 1
	// BrokerMetadata represent broker metadata.
	BrokerMetadata
	// MasterMetadata represent master metadata.
	MasterMetadata
	// StorageMetadata represent storage metadata.
	StorageMetadata
	// RootMetadata represent root metadata.
	RootMetadata
)

// Metadata represent show metadata lin query language.
type Metadata struct {
	MetadataType MetadataType
	Type         string     // broker/master/storage will be used.
	Source       SourceType // source(from state repo or state manager).
	StorageName  string     // storage will be used.
}

// StatementType returns metadata lin query language statement type.
func (m *Metadata) StatementType() StatementType {
	return MetadataStatement
}
