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

package influx

import (
	"github.com/klauspost/compress/gzip"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/tag"

	"bytes"
	"net/http"
	"strings"
	"testing"
)

const _testBody = `
# this is a comment
measurement value=12
measurement value=12 1439587925
measurement,foo=bar value=12 
measurement,foo=bar value=12 1439587925
measurement,foo=bar,bat=baz value=12,otherval=21 1439587925
total\ disk\ free,volumes=/net\,/home\,/ value=442221834240i 1435362189575692182
`

func makeGzipData(body []byte) []byte {
	var w bytes.Buffer
	gw := gzip.NewWriter(&w)
	_, _ = gw.Write(body)
	_ = gw.Close()
	return w.Bytes()
}

func Test_Parse(t *testing.T) {
	r := bytes.NewReader(makeGzipData([]byte(_testBody)))

	req, err := http.NewRequest(http.MethodPut, "", r)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	req.Header.Set("Content-Encoding", "gzip")

	enrichedTags := []tag.Tag{
		tag.NewTag([]byte("ip"), []byte("1.1.1.1")),
		tag.NewTag([]byte("region"), []byte("sh")),
	}
	metrics, err := Parse(req, enrichedTags, "ns")
	assert.Nil(t, err)
	assert.NotNil(t, metrics)
	assert.Len(t, metrics.Metrics, 6)
}

func Test_getGzipError(t *testing.T) {
	req, err := http.NewRequest(http.MethodPut, "", strings.NewReader(_testBody))
	assert.Nil(t, err)
	assert.NotNil(t, req)
	req.Header.Set("Content-Encoding", "gzip")
	_, err = Parse(req, nil, "ns")
	assert.NotNil(t, err)
}

func Test_parseError(t *testing.T) {
	const _badBody = `
# bad data
measurement value=12,bat=baz vvv=baz
measurement value=12 1439587925
`
	r := bytes.NewReader(makeGzipData([]byte(_badBody)))

	req, err := http.NewRequest(http.MethodPut, "", r)
	assert.Nil(t, err)
	assert.NotNil(t, req)
	req.Header.Set("Content-Encoding", "gzip")

	_, err = Parse(req, nil, "ns")
	assert.NotNil(t, err)
}

func Test_getPrecisionMultiplier(t *testing.T) {
	assert.Equal(t, int64(-1000000), getPrecisionMultiplier("ns"))
	assert.Equal(t, int64(-1000), getPrecisionMultiplier("us"))
	assert.Equal(t, int64(1), getPrecisionMultiplier("ms"))
	assert.Equal(t, int64(1000), getPrecisionMultiplier("s"))
	assert.Equal(t, int64(60000), getPrecisionMultiplier("m"))
	assert.Equal(t, int64(3600000), getPrecisionMultiplier("h"))
}
