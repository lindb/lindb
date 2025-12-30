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

package transfer

import (
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/spi/types"
)

var transferRegistry = make(map[option.EngineType]Transfer)

// RegisterTransfer registers transfer for given format.
func RegisterTransfer(engine option.EngineType, transfer Transfer) {
	transferRegistry[engine] = transfer
}

func GetTransfer(engine option.EngineType) Transfer {
	return transferRegistry[engine]
}

type Transfer interface {
	Schema() *types.TableSchema
	ToPage(data []byte) (*types.Page, error)
}
