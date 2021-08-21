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

package server

//go:generate mockgen -source=./service.go -destination=./service_mock.go -package=server

// Service represents an operational state of server, lifecycle methods to transition between states.
type Service interface {
	// Name returns the service's name
	Name() string
	// Run runs server
	Run() error
	// State returns current service state
	State() State
	// Stop shutdowns server, do some cleanup logic
	Stop()
}
