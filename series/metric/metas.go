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

import "encoding/binary"

// ID represents metric id.
type ID uint32

// EmptyMetricID represents empty value for metric id.
const EmptyMetricID = ID(0)

func (i ID) MarshalBinary() []byte {
	var scratch [4]byte
	binary.LittleEndian.PutUint32(scratch[:], uint32(i))
	return scratch[:]
}
