package metric

import (
	"fmt"
	"math"
	"sort"
	"strings"

	commonmodels "github.com/lindb/common/models"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/aggregation/function"
	"github.com/lindb/lindb/pkg/timeutil"
	protoCommonV1 "github.com/lindb/lindb/proto/gen/v1/common"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/tree"
)

type MetricMerger struct {
	groupAgg aggregation.GroupingAggregator
	stats    *commonmodels.NodeStats
	// field name -> aggregator spec
	// we will use it during intermediate tasks
	aggregatorSpecs map[string]*protoCommonV1.AggregatorSpec
	timeRange       timeutil.TimeRange
	interval        int64
}

func NewMetricMerger() spi.Merger {
	return &MetricMerger{
		aggregatorSpecs: make(map[string]*protoCommonV1.AggregatorSpec),
	}
}

// AddSplit implements spi.Merger
func (m *MetricMerger) AddSplit(split *spi.BinarySplit) {
	timeSeriesList := &protoCommonV1.TimeSeriesList{}
	// if err := timeSeriesList.Unmarshal(split.Page); err != nil {
	// 	panic(err)
	// }
	if len(timeSeriesList.FieldAggSpecs) == 0 {
		// if it gets empty aggregator spec(empty response), need ignore response.
		// if not ignore, will build empty group aggregator, and cannot aggregate real response data.
		return
	}

	for _, spec := range timeSeriesList.FieldAggSpecs {
		m.aggregatorSpecs[spec.FieldName] = spec
	}

	if m.groupAgg == nil {
		m.timeRange = timeutil.TimeRange{
			Start: timeSeriesList.Start,
			End:   timeSeriesList.End,
		}
		m.interval = timeSeriesList.Interval
		AggregatorSpecs := make(aggregation.AggregatorSpecs, len(timeSeriesList.FieldAggSpecs))
		for idx, aggSpec := range timeSeriesList.FieldAggSpecs {
			AggregatorSpecs[idx] = aggregation.NewAggregatorSpec(
				field.Name(aggSpec.FieldName),
				field.Type(aggSpec.FieldType),
			)
			for _, funcType := range aggSpec.FuncTypeList {
				AggregatorSpecs[idx].AddFunctionType(function.FuncType(funcType))
			}
		}
		m.groupAgg = aggregation.NewGroupingAggregator(
			timeutil.Interval(m.interval),
			1, // interval ratio is 1 when do merge result.
			m.timeRange,
			AggregatorSpecs,
		)
	}

	for _, ts := range timeSeriesList.TimeSeriesList {
		// if no field data, ignore this response
		if len(ts.Fields) == 0 {
			continue
		}
		fields := make(map[field.Name][]byte)
		for k, v := range ts.Fields {
			fmt.Printf("result field=%v\n", v)
			fields[field.Name(k)] = v
		}
		fmt.Println("group agg")
		m.groupAgg.Aggregate(series.NewGroupedIterator(ts.Tags, fields))
	}
}

// GetOutputPage implements spi.Merger
func (m *MetricMerger) GetOutputPage() *spi.Page {
	if m.groupAgg == nil {
		return nil
	}
	// groupIts := m.groupAgg.ResultSet()
	page := spi.NewPage()

	// for _, it := range groupIts {
	// 	for it.HasNext() {
	// 		fieldSeries := it.Next()
	// 		fieldName := fieldSeries.FieldName()
	// 		fieldType := fieldSeries.FieldType()
	// 		f := fields.NewDynamicField(fieldType, e.timeRange.Start, e.interval, e.pointCount)
	// 		e.fieldStore[fieldName] = f
	// 		fmt.Printf("fiiiii==%s\n", fieldName)
	// 		f.SetValue(fieldSeries)
	// 		var fieldValues *collections.FloatArray
	// 		ok := false
	// 		for fieldSeries.HasNext() {
	// 			startTime, it := fieldSeries.Next()
	// 			fmt.Printf("set v1=%v\n", it)
	// 			if it == nil {
	// 				continue
	// 			}
	// 			for it.HasNext() {
	// 				pIt := it.Next()
	// 				aggType := pIt.AggType()
	// 				fieldValues, ok = f.fields[aggType]
	// 				if !ok {
	// 					fieldValues = collections.NewFloatArray(f.capacity)
	// 					f.fields[aggType] = fieldValues
	// 				}
	// 				for pIt.HasNext() {
	// 					slot, val := pIt.Next()
	// 					idx := ((int64(slot)*f.interval + startTime) - f.startTime) / f.interval
	// 					fmt.Printf("value 2=%v\n", val)
	// 					fieldValues.SetValue(int(idx), val)
	// 				}
	// 			}
	// 		}
	// 	}
	// }
	//
	// page.AppendColumn(groupIts)
	return page
}

