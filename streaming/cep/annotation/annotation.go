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
	"context"

	"github.com/lindb/common/pkg/logger"

	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/tree"
)

var log = logger.GetLogger("CEP", "Annotation")

type Annotation struct {
	Name  string
	Props *collections.Properties
}

func ParseAnnotation(annotation *tree.Annotation) *Annotation {
	if annotation == nil {
		return nil
	}
	var propSlice []*tree.Property
	for _, elem := range annotation.Elements {
		// TODO: do nested annotation support?
		if prop, ok := elem.(*tree.Property); ok {
			propSlice = append(propSlice, prop)
		}
	}
	props, err := expression.EvalProps(expression.NewEvalContext(context.TODO()), propSlice)
	if err != nil {
		log.Error("failed to eval annotation properties", logger.Any("props", propSlice), logger.Error(err))
		return nil
	}
	return &Annotation{
		Name:  annotation.Name.Value,
		Props: props,
	}
}
