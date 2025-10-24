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

package store

import (
	"maps"
	"sync/atomic"
)

type DatabaseSet struct {
	value atomic.Value // map[string]Database
}

func NewDatabaseSet() *DatabaseSet {
	m := make(map[string]Database)
	set := &DatabaseSet{}
	set.value.Store(m)
	return set
}

func (ds *DatabaseSet) PutDatabase(newDBName string, newDB Database) {
	oldDBSet := ds.Entries()
	newDBSet := make(map[string]Database)
	maps.Copy(newDBSet, oldDBSet)
	newDBSet[newDBName] = newDB
	ds.value.Store(newDBSet)
}

func (ds *DatabaseSet) DropDatabase(newDBName string) {
	oldDBSet := ds.Entries()
	delete(oldDBSet, newDBName)

	newDBSet := make(map[string]Database)
	maps.Copy(newDBSet, oldDBSet)
	ds.value.Store(newDBSet)
}

func (ds *DatabaseSet) GetDatabase(dbName string) (Database, bool) {
	db, ok := ds.Entries()[dbName]
	return db, ok
}

func (ds *DatabaseSet) Entries() map[string]Database {
	return ds.value.Load().(map[string]Database)
}
