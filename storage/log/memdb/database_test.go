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

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestDatabase_Write(t *testing.T) {
	db := NewDatabase()
	err := db.Write("ns", 12, 1234, map[string]string{"a": "b"})
	if err != nil {
		panic(err)
	}
	err = db.Write("ns", 13, 1234, map[string]string{"a": "hi"})
	if err != nil {
		panic(err)
	}
	err = db.Write("test", 14, 1234, map[string]string{"a": "hi"})
	if err != nil {
		panic(err)
	}
	fmt.Println("write")
	indexDB := db.IndexDatabase()
	nsID, _ := indexDB.GetNamespaceID("test")
	logIDs := db.GetLogIDs(nsID)
	fmt.Println(logIDs)
	nsID, _ = indexDB.GetNamespaceID("ns")
	logIDs = db.GetLogIDs(nsID)
	fmt.Println(logIDs)
	keyID, _ := indexDB.GetFieldKeyID(nsID, "a")
	logIDs = db.GetLogIDs(keyID)
	fmt.Println(logIDs)
	valueID, _ := indexDB.GetFieldValueID(keyID, "hi")
	logIDs = db.GetLogIDs(valueID)
	fmt.Println(logIDs)
	name := fmt.Sprintf("%020d.index", 100)
	fmt.Println(name)
	seqNumStr := name[0 : strings.Index(name, "index")-1]
	seq, err := strconv.ParseInt(seqNumStr, 10, 64)
	fmt.Println(seq)
}
