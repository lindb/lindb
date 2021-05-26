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

package common

import (
	"github.com/lindb/lindb/mock"

	"github.com/stretchr/testify/assert"

	"net/http"
	"testing"
)

func _extractQueryHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := ExtractEnrichTags(r)
	if tags != nil {
		_, _ = w.Write([]byte(tags.String()))
	}
	if err == nil {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(500)
}

func Test_ExtractEnrichTags(t *testing.T) {
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodPost,
		URL:            "/write?enrich_tag=a=1",
		RequestBody:    nil,
		HandlerFunc:    _extractQueryHandler,
		ExpectHTTPCode: 200,
	})
}

func Test_extractTagsFromQuery(t *testing.T) {
	tags1, err := extractTagsFromQuery(make(map[string][]string))
	assert.Nil(t, err)
	assert.Empty(t, tags1)

	tags2, err := extractTagsFromQuery(map[string][]string{
		enrichTagsQueryKey: {"a=1", "b=2", "c=3="},
	})
	assert.Nil(t, err)
	assert.Equal(t, ",a=1,b=2,c=3\\=", tags2.String())

	_, err = extractTagsFromQuery(map[string][]string{
		enrichTagsQueryKey: {""},
	})
	assert.NotNil(t, err)

	tags4, err := extractTagsFromQuery(map[string][]string{
		enrichTagsQueryKey: {"=3", "a=1"},
	})
	assert.Nil(t, err)
	assert.Equal(t, ",a=1", tags4.String())
}
