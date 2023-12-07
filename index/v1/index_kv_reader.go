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

package v1

import (
	"github.com/lindb/lindb/index/model"
	"github.com/lindb/lindb/kv/version"
)

// IndexKVReader represents index kv reader.
type IndexKVReader interface {
	// GetBucket returns trie bucket by bucket id.
	GetBucket(bucket uint32) (trieBucket *model.TrieBucket, err error)
}

// indexKVReader implements IndexKVReader interface.
type indexKVReader struct {
	snapshot version.Snapshot
}

// NewIndexKVReader creates a IndexKVReader instance.
func NewIndexKVReader(snapshot version.Snapshot) IndexKVReader {
	return &indexKVReader{
		snapshot: snapshot,
	}
}

// GetBucket returns trie bucket by bucket id.
func (r *indexKVReader) GetBucket(bucket uint32) (trieBucket *model.TrieBucket, err error) {
	err = r.snapshot.Load(bucket, func(value []byte) error {
		if trieBucket == nil {
			trieBucket = model.NewTrieBucket()
		}
		if err0 := trieBucket.Unmarshal(value); err0 != nil {
			return err0
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return trieBucket, nil
}