// makeResultSet makes final result set from time series event(GroupedIterators).
// TODO: can opt use stream, leaf node need return grouping if completed.
func (m *MetricMerger) makeResultSet() (resultSet *commonmodels.ResultSet, err error) {
	// makeResultStartTime := time.Now()
	orderBy, err := m.buildOrderBy()
	if err != nil {
		return nil, err
	}

	// statement := exec.preparedStatement.Statement
	resultSet = new(commonmodels.ResultSet)
	// TODO: merge stats for cross idc query?
	// groupByKeys := statement.GroupBy
	// groupByKeysLength := 0
	// len(groupByKeys)
	fieldsMap := make(map[string]struct{})
	timeRange := m.timeRange
	interval := m.interval
	if m.groupAgg != nil {
		groupIts := m.groupAgg.ResultSet()
		selectItems := m.getSelectItems()
		fmt.Println(selectItems)
		fmt.Println(groupIts)
		for _, it := range groupIts {
			// TODO: reuse expression??
			expression := aggregation.NewExpression(
				timeRange,
				interval,
				selectItems,
			)
			// do expression eval
			expression.Eval(it)

			rs := expression.ResultSet()
			fmt.Printf("group result=%v\n", rs)
			// result order by/limit
			orderBy.Push(aggregation.NewOrderByRow(it.Tags(), rs))
		}

		rows := orderBy.ResultSet()
		for _, row := range rows {
			var tags map[string]string
			tagValues, fields := row.ResultSet()
			// if groupByKeysLength > 0 {
			// 	tagValues := tag.SplitTagValues(tagValues)
			// 	if groupByKeysLength != len(tagValues) {
			// 		// if tag values not match group by tag keys, ignore this time series
			// 		continue
			// 	}
			// 	// build group by tags for final result
			// 	tags = make(map[string]string)
			// 	for idx, tagKey := range groupByKeys {
			// 		tags[tagKey] = tagValues[idx]
			// 	}
			// }
			timeSeries := commonmodels.NewSeries(tags, tagValues)
			resultSet.AddSeries(timeSeries)
			fmt.Println("add series")

			// having := ctx.Deps.Statement.Having
			notHavingSlots := make(map[int]struct{})
			// slotValues := make(map[int]map[string]float64)

			// if having != nil {
			// 	for fieldName, values := range fields {
			// 		if values == nil {
			// 			continue
			// 		}
			// 		it := values.NewIterator()
			// 		for it.HasNext() {
			// 			slot, val := it.Next()
			// 			if math.IsNaN(val) {
			// 				continue
			// 			}
			// 			if v, ok := slotValues[slot]; ok {
			// 				v[fieldName] = val
			// 			} else {
			// 				slotValues[slot] = map[string]float64{fieldName: val}
			// 			}
			// 		}
			// 	}
			// 	// calc and fill
			// 	if len(slotValues) > 0 {
			// 		calc := sql.NewCalc(having)
			// 		for slot, fieldValue := range slotValues {
			// 			result, err := calc.CalcExpr(fieldValue)
			// 			if err != nil {
			// 				return resultSet, err
			// 			}
			// 			if r, ok := result.(bool); !ok {
			// 				return resultSet, fmt.Errorf("expected CalcExpr bool result got %v", reflect.TypeOf(result))
			// 			} else if !r {
			// 				notHavingSlots[slot] = struct{}{}
			// 			}
			// 		}
			// 	}
			// }

			fmt.Println(fields)
			for fieldName, values := range fields {
				if values == nil {
					fmt.Println("no values")
					continue
				}

				points := commonmodels.NewPoints()
				it := values.NewIterator()
				for it.HasNext() {
					slot, val := it.Next()
					if math.IsNaN(val) {
						// TODO: need check
						continue
					}
					fmt.Println("add point.....................")
					if _, ok := notHavingSlots[slot]; ok {
						continue
					}
					fmt.Println("add point.....................")
					points.AddPoint(timeutil.CalcTimestamp(timeRange.Start, slot, timeutil.Interval(interval)), val)
				}
				timeSeries.AddField(fieldName, points)
				fieldsMap[fieldName] = struct{}{}
			}
		}
	}

	sort.Slice(resultSet.Series, func(i, j int) bool {
		return resultSet.Series[i].TagValues < resultSet.Series[j].TagValues
	})

	// resultSet.MetricName = statement.MetricName
	// resultSet.GroupBy = statement.GroupBy
	for fName := range fieldsMap {
		resultSet.Fields = append(resultSet.Fields, fName)
	}
	resultSet.StartTime = timeRange.Start
	resultSet.EndTime = timeRange.End
	resultSet.Interval = interval

	// if ctx.stats != nil {
	// 	now := time.Now()
	// 	ctx.stats.Node = ctx.Deps.CurrentNode.Indicator()
	// 	ctx.stats.End = now.UnixNano()
	// 	ctx.stats.TotalCost = now.Sub(ctx.startTime).Nanoseconds()
	//
	// 	ctx.stats.Stages = append(ctx.stats.Stages, &commonmodels.StageStats{
	// 		Identifier: "Expression",
	// 		Start:      makeResultStartTime.UnixNano(),
	// 		End:        now.UnixNano(),
	// 		Cost:       now.Sub(makeResultStartTime).Nanoseconds(),
	// 		State:      tracker.CompleteState.String(),
	// 		Async:      false,
	// 	})
	// 	resultSet.Stats = ctx.stats
	// }
	return resultSet, nil
}

