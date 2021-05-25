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

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileName(t *testing.T) {
	assert.Equal(t, "000001.sst", Table(1))
	assert.Equal(t, "1234567891011.sst", Table(1234567891011))

	assert.Equal(t, "MANIFEST-000012", ManifestFileName(12))
	assert.Equal(t, "MANIFEST-123456789", ManifestFileName(123456789))
	assert.Equal(t, "CURRENT", current())
}

func Test_ParseFileName(t *testing.T) {
	assert.Nil(t, ParseFileName("xxx.tt"))
	assert.Nil(t, ParseFileName("aaa.sst"))
	fileDesc := ParseFileName("000001.sst")
	assert.Equal(t, TypeTable, fileDesc.FileType)
	assert.Equal(t, int64(1), fileDesc.FileNumber.Int64())
}
