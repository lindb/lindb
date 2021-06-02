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
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series/tag"

	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const enrichTagsQueryKey = "enrich_tag"

// ExtractEnrichTags extracts enriched tags from url query
// query: enriched_tag=host=test&enriched_tag=ip=1.1.1.1&enriched_tag=zone=bj
// extracted_tags: host:test, ip:1.1.1.1, zone=bj
func ExtractEnrichTags(req *http.Request) (tag.Tags, error) {
	q := req.URL.Query()
	return extractTagsFromQuery(q)
}

func extractTagsFromQuery(values url.Values) (tag.Tags, error) {
	var extracted tag.Tags
	for _, section := range values[enrichTagsQueryKey] {
		tagPair := strings.SplitN(section, "=", 2)
		if len(tagPair) != 2 {
			return extracted, fmt.Errorf("%w, query: %s", constants.ErrBadEnrichTagQueryFormat, section)
		}
		if len(tagPair[0]) == 0 || len(tagPair[1]) == 0 {
			continue
		}
		extracted = append(extracted, tag.NewTag(
			[]byte(tagPair[0]),
			[]byte(tagPair[1]),
		))
	}
	return extracted, nil
}
