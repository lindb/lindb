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

package analyzer

import "github.com/lindb/lindb/sql/tree"

type AnalyzerContext struct {
	Database    string // default database name
	Analysis    *Analysis
	IDAllocator *tree.NodeIDAllocator

	GetFuncReturnType tree.GetFuncReturnType
}

func NewAnalyzerContext(database string, stmt tree.Statement,
	idallocator *tree.NodeIDAllocator, streaming bool,
) *AnalyzerContext {
	ctx := &AnalyzerContext{
		Database:          database,
		Analysis:          NewAnalysis(stmt),
		IDAllocator:       idallocator,
		GetFuncReturnType: tree.GetDefaultFuncReturnType,
	}
	if streaming {
		ctx.GetFuncReturnType = tree.GetStreamingFuncReturnType
	}
	return ctx
}
