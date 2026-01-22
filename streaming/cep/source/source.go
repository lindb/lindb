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

package source

import (
	"fmt"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/streaming/cep/runtime"
)

type SourceType string

type createSourceFunc func(runtime runtime.Runtime) Source

var sourceRegistry = make(map[SourceType]createSourceFunc)

func RegisterSource(sourceType SourceType, source createSourceFunc) {
	sourceRegistry[sourceType] = source
}

func GetSource(sourceType SourceType, runtime runtime.Runtime) (Source, error) {
	fn, ok := sourceRegistry[sourceType]
	if !ok {
		return nil, fmt.Errorf("source type %s not registered", sourceType)
	}
	return fn(runtime), nil
}

type Source interface {
	Receive(e models.Event)
}
