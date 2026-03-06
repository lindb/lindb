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

package annotation

import (
	"strings"

	"github.com/apache/arrow-go/v18/arrow"

	"github.com/lindb/lindb/pkg/collections"
)

type Mapper interface {
	Map(record arrow.RecordBatch) arrow.RecordBatch
}

type createMapperFn func(props *collections.Properties) Mapper

var mappers = map[string]createMapperFn{
	"metric": NewMetricMapper,
}

func CreateMapper(annotation *Annotation) Mapper {
	if annotation == nil {
		return nil
	}
	if createFn, ok := mappers[strings.ToLower(annotation.Name)]; ok {
		return createFn(annotation.Props)
	}
	return nil
}
