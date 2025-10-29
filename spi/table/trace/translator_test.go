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

package trace

import (
	"encoding/hex"
	"fmt"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestTraceID(t *testing.T) {
	var traceID trace.TraceID
	b, err := hex.DecodeString("0fa9e36f12d744ff243a73dde06add94")
	if err != nil {
		panic(err)
	}
	if len(b) != len(traceID) {
		return
	}
	copy(traceID[:], b)
	fmt.Println(traceID.String())
}
