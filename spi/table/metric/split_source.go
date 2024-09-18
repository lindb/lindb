package metric

import (
	"errors"
	"fmt"

	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/tsdb"
)

type SplitSourceProvider struct {
	engine tsdb.Engine
}

func NewSplitSourceProvider(engine tsdb.Engine) *SplitSourceProvider {
	return &SplitSourceProvider{
		engine: engine,
	}
}

// getSchema returns table schema based on table handle.
func (msp *SplitSourceProvider) getSchema(db tsdb.Database, table *TableHandle) (metric.ID, *metric.Schema, error) {
	// find metric id(table id)
	metricID, err := db.MetaDB().GetMetricID(table.Namespace, table.Metric)
	if err != nil {
		return 0, nil, err
	}
	// find table schema
	schema, err := db.MetaDB().GetSchema(metricID)
	if err != nil {
		return 0, nil, err
	}
	return metricID, schema, nil
}

func (msp *SplitSourceProvider) buildTableScan(table spi.TableHandle, outputColumns []spi.ColumnMetadata) *TableScan {
	metricTable, ok := table.(*TableHandle)
	if !ok {
		panic(fmt.Sprintf("not support table handle<%T>", table))
	}
	db, ok := msp.engine.GetDatabase(metricTable.Database)
	if !ok {
		panic(fmt.Errorf("%w: %s", constants.ErrDatabaseNotFound, metricTable.Database))
	}

	// find table(metric) schema
	metricID, schema, err := msp.getSchema(db, metricTable)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound) {
			return nil
		}
		// if isn't not found error, throw it
		panic(err)
	}
	// mapping fields for searching
	fields := lo.Filter(schema.Fields, func(field field.Meta, index int) bool {
		return lo.ContainsBy(outputColumns, func(column spi.ColumnMetadata) bool {
			return column.Name == field.Name.String() && column.DataType == types.DTTimeSeries
		})
	})
	// mpaaing tags for grouping
	groupingTags := lo.Filter(schema.TagKeys, func(tagMeta tag.Meta, index int) bool {
		return lo.ContainsBy(outputColumns, func(column spi.ColumnMetadata) bool {
			return column.Name == tagMeta.Key && column.DataType == types.DTString
		})
	})
	fmt.Printf("all fields=%v, group key=%v, select field=%v,output=%v\n", schema.Fields, groupingTags, fields, outputColumns)

	if len(fields)+len(groupingTags) != len(outputColumns) {
		// output columns size not match
		return nil
	}

	var grouping *Grouping
	if len(groupingTags) > 0 {
		grouping = NewGrouping(db, groupingTags)
	}

	return &TableScan{
		db:              db,
		schema:          schema,
		metricID:        metricID,
		timeRange:       metricTable.TimeRange,
		interval:        metricTable.Interval,
		storageInterval: metricTable.StorageInterval,
		fields:          fields,
		grouping:        grouping,
		output:          outputColumns,
	}
}

func (msp *SplitSourceProvider) findPartitions(tableScan *TableScan, partitionIDs []int) (partitions []*Partition) {
	for _, id := range partitionIDs {
		shard, ok := tableScan.db.GetShard(models.ShardID(id))
		if ok {
			// TODO: use storage interval?
			dataFamilies := shard.GetDataFamilies(tableScan.interval.Type(), tableScan.timeRange)
			if len(dataFamilies) > 0 {
				partitions = append(partitions, &Partition{
					shard:    shard,
					families: dataFamilies,
				})
			}
		}
	}
	return
}

// 1. find database/table(metric) schema
// 2. find columns(tags)' values if has predicate
// 3. find partitions
func (msp *SplitSourceProvider) CreateSplitSources(table spi.TableHandle, partitionIDs []int,
	outputColumns []spi.ColumnMetadata, predicate tree.Expression,
) (splits []spi.SplitSource) {
	tableScan := msp.buildTableScan(table, outputColumns)
	if tableScan == nil {
		return
	}
	// find partitions
	partitions := msp.findPartitions(tableScan, partitionIDs)
	if len(partitions) == 0 {
		return
	}

	tableScan.predicate = predicate
	tableScan.lookupColumnValues()
	// TODO: if grouping start tag value collect

	for i := range partitions {
		splits = append(splits, NewSplitSource(tableScan, partitions[i]))
	}

	return
}

type SplitSource struct {
	tableScan *TableScan
	partition *Partition

	// prepare
	seriesIDs       *roaring.Bitmap
	groupingContext flow.GroupingContext
	resultSet       []flow.FilterResultSet
	highKeys        []uint16
	index           int
}

func NewSplitSource(tableScan *TableScan, partition *Partition) *SplitSource {
	return &SplitSource{
		tableScan: tableScan,
		partition: partition,
	}
}

