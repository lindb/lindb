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
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
)

var IndexKVMerger kv.MergerType = "IndexKVMergerV1"

func init() {
	// register index kv merger
	kv.RegisterMerger(IndexKVMerger, NewIndexKVMerger)
}

type indexKVMerger struct {
	flusher  kv.Flusher
	kvWriter table.StreamWriter
}

func NewIndexKVMerger(kvFlusher kv.Flusher) (kv.Merger, error) {
	kvWriter, err := kvFlusher.StreamWriter()
	if err != nil {
		return nil, err
	}
	return &indexKVMerger{
		flusher:  kvFlusher,
		kvWriter: kvWriter,
	}, nil
}

func (m *indexKVMerger) Init(_ map[string]interface{}) {}

func (m *indexKVMerger) Merge(bucketID uint32, buckets [][]byte) error {
	// TODO: reuse bucket
	trieBucket := model.NewTrieBucket()
	for _, bucket := range buckets {
		err := trieBucket.Unmarshal(bucket)
		if err != nil {
			return err
		}
	}
	m.kvWriter.Prepare(bucketID)
	// write index bucket
	if err := trieBucket.Write(m.kvWriter); err != nil {
		return err
	}
	// commit index bucket
	return m.kvWriter.Commit()
}
