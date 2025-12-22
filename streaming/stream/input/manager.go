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

package input

import (
	"sync"
)

var (
	instance *Manager
	once     sync.Once
)

func GetManager() *Manager {
	once.Do(func() {
		instance = &Manager{
			apps: make(map[string]InputManager),
		}
	})
	return instance
}

type Manager struct {
	apps map[string]InputManager

	mutex sync.Mutex
}

func (mgr *Manager) GetInputHandler(app, stream string) InputHandler {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()

	inputMgr, ok := mgr.apps[app]
	if !ok {
		inputMgr = NewInputManager()
		mgr.apps[app] = inputMgr
	}
	return inputMgr.GetInputHandler(stream)
}
