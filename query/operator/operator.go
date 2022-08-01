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

package operator

//go:generate mockgen -source=./operator.go -destination=./operator_mock.go -package=operator

// Operator represents the query operator.
type Operator interface {
	// Identifier returns identifier value of the operator.
	Identifier() string
	// Execute executes current query operator, return error if failure.
	Execute() error
}

// TrackableOperator represents operator can be tracked.
type TrackableOperator interface {
	// Stats returns the stats of operator.
	Stats() interface{}
}
