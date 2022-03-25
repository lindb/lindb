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

package flow

// Stage is the definition of query stage
type Stage int

const (
	FilteringStage Stage = iota + 1
	GroupingStage
	ScannerStage
	DownSamplingStage
)

// String returns string value of stage.
func (qs Stage) String() string {
	switch qs {
	case FilteringStage:
		return "filtering"
	case GroupingStage:
		return "grouping"
	case ScannerStage:
		return "scanner"
	case DownSamplingStage:
		return "downSampling"
	default:
		return "unknown"
	}
}
