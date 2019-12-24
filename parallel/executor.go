package parallel

//go:generate mockgen -source=./executor.go -destination=./executor_mock.go -package=parallel

// Executor represents a query executor both storage/broker side.
// When returning query results the following is the order in which processing takes place:
// 1) filtering
// 2) Scanning
// 3) Grouping if need
// 4) Down sampling
// 5) Aggregation
// 6) Functions
// 7) Expressions
type Executor interface {
	// Execute execute query
	// 1) plan query language
	// 2) aggregator data from time series(memory/file/network)
	Execute()
}

// BrokerExecutor represents the broker query executor,
// 1) chooses the storage nodes that the data is relatively complete
// 2) chooses broker nodes for root and intermediate computing from all available broker nodes
// 3) storage node as leaf computing node does filtering and atomic compute
// 4) intermediate computing nodes are optional, only need if has group by query, does order by for grouping
// 4) root computing node does function and expression computing ???? //TODO  need?
// 5) finally returns result set to user  ???? //TODO  need?
//
// NOTICE: there are some scenarios:
// 1) some assignment shards not in query replica shards,
//    maybe some expectant results are lost in data in offline shard, WHY can query not completely data,
//    because of for the system availability.
type BrokerExecutor interface {
	Executor
	// ExecuteContext returns the broker execute context
	ExecuteContext() BrokerExecuteContext
}

// MetadataExecutor represents the metadata query executor, includes:
// 1. suggest metric name
// 2. suggest tag keys by spec metric name
// 3. suggest tag values by spec metric name and tag key
// 4. suggest fields by spec metric name
type MetadataExecutor interface {
	// Execute executes metadata query logic, (both broker and storage need implement it)
	Execute() ([]string, error)
}
