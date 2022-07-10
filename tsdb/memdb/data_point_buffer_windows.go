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

//go:build windows

package memdb

import (
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
)

var (
	unmapFunc = fileutil.Unmap
)

// closeBuffer just closes file and unmap file.
func (d *dataPointBuffer) closeBuffer() {
	for i, buf := range d.buf {
		if err := unmapFunc(d.files[i], buf); err != nil {
			memDBLogger.Error("unmap file in memory database err",
				logger.String("file", d.path), logger.Error(err))
		}
	}
}
