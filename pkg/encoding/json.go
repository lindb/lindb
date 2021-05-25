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

package encoding

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/logger"
)

var (
	log  = logger.GetLogger("pkg/encoding", "JSONMarshaller")
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

// JSONMarshal returns the JSON encoding of v.
func JSONMarshal(v interface{}) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Error("json marshal error")
	}
	return data
}

// JSONUnmarshal parses the JSON-encoded data and stores the result
func JSONUnmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
