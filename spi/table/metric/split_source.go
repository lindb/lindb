package metric

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/tsdb"
)

type MetricSplitSourceProvider struct {
	engine tsdb.Engine
}

func NewMetricSplitSourceProvider(engine tsdb.Engine) *MetricSplitSourceProvider {
	return &MetricSplitSourceProvider{
		engine: engine,
	}
}

func (msp *MetricSplitSourceProvider) CreateSplitSources(table spi.TableHandle, partitions []int,
	outputColumns []spi.ColumnMetadata, filter tree.Expression,
) (splits []spi.SplitSource) {
	var (
		metricTable *MetricTableHandle
		ok          bool
	)
	if metricTable, ok = table.(*MetricTableHandle); !ok {
		panic("not support table handle")
	}
	db, ok := msp.engine.GetDatabase(metricTable.Database)
	if !ok {
		panic(fmt.Errorf("%w: %s", constants.ErrDatabaseNotFound, metricTable.Database))
	}
	// find metric id(table id)
	metricID, err := db.MetaDB().GetMetricID(metricTable.Namespace, metricTable.Metric)
	if err != nil {
		panic(err)
	}
	// find table schema
	schema, err := db.MetaDB().GetSchema(metricID)
	if err != nil {
		panic(err)
	}

	var filterResult map[tree.NodeID]*flow.TagFilterResult
	if filter != nil {
		// lookup column if filter not nil
		lookup := NewColumnValuesLookVisitor(db, schema)
		lookup.filterResult = make(map[tree.NodeID]*flow.TagFilterResult)
		_ = filter.Accept(nil, lookup)
		filterResult = lookup.ResultSet()
		// TODO: check filter result if empty????
	}

	var shards []tsdb.Shard
	var families [][]tsdb.DataFamily
	for _, partition := range partitions {
		shard, ok := db.GetShard(models.ShardID(partition))
		if ok {
			dataFamilies := shard.GetDataFamilies(metricTable.Interval.Type(), metricTable.TimeRange)
			if len(dataFamilies) > 0 {
				shards = append(shards, shard)
				families = append(families, dataFamilies)
			}
		}
	}

	if len(shards) == 0 {
		panic(constants.ErrShardNotFound)
	}

	groupingTags := lo.Filter(schema.TagKeys, func(item tag.Meta, index int) bool {
		return lo.ContainsBy(outputColumns, func(column spi.ColumnMetadata) bool {
			return column.Name == item.Key
		})
	})
	lengthOfGroupByTagKeys := len(groupingTags)
	stargetCtx := &flow.StorageExecuteContext{}

	stargetCtx.GroupByTagKeyIDs = make([]tag.KeyID, lengthOfGroupByTagKeys)
	stargetCtx.GroupByTags = make(tag.Metas, lengthOfGroupByTagKeys)

	for idx, tagKey := range groupingTags {
		stargetCtx.GroupByTags[idx] = tagKey
		stargetCtx.GroupByTagKeyIDs[idx] = tagKey.ID
	}

	// init grouping tag value collection, need cache found grouping tag value id
	stargetCtx.GroupingTagValueIDs = make([]*roaring.Bitmap, lengthOfGroupByTagKeys)

	for i := range shards {
		splits = append(splits, NewMetricSplitSource(metricTable, db, stargetCtx, shards[i], metricID, schema, outputColumns, families[i], filter, filterResult))
	}

	return
}

type MetricSplitSource struct {
	where        tree.Expression
	shard        tsdb.Shard
	db           tsdb.Database
	seriesIDs    *roaring.Bitmap
	schema       *metric.Schema
	table        *MetricTableHandle
	filterResult map[tree.NodeID]*flow.TagFilterResult
	fields       field.Metas
	groupingTags tag.Metas
	families     []tsdb.DataFamily
	resultSet    []flow.FilterResultSet
	highKeys     []uint16
	index        int
	metricID     metric.ID

	shardCtx   *flow.ShardExecuteContext
	storageCtx *flow.StorageExecuteContext
}

// TODO: remove schema
func NewMetricSplitSource(table *MetricTableHandle, db tsdb.Database, storageCtx *flow.StorageExecuteContext, shard tsdb.Shard, metricID metric.ID,
	schema *metric.Schema, outputColumns []spi.ColumnMetadata, families []tsdb.DataFamily, where tree.Expression, filterResult map[tree.NodeID]*flow.TagFilterResult,
) *MetricSplitSource {
	return &MetricSplitSource{
		table:    table,
		db:       db,
		shard:    shard,
		families: families,
		metricID: metricID,
		schema:   schema,
		fields: lo.Filter(schema.Fields, func(item field.Meta, index int) bool {
			return lo.ContainsBy(outputColumns, func(column spi.ColumnMetadata) bool {
				return column.Name == item.Name.String()
			})
		}),
		groupingTags: lo.Filter(schema.TagKeys, func(item tag.Meta, index int) bool {
			return lo.ContainsBy(outputColumns, func(column spi.ColumnMetadata) bool {
				return column.Name == item.Key
			})
		}),
		where:        where,
		filterResult: filterResult,
		storageCtx:   storageCtx,
	}
}

