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

package storagequery

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
)

//go:generate mockgen -source ./tag_search.go -destination=./tag_search_mock.go -package=storagequery

// tagFilterResult represents the tag filter result, include tag key id and tag value ids
type tagFilterResult struct {
	tagKey      uint32
	tagValueIDs *roaring.Bitmap
}

// TagSearch represents the tag filtering by tag filter expr
type TagSearch interface {
	// Filter filters tag value ids base on tag filter expr, if fail return nil, else return tag value ids
	Filter() (map[string]*tagFilterResult, error)
}

// tagSearch implements TagSearch
type tagSearch struct {
	namespace  string
	metricName string
	condition  stmt.Expr
	metadata   metadb.Metadata

	result map[string]*tagFilterResult
	tags   map[string]uint32 // for cache tag key
	err    error
}

// newTagSearch creates tag search
func newTagSearch(namespace, metricName string, condition stmt.Expr, metadata metadb.Metadata) TagSearch {
	return &tagSearch{
		namespace:  namespace,
		metricName: metricName,
		condition:  condition,
		metadata:   metadata,
		tags:       make(map[string]uint32),
		result:     make(map[string]*tagFilterResult),
	}
}

// Filter filters tag value ids base on tag filter expr, if fail return nil, else return tag value ids
func (s *tagSearch) Filter() (map[string]*tagFilterResult, error) {
	s.findTagValueIDsByExpr(s.condition)
	if s.err != nil {
		return nil, s.err
	}
	return s.result, nil
}

// findTagValueIDsByExpr finds tag value ids by expr, recursion filter for expr
func (s *tagSearch) findTagValueIDsByExpr(expr stmt.Expr) {
	if expr == nil {
		return
	}
	if s.err != nil {
		return
	}
	switch expr := expr.(type) {
	case stmt.TagFilter:
		tagKeyID, err := s.getTagKeyID(expr.TagKey())
		if err != nil {
			s.err = err
			return
		}
		tagValueIDs, err := s.metadata.TagMetadata().FindTagValueDsByExpr(tagKeyID, expr)
		if err != nil {
			s.err = err
			return
		}
		if tagValueIDs != nil && !tagValueIDs.IsEmpty() {
			// save atomic tag filter result
			s.result[expr.Rewrite()] = &tagFilterResult{
				tagKey:      tagKeyID,
				tagValueIDs: tagValueIDs,
			}
		}
	case *stmt.ParenExpr:
		s.findTagValueIDsByExpr(expr.Expr)
	case *stmt.NotExpr:
		// find tag value id by expr => (not tag filter) => tag filter
		s.findTagValueIDsByExpr(expr.Expr)
	case *stmt.BinaryExpr:
		if expr.Operator != stmt.AND && expr.Operator != stmt.OR {
			s.err = fmt.Errorf("wrong binary operator in tag filter: %s", stmt.BinaryOPString(expr.Operator))
			return
		}
		s.findTagValueIDsByExpr(expr.Left)
		s.findTagValueIDsByExpr(expr.Right)
	}
}

// getTagKeyID returns the tag key id by tag key
func (s *tagSearch) getTagKeyID(tagKey string) (uint32, error) {
	tagKeyID, ok := s.tags[tagKey]
	if ok {
		return tagKeyID, nil
	}
	tagKeyID, err := s.metadata.MetadataDatabase().GetTagKeyID(s.namespace, s.metricName, tagKey)
	if err != nil {
		return 0, err
	}
	s.tags[tagKey] = tagKeyID
	return tagKeyID, nil
}
