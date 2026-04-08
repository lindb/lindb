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

package memdb

import "github.com/lindb/lindb/series/field"

// FlushEvent represents flush metadata/index event.
type FlushEvent struct {
	Callback func(err error)
}

// arrowIndexEvent carries data needed to index a new time series from the Arrow path.
type arrowIndexEvent struct {
	nameHash    uint64
	memSeriesID uint32
	namespace   string
	name        string
	attrHash    uint64
	attrs       []struct{ key, value string }
}

// arrowMetaEvent carries data needed to persist field metadata from the Arrow path.
type arrowMetaEvent struct {
	nameHash   uint64
	namespace  []byte
	name       []byte
	fieldMetas []field.Meta
}