func (mss *MetricSplitSource) Prepare() {
	var (
		seriesIDs *roaring.Bitmap
		err       error
		ok        bool
	)

	if mss.where == nil {
		// if where condition nil, find all series ids under metric
		seriesIDs, err = mss.shard.IndexDB().GetSeriesIDsForMetric(mss.metricID)
		if err != nil {
			panic(err)
		}
	} else {
		// find series ids based on where condition
		lookup := NewRowLookupVisitor(mss)
		if seriesIDs, ok = mss.where.Accept(nil, lookup).(*roaring.Bitmap); !ok {
			panic(constants.ErrSeriesIDNotFound)
		}
	}

	if seriesIDs == nil || seriesIDs.IsEmpty() {
		panic(constants.ErrSeriesIDNotFound)
	}

	mss.seriesIDs = roaring.New()

	for i := range mss.families {
		family := mss.families[i]
		// check family data if matches condition(series ids)
		resultSet, err := family.Filter(&flow.MetricScanContext{
			MetricID:                mss.metricID,
			SeriesIDs:               seriesIDs,
			SeriesIDsAfterFiltering: seriesIDs,
			Fields:                  mss.fields,
			TimeRange:               mss.table.TimeRange,
			StorageInterval:         mss.table.StorageInterval,
		})

		if !errors.Is(err, constants.ErrNotFound) && err != nil {
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
	mss.shardCtx = &flow.ShardExecuteContext{
		StorageExecuteCtx:       mss.storageCtx,
		SeriesIDsAfterFiltering: mss.seriesIDs,
	}

	if mss.groupingTags.Len() > 0 {
		fmt.Println("grouping tag keys")
		// if it has grouping, do group by tag keys, else just split series ids as batch first.
		mss.shard.IndexDB().GetGroupingContext(mss.shardCtx)
	}

	mss.highKeys = mss.seriesIDs.GetHighKeys()

	fmt.Printf("series id====%v,rs=%v\n", mss.seriesIDs.String(), mss.resultSet)
}

func (mss *MetricSplitSource) HasSplit() bool {
	return mss.index < len(mss.highKeys)
}

func (mss *MetricSplitSource) GetNextSplit() spi.Split {
	if mss.index >= len(mss.highKeys) {
		return nil
	}
	highSeriesID := mss.highKeys[mss.index]
	lowSeriesIDsContainer := mss.seriesIDs.GetContainerAtIndex(mss.index)
	mss.index++
	return &MetricScanSplit{
		MinSeriesID:           lowSeriesIDsContainer.Minimum(),
		MaxSeriesID:           lowSeriesIDsContainer.Maximum(),
		HighSeriesID:          highSeriesID,
		LowSeriesIDsContainer: lowSeriesIDsContainer,
		Fields:                mss.fields,
		GroupingTags:          mss.groupingTags,
		ResultSet:             mss.resultSet,
		ShardExecuteContext:   mss.shardCtx,
	}
}

type RowLookupVisitor struct {
	split *MetricSplitSource
}

func NewRowLookupVisitor(split *MetricSplitSource) *RowLookupVisitor {
	return &RowLookupVisitor{
		split: split,
	}
}

func (v *RowLookupVisitor) Visit(context any, n tree.Node) any {
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		return v.visitComparisonExpression(context, node)
	case *tree.Cast:
		return node.Expression.Accept(context, v)
	}
	return nil
}

func (v *RowLookupVisitor) visitComparisonExpression(_ any, node *tree.ComparisonExpression) (r any) {
	columnResult, ok := v.split.filterResult[node.ID]
	if !ok {
		return nil
	}
	indexDB := v.split.shard.IndexDB()
	seriesIDs, err := indexDB.GetSeriesIDsByTagValueIDs(columnResult.TagKeyID, columnResult.TagValueIDs)
	if err != nil {
		panic(err)
	}
	return seriesIDs
}

type ColumnValuesLookupVisitor struct {
	db     tsdb.Database
	schema *metric.Schema

	// result which after column condition metadata filter
	// set value in column search, the where clause condition that user input
	// first find all column values in where clause, then do column match
	filterResult map[tree.NodeID]*flow.TagFilterResult
}

func NewColumnValuesLookVisitor(db tsdb.Database, schema *metric.Schema) *ColumnValuesLookupVisitor {
	return &ColumnValuesLookupVisitor{
		db:           db,
		schema:       schema,
		filterResult: make(map[tree.NodeID]*flow.TagFilterResult),
	}
}

func (v *ColumnValuesLookupVisitor) Visit(context any, n tree.Node) any {
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		return v.visitComparisonExpression(context, node)
	case *tree.Cast:
		return node.Expression.Accept(context, v)
	case *tree.Identifier:
		return node.Value
	case *tree.StringLiteral:
		return node.Value
	}
	return nil
}

func (v *ColumnValuesLookupVisitor) ResultSet() map[tree.NodeID]*flow.TagFilterResult {
	return v.filterResult
}

func (v *ColumnValuesLookupVisitor) visitComparisonExpression(context any, node *tree.ComparisonExpression) (r any) {
	column := node.Left.Accept(context, v)
	columnValue := node.Right.Accept(context, v)

	columnName, ok := column.(string)
	if !ok {
		panic(fmt.Sprintf("column name '%v' not support type '%s'", column, reflect.TypeOf(column)))
	}

	tagMeta, ok := v.schema.TagKeys.Find(columnName)
	if !ok {
		panic(fmt.Errorf("%w, column name: %s", constants.ErrColumnNotFound, columnName))
	}
	tagKeyID := tagMeta.ID
	var tagValueIDs *roaring.Bitmap
	var err error
	switch val := columnValue.(type) {
	case string:
		// FIXME: impl other expr
		tagValueIDs, err = v.db.MetaDB().FindTagValueDsByExpr(tagKeyID, &tree.EqualsExpr{
			Name:  columnName,
			Value: val,
		})
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("value of column '%v' not support type '%s'", columnName, reflect.TypeOf(val)))
	}

	if tagValueIDs == nil || tagValueIDs.IsEmpty() {
		return nil
	}

	v.filterResult[node.ID] = &flow.TagFilterResult{
		TagKeyID:    tagKeyID,
		TagValueIDs: tagValueIDs,
	}
	return nil
}
