package metric

import (
	"errors"
	"fmt"
	"reflect"

	common_timeutil "github.com/lindb/common/pkg/timeutil"
	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
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

func (msp *MetricSplitSourceProvider) CreateSplitSources(database string, table spi.TableHandle, partitions []int, columns []spi.ColumnMetadata, filter tree.Expression) (splits []spi.SplitSource) {
	var (
		metricTable *MetricTableHandle
		ok          bool
	)
	if metricTable, ok = table.(*MetricTableHandle); !ok {
		panic("not support table handle")
	}
	db, ok := msp.engine.GetDatabase(database)
	if !ok {
		panic(fmt.Errorf("%w: %s", constants.ErrDatabaseNotFound, database))
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

	metricID, err := db.MetaDB().GetMetricID(metricTable.Namespace, metricTable.Metric)
	if err != nil {
		panic(err)
	}
	// check table schema
	schema, err := db.MetaDB().GetSchema(metricID)
	if err != nil {
		panic(err)
	}

	var filterResult map[tree.NodeID]*flow.TagFilterResult
	if filter != nil {
		// lookup column if filter not nil
		lookup := NewColumnLookVisitor(db, schema)
		lookup.filterResult = make(map[tree.NodeID]*flow.TagFilterResult)
		_ = filter.Accept(nil, lookup)
		filterResult = lookup.ResultSet()
		// TODO: check filter result if empty????
	}

	for i := range shards {
		splits = append(splits, NewMetricSplitSource(metricTable, db, shards[i], metricID, schema, columns, families[i], filter, filterResult))
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
	families     []tsdb.DataFamily
	resultSet    []flow.FilterResultSet
	highKeys     []uint16
	index        int
	metricID     metric.ID
}

// TODO: remove schema
func NewMetricSplitSource(table *MetricTableHandle, db tsdb.Database, shard tsdb.Shard, metricID metric.ID, schema *metric.Schema, columns []spi.ColumnMetadata, families []tsdb.DataFamily, where tree.Expression, filterResult map[tree.NodeID]*flow.TagFilterResult) *MetricSplitSource {
	return &MetricSplitSource{
		table:    table,
		db:       db,
		shard:    shard,
		families: families,
		metricID: metricID,
		schema:   schema,
		fields: lo.Filter(schema.Fields, func(item field.Meta, index int) bool {
			return lo.ContainsBy(columns, func(column spi.ColumnMetadata) bool {
				return column.Name == item.Name.String()
			})
		}),
		where:        where,
		filterResult: filterResult,
	}
}

func (mss *MetricSplitSource) Prepare() {
	var (
		seriesIDs *roaring.Bitmap
		err       error
		ok        bool
	)

	if mss.where == nil {
		seriesIDs, err = mss.shard.IndexDB().GetSeriesIDsForMetric(mss.metricID)
		if err != nil {
			panic(err)
		}
	} else {
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
		resultSet, err := family.Filter(&flow.MetricScanContext{
			MetricID:                mss.metricID,
			SeriesIDs:               seriesIDs,
			SeriesIDsAfterFiltering: seriesIDs,
			Fields:                  mss.fields,
			TimeRange:               mss.table.TimeRange,
			StorageInterval:         timeutil.Interval(10 * common_timeutil.OneSecond),
		})
		if !errors.Is(err, constants.ErrNotFound) && err != nil {
			panic(err)
		}

		// FIXME: group by
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
	if len(mss.resultSet) == 0 {
		panic(constants.ErrNotFound)
	}

	if mss.seriesIDs.IsEmpty() {
		panic(constants.ErrSeriesIDNotFound)
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
		ResultSet:             mss.resultSet,
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

type ColumnLookupVisitor struct {
	db     tsdb.Database
	schema *metric.Schema

	// result which after column condition metadata filter
	// set value in column search, the where clause condition that user input
	// first find all column values in where clause, then do column match
	filterResult map[tree.NodeID]*flow.TagFilterResult
}

func NewColumnLookVisitor(db tsdb.Database, schema *metric.Schema) *ColumnLookupVisitor {
	return &ColumnLookupVisitor{
		db:           db,
		schema:       schema,
		filterResult: make(map[tree.NodeID]*flow.TagFilterResult),
	}
}

func (v *ColumnLookupVisitor) Visit(context any, n tree.Node) any {
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

func (v *ColumnLookupVisitor) ResultSet() map[tree.NodeID]*flow.TagFilterResult {
	return v.filterResult
}

func (v *ColumnLookupVisitor) visitComparisonExpression(context any, node *tree.ComparisonExpression) (r any) {
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
