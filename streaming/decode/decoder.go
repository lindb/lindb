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

package decode

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
)

var decoderRegistry = make(map[option.EngineType]Decoder)

// RegisterDecoder registers decoder for given format.
func RegisterDecoder(engine option.EngineType, decoder Decoder) {
	decoderRegistry[engine] = decoder
}

func GetDecoder(engine option.EngineType) Decoder {
	return decoderRegistry[engine]
}

type Decoder interface {
	ToEvent(data []byte) (models.Event, error)
}
