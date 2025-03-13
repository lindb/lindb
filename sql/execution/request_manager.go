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

package execution

import (
	"sync"

	"github.com/lindb/lindb/sql/execution/model"
)

type RequestManager interface {
	CompleteRequet(id model.RequestID, err error)
}

type requestManager struct {
	requests map[model.RequestID]any
	lock     sync.RWMutex
}

func NewRequestManager() RequestManager {
	return &requestManager{
		requests: make(map[model.RequestID]any),
	}
}

func (mgr *requestManager) CompleteRequet(id model.RequestID, err error) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	delete(mgr.requests, id)
}
