package aggregation

//go:generate mockgen -source=./container_agg.go -destination=./container_agg_mock.go -package=aggregation

// ContainerAggregator represents the aggregator's container with the aggregates of fields.
type ContainerAggregator interface {
	// GetFieldAggregates returns the aggregates of fields that need query.
	GetFieldAggregates() FieldAggregates
}