func (mss *SplitSource) lookupSeriesIDs() *roaring.Bitmap {
	var (
		seriesIDs *roaring.Bitmap
		err       error
		ok        bool
	)
	predicate := mss.tableScan.predicate

	if predicate == nil {
		// if predicate nil, find all series ids under metric
		seriesIDs, err = mss.partition.shard.IndexDB().GetSeriesIDsForMetric(mss.tableScan.metricID)
		if err != nil {
			panic(err)
		}
	} else {
		// find series ids based on where condition
		lookup := NewRowLookupVisitor(mss)
		if seriesIDs, ok = predicate.Accept(nil, lookup).(*roaring.Bitmap); !ok {
			panic(constants.ErrSeriesIDNotFound)
		}
	}

	if seriesIDs == nil || seriesIDs.IsEmpty() {
		panic(constants.ErrSeriesIDNotFound)
	}
	return seriesIDs
}

// 1. find series ids
func (mss *SplitSource) Prepare() {
	seriesIDs := mss.lookupSeriesIDs()

	mss.seriesIDs = roaring.New()

	for i := range mss.partition.families {
		family := mss.partition.families[i]
		// check family data if matches condition(series ids)
		resultSet, err := family.Filter(&flow.MetricScanContext{
			MetricID:                mss.tableScan.metricID,
			SeriesIDs:               seriesIDs,
			SeriesIDsAfterFiltering: seriesIDs,
			Fields:                  mss.tableScan.fields,
			TimeRange:               mss.tableScan.timeRange,
			StorageInterval:         mss.tableScan.storageInterval,
		})

		if err != nil && !errors.Is(err, constants.ErrNotFound) {
			panic(err)
		}

		for i := range resultSet {
			rs := resultSet[i]

			// check double, maybe some series ids be filtered out when do grouping.
			finalSeriesIDs := roaring.FastAnd(seriesIDs, rs.SeriesIDs())
			if finalSeriesIDs.IsEmpty() {
				continue
			}

			mss.resultSet = append(mss.resultSet, rs)
			mss.seriesIDs.Or(finalSeriesIDs)
		}
	}

	if mss.seriesIDs.IsEmpty() {
		panic(constants.ErrSeriesIDNotFound)
	}

	if mss.tableScan.hasGrouping() {
		// if it has grouping, do group by tag keys, else just split series ids as batch first.
		seriesIDsAfterGrouping, groupingContext, err := mss.partition.shard.IndexDB().
			GetGroupingContext(mss.tableScan.grouping.tags, mss.seriesIDs)
		if err != nil {
			// TODO: add not found check
			panic(err)
		}
		// maybe filtering some series ids after grouping that is result of filtering.
		// if not found, return empty series ids.
		mss.seriesIDs = seriesIDsAfterGrouping
		mss.groupingContext = groupingContext
	}
	fmt.Printf("final series id=%s\n", mss.seriesIDs)

	mss.highKeys = mss.seriesIDs.GetHighKeys()
}

func (mss *SplitSource) HasNext() bool {
	return mss.index < len(mss.highKeys)
}

func (mss *SplitSource) Next() spi.Split {
	if mss.index >= len(mss.highKeys) {
		return nil
	}
	highSeriesID := mss.highKeys[mss.index]
	lowSeriesIDsContainer := mss.seriesIDs.GetContainerAtIndex(mss.index)
	mss.index++
	return &ScanSplit{
		MinSeriesID:           lowSeriesIDsContainer.Minimum(),
		MaxSeriesID:           lowSeriesIDsContainer.Maximum(),
		HighSeriesID:          highSeriesID,
		LowSeriesIDsContainer: lowSeriesIDsContainer,
		tableScan:             mss.tableScan,
		groupingContext:       mss.groupingContext,
		ResultSet:             mss.resultSet,
	}
}

type RowsLookupVisitor struct {
	split *SplitSource
}

func NewRowLookupVisitor(split *SplitSource) *RowsLookupVisitor {
	return &RowsLookupVisitor{
		split: split,
	}
}

func (v *RowsLookupVisitor) Visit(context any, n tree.Node) any {
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		return v.visitComparisonExpression(context, node)
	case *tree.Cast:
		return node.Expression.Accept(context, v)
	}
	return nil
}

func (v *RowsLookupVisitor) visitComparisonExpression(_ any, node *tree.ComparisonExpression) (r any) {
	columnResult, ok := v.split.tableScan.filterResult[node.ID]
	if !ok {
		return nil
	}
	indexDB := v.split.partition.shard.IndexDB()
	seriesIDs, err := indexDB.GetSeriesIDsByTagValueIDs(columnResult.TagKeyID, columnResult.TagValueIDs)
	if err != nil {
		panic(err)
	}
	return seriesIDs
}
