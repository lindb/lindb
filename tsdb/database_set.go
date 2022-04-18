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

package tsdb

import "sync/atomic"

type databaseSet struct {
	value atomic.Value // map[string]Database
}

func newDatabaseSet() *databaseSet {
	var m = make(map[string]Database)
	set := &databaseSet{}
	set.value.Store(m)
	return set
}

func (ds *databaseSet) PutDatabase(newDBName string, newDB Database) {
	oldDBSet := ds.Entries()
	var newDBSet = make(map[string]Database)
	for dbName, db := range oldDBSet {
		newDBSet[dbName] = db
	}
	newDBSet[newDBName] = newDB
	ds.value.Store(newDBSet)
}

func (ds *databaseSet) DropDatabase(newDBName string) {
	oldDBSet := ds.Entries()
	delete(oldDBSet, newDBName)

	var newDBSet = make(map[string]Database)
	for dbName, db := range oldDBSet {
		newDBSet[dbName] = db
	}
	ds.value.Store(newDBSet)
}

func (ds *databaseSet) GetDatabase(dbName string) (Database, bool) {
	db, ok := ds.Entries()[dbName]
	return db, ok
}

func (ds *databaseSet) Entries() map[string]Database {
	return ds.value.Load().(map[string]Database)
}
