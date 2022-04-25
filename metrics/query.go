package metrics

import "github.com/lindb/lindb/internal/linmetric"

// BrokerQueryStatistics represents broker query statistics.
type BrokerQueryStatistics struct {
	CreatedTasks         *linmetric.BoundCounter // create query task
	ExpireTasks          *linmetric.BoundCounter // task expire, long-term no response
	AliveTask            *linmetric.BoundGauge   // current executing task(alive)
	EmitResponse         *linmetric.BoundCounter // emit response to parent node
	OmitResponse         *linmetric.BoundCounter // omit response because task evicted
	SentRequest          *linmetric.BoundCounter // send request success
	SentRequestFailures  *linmetric.BoundCounter // send request failure
	SentResponses        *linmetric.BoundCounter // send response to parent success
	SentResponseFailures *linmetric.BoundCounter // send response failure
}

// StorageQueryStatistics represents storage query statistics.
type StorageQueryStatistics struct {
	MetricQuery         *linmetric.BoundCounter // execute metric query success(just plan it)
	MetricQueryFailures *linmetric.BoundCounter // execute metric query failure
	MetaQuery           *linmetric.BoundCounter // metadata query success
	MetaQueryFailures   *linmetric.BoundCounter // metadata query failure
	OmitRequest         *linmetric.BoundCounter // omit request(task no belong to current node, wrong stream etc.)
}

// NewBrokerQueryStatistics creates broker query statistics.
func NewBrokerQueryStatistics() *BrokerQueryStatistics {
	scope := linmetric.BrokerRegistry.NewScope("lindb.broker.query")
	return &BrokerQueryStatistics{
		CreatedTasks:         scope.NewCounter("created_tasks"),
		AliveTask:            scope.NewGauge("alive_tasks"),
		ExpireTasks:          scope.NewCounter("expire_tasks"),
		EmitResponse:         scope.NewCounter("emitted_responses"),
		OmitResponse:         scope.NewCounter("omitted_responses"),
		SentRequest:          scope.NewCounter("sent_requests"),
		SentResponses:        scope.NewCounter("sent_responses"),
		SentResponseFailures: scope.NewCounter("sent_responses_failures"),
		SentRequestFailures:  scope.NewCounter("sent_requests_failures"),
	}
}

// NewStorageQueryStatistics creates a storage query statistics.
func NewStorageQueryStatistics() *StorageQueryStatistics {
	scope := linmetric.StorageRegistry.NewScope("lindb.storage.query")
	return &StorageQueryStatistics{
		MetricQuery:         scope.NewCounter("metric_queries"),
		MetricQueryFailures: scope.NewCounter("metric_query_failures"),
		MetaQuery:           scope.NewCounter("meta_queries"),
		MetaQueryFailures:   scope.NewCounter("meta_query_failures"),
		OmitRequest:         scope.NewCounter("omitted_requests"),
	}
}
