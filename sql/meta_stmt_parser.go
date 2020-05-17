package sql

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/grammar"
	"github.com/lindb/lindb/sql/stmt"
)

// metaStmtParser represents metadata statement parser
type metaStmtParser struct {
	baseStmtParser
	metadataType stmt.MetadataType
	tagKey       string
	prefix       string
}

// newMetaStmtParser creates a new metadata statement parser
func newMetaStmtParser(metadataType stmt.MetadataType) *metaStmtParser {
	return &metaStmtParser{
		metadataType: metadataType,
		baseStmtParser: baseStmtParser{
			exprStack: collections.NewStack(),
			namespace: constants.DefaultNamespace,
		},
	}
}

// build builds the metadata statement
func (s *metaStmtParser) build() (stmt.Statement, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.limit <= 0 {
		s.limit = 100
	}
	return &stmt.Metadata{
		Namespace:  s.namespace,
		MetricName: s.metricName,
		Type:       s.metadataType,
		TagKey:     s.tagKey,
		Prefix:     s.prefix,
		Condition:  s.condition,
		Limit:      s.limit,
	}, nil
}

// visitPrefix visits when production prefix expression is entered
func (s *metaStmtParser) visitPrefix(ctx *grammar.PrefixContext) {
	s.prefix = strutil.GetStringValue(ctx.Ident().GetText())
}

// visitWithTagKey visits when production with tag key expression is entered
func (s *metaStmtParser) visitWithTagKey(ctx *grammar.WithTagKeyContext) {
	s.tagKey = strutil.GetStringValue(ctx.Ident().GetText())
}
