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

package metric

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type BrokerRow struct {
	readOnlyRow
}

func (br *BrokerRow) Unmarshal(data []byte) {
	br.m.Init(data, flatbuffers.GetUOffsetT(data))
}

// UnmarshalBrokerRowsInto unmarshal rows into dst, return error if data is invalid
//func UnmarshalBrokerRowsInto(dst []BrokerRow, src []byte) ([]BrokerRow, error) {
//
//}
