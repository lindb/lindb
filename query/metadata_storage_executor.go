package query

import (
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

// metadataStorageExecutor represents the executor which executes metric metadata suggest in storage side
type metadataStorageExecutor struct {
	database tsdb.Database
	request  *stmt.Metadata
	shardIDs []int32
}

// newMetadataStorageExecutor creates a metadata suggest executor in storage side
func newMetadataStorageExecutor(database tsdb.Database, shardIDs []int32,
	request *stmt.Metadata,
) parallel.MetadataExecutor {
	return &metadataStorageExecutor{
		database: database,
		request:  request,
		shardIDs: shardIDs,
	}
}

// Execute executes the metadata suggest query based on query type
func (e *metadataStorageExecutor) Execute() (result []string, err error) {
	req := e.request
	limit := req.Limit

	switch req.Type {
	case stmt.Namespace:
		return e.database.Metadata().MetadataDatabase().SuggestNamespace(req.Prefix, limit)
	case stmt.Metric:
		return e.database.Metadata().MetadataDatabase().SuggestMetrics(req.Namespace, req.Prefix, limit)
	case stmt.TagKey:
		return e.database.Metadata().MetadataDatabase().SuggestTagKeys(req.Namespace, req.MetricName, req.Prefix, limit)
	case stmt.Field:
		fields, err := e.database.Metadata().MetadataDatabase().GetAllFields(req.Namespace, req.MetricName)
		if err != nil {
			return nil, err
		}
		result = append(result, string(encoding.JSONMarshal(fields)))
		return result, nil
	case stmt.TagValue:
		tagKeyID, err := e.database.Metadata().MetadataDatabase().GetTagKeyID(req.Namespace, req.MetricName, req.TagKey)
		if err != nil {
			return nil, err
		}
		result = e.database.Metadata().TagMetadata().SuggestTagValues(tagKeyID, req.TagValue, limit)
	}
	return result, nil
}
