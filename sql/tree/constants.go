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

package tree

const (
	field   = "field"
	unknown = "unknown"
)

// StateType represents state statement type.
type StateType uint8

const (
	// RootAlive represents show root alive(node)  statement.
	RootAlive StateType = iota + 1
	// BrokerAlive represents show broker alive(node)  statement.
	BrokerAlive
	// StorageAlive represents show storage alive(node) statement.
	StorageAlive
)
