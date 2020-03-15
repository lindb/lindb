package query

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/stmt"
)

// storageExecuteContext represents storage query execute context
type storageExecuteContext struct {
	namespace string
	query     *stmt.Query
	shardIDs  []int32

	tagFilterResult map[string]*tagFilterResult

	stats *models.StorageStats // storage query stats track for explain query
}

// newStorageExecuteContext creates storage execute context
func newStorageExecuteContext(namespace string, shardIDs []int32, query *stmt.Query) *storageExecuteContext {
	ctx := &storageExecuteContext{
		namespace: namespace,
		query:     query,
		shardIDs:  shardIDs,
	}
	if query.Explain {
		// if explain query, create storage query stats
		ctx.stats = models.NewStorageStats()
	}
	return ctx
}

// QueryStats returns the storage query stats
func (ctx *storageExecuteContext) QueryStats() *models.StorageStats {
	if ctx.stats != nil {
		ctx.stats.Complete()
	}
	return ctx.stats
}

// setTagFilterResult sets tag filter result
func (ctx *storageExecuteContext) setTagFilterResult(tagFilterResult map[string]*tagFilterResult) {
	ctx.tagFilterResult = tagFilterResult
}
