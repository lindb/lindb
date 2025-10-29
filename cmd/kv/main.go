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

package main

import (
	"log"

	"github.com/lindb/roaring"
	"github.com/linxGnu/grocksdb"

	"github.com/lindb/lindb/pkg/encoding"
)

type IntAddMergeOperator struct{}

func (op *IntAddMergeOperator) Name() string {
	return "IntAddMergeOperator"
}

func (op *IntAddMergeOperator) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	value := roaring.New()
	value2 := roaring.New()
	if len(existingValue) > 0 {
		encoding.BitmapUnmarshal(value, existingValue)
	}
	for _, op := range operands {
		encoding.BitmapUnmarshal(value2, op)
		value.Or(value2)
		value2.Clear()
	}
	d, _ := value.ToBytes()
	return d, true
}

func (op *IntAddMergeOperator) PartialMerge(key, leftOperand, rightOperand []byte) ([]byte, bool) {
	value := roaring.New()
	value2 := roaring.New()
	encoding.BitmapUnmarshal(value, leftOperand)
	encoding.BitmapUnmarshal(value2, rightOperand)
	value.Or(value2)
	d, _ := value.ToBytes()
	return d, true
}

func main() {
	opts := grocksdb.NewDefaultOptions()
	opts.SetCreateIfMissing(true)
	opts.SetMergeOperator(&IntAddMergeOperator{})

	db, err := grocksdb.OpenDb(opts, "mydb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	wo := grocksdb.NewDefaultWriteOptions()
	ro := grocksdb.NewDefaultReadOptions()
	v1, _ := roaring.BitmapOf(1, 2, 3).ToBytes()
	v2, _ := roaring.BitmapOf(5, 6).ToBytes()

	db.Merge(wo, []byte("counter"), v1)
	db.Merge(wo, []byte("counter"), v2)

	v, err := db.Get(ro, []byte("counter"))
	if err == nil && v.Exists() {
		value := roaring.New()
		encoding.BitmapUnmarshal(value, v.Data())
		log.Println("counter value =", value.String())
		v.Free()
	} else {
		log.Println("get err:", err)
	}

	// 手动compact
	db.CompactRange(grocksdb.Range{Start: nil, Limit: nil})
	log.Println("compact done.")
}