// buildOrderBy builds order by container.
func (m *MetricMerger) buildOrderBy() (aggregation.OrderBy, error) {
	// statement := exec.preparedStatement.Statement
	// build order by items if need do order by query
	// orderByExprs := statement.OrderByItems
	// if len(orderByExprs) == 0 {
	// use default limiter
	return aggregation.NewResultLimiter(100), nil
	// }
	// var orderByItems []*aggregation.OrderByItem
	// fields := exec.aggregatorSpecs
	// for _, orderBy := range orderByExprs {
	// 	expr := orderBy.(*tree.OrderByExpr)
	// 	funcType := function.Unknown
	// 	var fieldName string
	// 	switch e := expr.Expr.(type) {
	// 	case *tree.FieldExpr:
	// 		aggSpec, ok := fields[e.Name]
	// 		if ok {
	// 			funcType = field.Type(aggSpec.FieldType).GetOrderByFunc()
	// 			fieldName = e.Name
	// 		}
	// 	case *tree.CallExpr:
	// 		funcType = e.FuncType
	// 		fieldName = e.Params[0].Rewrite()
	// 	}
	// 	if funcType == function.Unknown {
	// 		return nil, errors.New("cannot parse order by function")
	// 	}
	// 	orderByItems = append(orderByItems, &aggregation.OrderByItem{
	// 		Expr:     expr,
	// 		Name:     fieldName,
	// 		FuncType: funcType,
	// 		Desc:     expr.Desc,
	// 	})
	// }
	// return aggregation.NewTopNOrderBy(orderByItems, statement.Limit), nil
}

// getSelectItems returns select field items.
func (m *MetricMerger) getSelectItems() []tree.Expr {
	// statement := ctx.Deps.Statement
	// selectItems := statement.SelectItems
	selectItems := []tree.Expr{}
	// if statement.AllFields {
	// if select all fields, read field names from aggregator
	allAggFields := m.groupAgg.Fields()
	selectItems = []tree.Expr{}
	isHistogram := false
	for _, fieldName := range allAggFields {
		if strings.HasPrefix(string(fieldName), "__bucket_") {
			// filter histogram raw field
			isHistogram = true
			continue
		}
		selectItems = append(selectItems, &tree.SelectItem2{Expr: &tree.FieldExpr{Name: fieldName.String()}})
	}
	if isHistogram {
		// add histogram functions
		addQuantileFn := func(as string, num float64) {
			selectItems = append(selectItems, &tree.SelectItem2{
				Expr:  &tree.CallExpr{FuncType: function.Quantile, Params: []tree.Expr{&tree.NumberLiteral{Val: num}}},
				Alias: as,
			})
		}
		addQuantileFn("p99", 0.99)
		addQuantileFn("p95", 0.95)
		addQuantileFn("p90", 0.90)
		addQuantileFn("mean", 0.50)
	}
	// }
	return selectItems
}
