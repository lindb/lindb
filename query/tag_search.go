package query

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/metadb"
)

//go:generate mockgen -source ./tag_search.go -destination=./tag_search_mock.go -package=query

// filterResult represents the tag filter result, include tag key id and tag value ids
type filterResult struct {
	tagKey      uint32
	tagValueIDs *roaring.Bitmap
}

// TagSearch represents the tag filtering by tag filter expr
type TagSearch interface {
	// Filter filters tag value ids base on tag filter expr, if fail return nil, else return tag value ids
	Filter() (map[string]*filterResult, error)
}

// tagSearch implements TagSearch
type tagSearch struct {
	namespace string
	query     *stmt.Query
	metadata  metadb.Metadata

	result map[string]*filterResult
	tags   map[string]uint32 // for cache tag key
	err    error
}

// newTagSearch creates tag search
func newTagSearch(namespace string, query *stmt.Query, metadata metadb.Metadata) TagSearch {
	return &tagSearch{
		namespace: namespace,
		query:     query,
		metadata:  metadata,
		tags:      make(map[string]uint32),
		result:    make(map[string]*filterResult),
	}
}

// Filter filters tag value ids base on tag filter expr, if fail return nil, else return tag value ids
func (s *tagSearch) Filter() (map[string]*filterResult, error) {
	s.findTagValueIDsByExpr(s.query.Condition)
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
		// save atomic tag filter result
		s.result[expr.Rewrite()] = &filterResult{
			tagKey:      tagKeyID,
			tagValueIDs: tagValueIDs,
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
	tagKeyID, err := s.metadata.MetadataDatabase().GetTagKeyID(s.namespace, s.query.MetricName, tagKey)
	if err != nil {
		return 0, err
	}
	s.tags[tagKey] = tagKeyID
	return tagKeyID, nil
}
