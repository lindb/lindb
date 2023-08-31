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

package table

import (
	"errors"

	"github.com/lindb/common/pkg/logger"
)

var (
	ErrEmptyKeys = errors.New("empty keys under store builder")
)

const (
	// magic-number in the footer of sst file
	magicNumberOffsetFile uint64 = 0x69632d656d656c65
	// current file layout version
	version0 = 0

	sstFileFooterSize = 4 + // posOfOffset(4)
		4 + // posOfKeys(4)
		1 + // version(1)
		8 // magicNumber(8)
	magicNumberAtFooter = 9
)

var tableLogger = logger.GetLogger("KV", "Table")
