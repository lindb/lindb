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

package brokerquery

import (
	"sync"

	"github.com/google/uuid"

	"github.com/lindb/lindb/models"
)

var (
	rManager            RequestManager
	once4RequestManager sync.Once
)

// RequestManager represents the request manager which store lin query reuqest.
type RequestManager interface {
	// NewRequest creates a new request and returns request id.
	NewRequest(req *models.Request) string
	// CompleteRequest completes a request by given request id.
	CompleteRequest(requestID string)
	// GetAliveRequests returns all alive request.
	GetAliveRequests() []*models.Request
}

// GetRequestManager returns a singleton RequestManager instance.
func GetRequestManager() RequestManager {
	if rManager != nil {
		return rManager
	}
	once4RequestManager.Do(func() {
		rManager = newRequestManager()
	})
	return rManager
}

// requestManager implements RequestManager interface.
type requestManager struct {
	requests map[string]*models.Request

	mutex sync.RWMutex
}

// newRequestManager creates a RequestManager instance.
func newRequestManager() RequestManager {
	return &requestManager{
		requests: make(map[string]*models.Request),
	}
}

// NewRequest creates a new request and returns request id.
func (r *requestManager) NewRequest(req *models.Request) string {
	requestID := uuid.New().String()
	req.RequestID = requestID

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.requests[requestID] = req
	return requestID
}

// CompleteRequest completes a request by given request id.
func (r *requestManager) CompleteRequest(requestID string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.requests, requestID)
}

// GetAliveRequests returns all alive request.
func (r *requestManager) GetAliveRequests() (rs []*models.Request) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, v := range r.requests {
		rs = append(rs, v)
	}
	return
}
