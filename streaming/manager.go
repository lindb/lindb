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

package streaming

import (
	"sync"
)

var (
	mgr          Mananger
	once4Manager sync.Once
)

func GetManager() Mananger {
	if mgr != nil {
		return mgr
	}

	once4Manager.Do(func() {
		mgr = NewManager()
	})
	return mgr
}

type Mananger interface {
	GetDataSource(name string) (DataSource, bool)
	AddDataSource(ds DataSource)
}

type manager struct {
	dataSources map[string]DataSource

	lock sync.RWMutex
}

func NewManager() Mananger {
	return &manager{
		dataSources: make(map[string]DataSource),
	}
}

func (m *manager) GetDataSource(name string) (DataSource, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	ds, ok := m.dataSources[name]
	return ds, ok
}

func (m *manager) AddDataSource(ds DataSource) {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, ok := m.dataSources[ds.Name()]
	if ok {
		return
	}
	// initialize data source
	ds.Initialize()

	m.dataSources[ds.Name()] = ds
}
