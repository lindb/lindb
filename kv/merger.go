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

package kv

//go:generate mockgen -source ./merger.go -destination=./merger_mock.go -package kv

// MergerType represents the merger type
type MergerType string

// NewMerger represents create merger instance function
type NewMerger func(flusher Flusher) Merger

var mergers = make(map[MergerType]NewMerger)

// RegisterMerger registers family merger
// NOTICE: must register before create family
func RegisterMerger(name MergerType, merger NewMerger) {
	_, ok := mergers[name]
	if ok {
		panic("merger already register")
	}
	mergers[name] = merger
}

// Merger represents merger values of same key when do compaction job(compact/rollup etc.)
type Merger interface {
	// Init initializes merger params or context, before does merge operation
	Init(params map[string]interface{})
	// Merge merges values for same key,
	// merged data will be written into Flusher directly
	// return err if failure
	Merge(key uint32, values [][]byte) error
}
