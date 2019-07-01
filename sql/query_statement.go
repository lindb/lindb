package sql

import (
	"fmt"

	"github.com/eleme/lindb/pkg/proto"
	parser "github.com/eleme/lindb/sql/grammar"
	"github.com/eleme/lindb/sql/util"

	"strconv"
	"strings"
	"time"
)

type Unit struct {
	unit int32
}

type QueryStatement struct {
	measurement         string
	startTime           int64
	endTime             int64
	condition           *proto.Condition
	fieldExprList       []proto.FieldExpr
	orderByExpr         *proto.Expr
	desc                bool
	limit               int32
	conditionAggregates map[*proto.Condition]map[util.AggregatorUnit]Unit
	groupBy             map[string]bool
	interval            int64
	fieldID             int
}

// NewDefaultQueryStatement build lindb default query statement
func NewDefaultQueryStatement() *QueryStatement {
	return &QueryStatement{
		limit:               20,
		fieldID:             1,
		condition:           new(proto.Condition),
		fieldExprList:       make([]proto.FieldExpr, 0, 10),
		groupBy:             make(map[string]bool),
		conditionAggregates: make(map[*proto.Condition]map[util.AggregatorUnit]Unit),
	}
}

// Parse lindb parse query sql
func (qs *QueryStatement) Parse(ctx *parser.QueryStmtContext) {
	// parse measurement like 'table'
	if ctx.FromClause() != nil {
		clauseContext := ctx.FromClause().(*parser.FromClauseContext)
		qs.measurement = util.GetStringValue(clauseContext.MetricName().GetText())
	}
	// parse select fields
	if ctx.Fields() != nil {
		fieldsContext := ctx.Fields().(*parser.FieldsContext)
		qs.parseFields(fieldsContext)
	}
	//parse where clause
	if ctx.WhereClause() != nil {
		whereClause := ctx.WhereClause().(*parser.WhereClauseContext)
		qs.parseClause(whereClause)
	}
	//parse group by
	if ctx.GroupByClause() != nil {
		groupBy := ctx.GroupByClause().(*parser.GroupByClauseContext)
		qs.parseGroupBy(groupBy)
	}
	//parse order by
	if ctx.OrderByClause() != nil {
		orderBy := ctx.OrderByClause().(*parser.OrderByClauseContext)
		qs.parseOrderBy(orderBy)
	}
	//parse limit
	if ctx.LimitClause() != nil {
		limitClause := ctx.LimitClause().(*parser.LimitClauseContext)
		qs.parseLimit(limitClause)
	}
}

// parseLimit parse lindb sql limit
func (qs *QueryStatement) parseLimit(limitClause *parser.LimitClauseContext) {
	if nil != limitClause {
		limit, _ := strconv.ParseInt(limitClause.L_INT().GetText(), 10, 32)
		qs.limit = int32(limit)
	}
}

// build build lindb query statement
func (qs *QueryStatement) build() *proto.Query {
	query := new(proto.Query)
	query.Measurement = qs.measurement
	query.Interval = qs.interval
	now := util.NowTimestamp()
	timeRange := new(proto.TimeRange)
	timeRange.StartTime = qs.startTime
	timeRange.EndTime = qs.endTime

	if timeRange.StartTime <= 0 {
		timeRange.StartTime = now - util.OneHour
	}
	if timeRange.EndTime <= 0 {
		timeRange.EndTime = now
	}
	if timeRange.EndTime < timeRange.StartTime {
		panic(fmt.Sprintln("start time can not be larger than end time"))
	}
	query.TimeRange = timeRange
	query.Condition = qs.condition
	groupByExpr := new(proto.GroupByExpr)

	for group := range qs.groupBy {
		groupByExpr.GroupBy = append(groupByExpr.GroupBy, group)
	}
	query.GroupByExpr = groupByExpr
	for entry := range qs.conditionAggregates {
		for aggregatorUnit := range qs.conditionAggregates[entry] {
			builder := new(proto.ConditionAggregator)
			builder.Condition = entry
			builder.Field = aggregatorUnit.GetField()
			if aggregatorUnit.GetDownSampling() != nil {
				builder.DownSampling = aggregatorUnit.GetDownSampling().String()
			}
			if aggregatorUnit.GetAggregator() != nil {
				builder.Aggregator = aggregatorUnit.GetAggregator().String()
			}
			builder.UnitId = qs.conditionAggregates[entry][aggregatorUnit].unit
			query.ConditionAggregators = append(query.ConditionAggregators, builder)
		}
	}
	if qs.orderByExpr != nil {
		orderByExprBuilder := new(proto.OrderByExpr)
		orderByExprBuilder.Expr = qs.orderByExpr
		orderByExprBuilder.Desc = qs.desc
		query.OrderBy = orderByExprBuilder
		query.Limit = qs.limit
	}
	for i := range qs.fieldExprList {
		query.FieldExprList = append(query.FieldExprList, &qs.fieldExprList[i])
	}
	return query
}

// parseClause lindb sql parse where clause
func (qs *QueryStatement) parseClause(whereClauseContext *parser.WhereClauseContext) {
	if nil == whereClauseContext {
		return
	}
	clause := whereClauseContext.ClauseBooleanExpr().(*parser.ClauseBooleanExprContext)
	allClause := clause.AllClauseBooleanExpr()
	switch {
	case len(allClause) > 0:
		var exprs = make([]parser.ClauseBooleanExprContext, 0, len(allClause))
		for i := range allClause {
			expr := allClause[i].(*parser.ClauseBooleanExprContext)
			exprs = append(exprs, *expr)
		}
		qs.parseConditionClause(exprs)
	case nil != clause.TimeExpr():
		expr := clause.TimeExpr().(*parser.TimeExprContext)
		qs.parseTimeExpr(expr)
	case nil != clause.TagBooleanExpr():
		expr := clause.TagBooleanExpr().(*parser.TagBooleanExprContext)
		parseTagExpr(expr, qs.condition)
	}
	if qs.condition.Operator == proto.LogicOperator_UNKNOWN {
		qs.condition.Operator = proto.LogicOperator_AND
	}
}

// parseConditionClause parse lindb sql condition clause
func (qs *QueryStatement) parseConditionClause(clauseList []parser.ClauseBooleanExprContext) {
	if clauseList == nil {
		return
	}
	if 2 == len(clauseList) {
		for i := 0; i < 2; i++ {
			clause := clauseList[i]
			if nil != clause.TimeExpr() {
				timeExpr := clause.TimeExpr().(*parser.TimeExprContext)
				qs.parseTimeExpr(timeExpr)
			} else if nil != clause.TagBooleanExpr() {
				exprContext := clause.TagBooleanExpr().(*parser.TagBooleanExprContext)
				parseTagExpr(exprContext, qs.condition)
			}
		}
	}
}

// parseTimeExpr lindb Parse time range condition
func (qs *QueryStatement) parseTimeExpr(ctx *parser.TimeExprContext) {
	if ctx == nil {
		return
	}
	timeBooleanExprContexts := ctx.AllTimeBooleanExpr()
	for i := range timeBooleanExprContexts {
		timeExpr := timeBooleanExprContexts[i].(*parser.TimeBooleanExprContext)
		timestamp := timeExpr.Ident()
		var times int64
		if timestamp != nil {
			times = util.ParseTimestamp(util.GetStringValue(timestamp.GetText()))
		}
		iNowExpr := timeExpr.NowExpr()
		if iNowExpr != nil {
			now := time.Now().Unix() * 1000
			nowExpr := iNowExpr.(*parser.NowExprContext)
			if nowExpr.DurationLit() != nil {
				durationList := nowExpr.DurationLit().(*parser.DurationLitContext)
				duration := parseDuration(durationList)
				times = now + duration
			} else {
				times = now
			}
		}
		iOpCtx := timeExpr.BoolExprBinaryOperator()
		opCtx := iOpCtx.(*parser.BoolExprBinaryOperatorContext)
		if opCtx.T_GREATER() != nil || opCtx.T_GREATEREQUAL() != nil {
			qs.startTime = times
		}
		if opCtx.T_LESS() != nil || opCtx.T_LESSEQUAL() != nil {
			qs.endTime = times
		}

	}
}

// parseFields lindb sql parse a list of one or more fields
func (qs *QueryStatement) parseFields(ctx *parser.FieldsContext) {
	if ctx == nil {
		return
	}
	fieldContexts := ctx.AllField()
	if len(fieldContexts) == 0 {
		return
	}
	for i := range fieldContexts {
		fieldContext := fieldContexts[i].(*parser.FieldContext)
		qs.parseField(fieldContext)
	}
}

// parseOrderBy lindb sql parse order by
func (qs *QueryStatement) parseOrderBy(orderByClause *parser.OrderByClauseContext) {
	if orderByClause != nil {
		iSortFieldsContext := orderByClause.SortFields()
		if iSortFieldsContext != nil {
			sortFieldsContext := iSortFieldsContext.(*parser.SortFieldsContext)
			sortFieldContextList := sortFieldsContext.AllSortField()
			if len(sortFieldContextList) > 0 {
				sortFieldContext := sortFieldContextList[0].(*parser.SortFieldContext)
				if len(sortFieldContext.AllT_ASC()) > 0 {
					qs.desc = false
				} else {
					qs.desc = true
				}

				expr := new(proto.Expr)
				exprContext := sortFieldContext.Expr().(*parser.ExprContext)
				qs.parseExpr(exprContext)
				//todo
				//if new(proto.Expr_Ref).Ref == expr.GetRef() {
				//	ref := expr.GetRef().RefName
				//	i := 0
				//	length := len(qs.fieldExprList)
				//	for ; i < length; i++ {
				//		fieldExpr := qs.fieldExprList[0]
				//		if strings.EqualFold(fieldExpr.Alias, ref) {
				//			qs.orderByExpr = fieldExpr.GetExpr()
				//			return
				//		}
				//	}
				//}
				qs.orderByExpr = expr
			}
		}
	}
}

// parseField lindb sql parse a single field
func (qs *QueryStatement) parseField(ctx *parser.FieldContext) {
	exprCtx := ctx.Expr().(*parser.ExprContext)
	field := new(proto.FieldExpr)
	expr := new(proto.Expr)
	qs.parseExpr(exprCtx)
	if ctx.Alias() != nil {
		alias := ctx.Alias().(*parser.AliasContext)
		field.Alias = parseAlias(alias)
	} else {
		field.Alias = exprCtx.GetText()
	}
	field.Expr = expr
	qs.fieldExprList = append(qs.fieldExprList, *field)
}

// parseAlias lindb sql parse the "as alias" alias for fields
func parseAlias(ctx *parser.AliasContext) string {
	alias := ctx.Ident().GetText()
	return util.GetStringValue(alias)
}

func (qs *QueryStatement) parseExpr(ctx *parser.ExprContext) {
	if nil == ctx {
		return
	}
	switch {
	case nil != ctx.ExprAtom():
		atom := ctx.ExprAtom().(*parser.ExprAtomContext)
		qs.parseExprAtom(atom)
	case nil != ctx.ExprFunc():
		fun := ctx.ExprFunc().(*parser.ExprFuncContext)
		qs.parseCall(fun)
	case nil != ctx.DurationLit():
		lin := ctx.DurationLit().(*parser.DurationLitContext)
		interval := parseDuration(lin)
		longVal := new(proto.LongExpr)
		longVal.Value = interval
		//todo set long value
	case nil != ctx.AllExpr():
		exprContexts := ctx.AllExpr()
		size := len(exprContexts)
		switch size {
		case 2:
			binaryExp := new(proto.BinaryExpr)
			left := new(proto.Expr)
			exprContext0 := exprContexts[0].(*parser.ExprContext)
			qs.parseExpr(exprContext0)
			binaryExp.Left = left
			var op proto.Operator
			switch {
			case ctx.T_ADD() != nil:
				op = proto.Operator_ADD
			case ctx.T_SUB() != nil:
				op = proto.Operator_SUB
			case ctx.T_MUL() != nil:
				op = proto.Operator_MUL
			case ctx.T_DIV() != nil:
				op = proto.Operator_DIV
			default:
				panic(fmt.Sprintln("Unknown operator type in expr"))
			}
			binaryExp.Op = op
			right := new(proto.Expr)
			exprContext1 := exprContexts[1].(*parser.ExprContext)
			qs.parseExpr(exprContext1)
			binaryExp.Right = right
			//todo set binary value
			return
		case 1:
			exprContext := exprContexts[0].(*parser.ExprContext)
			qs.parseExpr(exprContext)
			return
		default:
			panic(fmt.Sprintln("Unknown expr type"))
		}
	}
}

func (qs *QueryStatement) parseExprAtom(ctx *parser.ExprAtomContext) {
	switch {
	case ctx.Ident() != nil:
		valRef := new(proto.ValRefExpr)
		field := util.GetStringValue(ctx.Ident().GetText())
		var currentFieldID int32
		if ctx.IdentFilter() != nil {
			filterContext, _ := ctx.IdentFilter().(*parser.IdentFilterContext)
			condition := parseFilter(filterContext)
			currentFieldID = qs.addField(condition, *util.NewAggregatorUnit(field, nil, nil))
		} else {
			currentFieldID = qs.addField(new(proto.Condition), *util.NewAggregatorUnit(field, nil, nil))
		}
		valRef.RefName = strconv.Itoa(int(currentFieldID))
		//todo set ref value
	case ctx.IntNumber() != nil:
		longVal := new(proto.LongExpr)
		dec, _ := strconv.ParseInt(ctx.IntNumber().GetText(), 10, 64)
		longVal.Value = dec
		//todo set long value
	case ctx.DecNumber() != nil:
		doubleVal := new(proto.DoubleExpr)
		dec, _ := strconv.ParseFloat(ctx.DecNumber().GetText(), 64)
		doubleVal.Value = dec
		//todo set double value
	default:
		panic(fmt.Sprintf("Unknown expr atom type"))
	}
}

// addField parse lindb sql field
func (qs *QueryStatement) addField(condition *proto.Condition, unit util.AggregatorUnit) int32 {
	aggregatorUnits := qs.conditionAggregates[condition]
	if aggregatorUnits == nil {
		aggregatorUnits = make(map[util.AggregatorUnit]Unit)
	}
	id := aggregatorUnits[unit]
	if id.unit <= 0 {
		id.unit = int32(qs.fieldID)
		aggregatorUnits[unit] = id
		qs.fieldID++
	}
	qs.conditionAggregates[condition] = aggregatorUnits
	return id.unit
}

// parseCall lindb sql parse a function call
func (qs *QueryStatement) parseCall(ctx *parser.ExprFuncContext) {
	name := util.GetStringValue(ctx.Ident().GetText())
	if !util.IsDownSamplingOrAggregator(name) {
		callExpr := new(proto.CallExpr)
		callExpr.Name = name
		if ctx.ExprFuncParams() != nil {
			funcParams := ctx.ExprFuncParams().(*parser.ExprFuncParamsContext)
			funcParam := funcParams.AllFuncParam()
			for i := range funcParam {
				param := funcParam[i].(*parser.FuncParamContext)
				callParam := new(proto.Expr)
				qs.parseCallParam(param)
				callExpr.Args = append(callExpr.Args, callParam)
			}
		}
		//todo expr call
		return
	}
	varRefBuilder := new(proto.ValRefExpr)
	currentFieldID := qs.parseDownSamplingOrAggregator(ctx)
	varRefBuilder.RefName = strconv.Itoa(int(currentFieldID))
	//todo set ref value
}

// parseDownSamplingOrAggregator parse lindb down sampling or aggregator
func (qs *QueryStatement) parseDownSamplingOrAggregator(ctx *parser.ExprFuncContext) int32 {
	var values = make(map[int]string)
	filterCondition := getCallNames(ctx, values)
	size := len(values)
	var currentFieldID int32
	var aggregatorUnit util.AggregatorUnit
	switch size {
	case 2:
		fun := values[0]
		function := util.ValueOf(fun)
		if strings.HasPrefix(fun, util.DownSampling) {
			aggregatorUnit = *util.NewAggregatorUnit(values[1], &function, nil)
		} else {
			aggregatorUnit = *util.NewAggregatorUnit(values[1], nil, &function)
		}
	case 3:
		fieldName := values[2]
		fun1 := util.ValueOf(values[1])
		fun0 := util.ValueOf(values[0])
		aggregatorUnit = *util.NewAggregatorUnit(fieldName, &fun1, &fun0)
	case 1:
		aggregatorUnit = *util.NewAggregatorUnit(values[0], nil, nil)
	default:
		panic(fmt.Sprintln("sql is not valid"))
	}

	if filterCondition == nil {
		currentFieldID = qs.addField(new(proto.Condition), aggregatorUnit)
	} else {
		currentFieldID = qs.addField(filterCondition, aggregatorUnit)
	}
	return currentFieldID
}

func getCallNames(ctx *parser.ExprFuncContext, values map[int]string) *proto.Condition {
	var filterCondition *proto.Condition
	callName := util.GetStringValue(ctx.Ident().GetText())
	valuesLen := len(values)
	values[valuesLen] = callName
	if ctx.ExprFuncParams() != nil {
		paramsContext := ctx.ExprFuncParams().(*parser.ExprFuncParamsContext)
		funcParams := paramsContext.AllFuncParam()
		for i := range funcParams {
			param := funcParams[i].(*parser.FuncParamContext)
			if param.Expr() == nil {
				continue
			}
			exprContext := param.Expr().(*parser.ExprContext)
			if exprContext.ExprFunc() != nil {
				exprFuncContext := exprContext.ExprFunc().(*parser.ExprFuncContext)
				getCallNames(exprFuncContext, values)
			} else if exprContext.ExprAtom() != nil {
				atomContext := exprContext.ExprAtom().(*parser.ExprAtomContext)
				identContext := atomContext.Ident()
				if identContext != nil {
					valuesLen = len(values)
					values[valuesLen] = util.GetStringValue(identContext.GetText())
				}
				if atomContext.IdentFilter() != nil {
					filterContext := atomContext.IdentFilter().(*parser.IdentFilterContext)
					filterCondition = parseFilter(filterContext)
				}
			}
		}
	}
	return filterCondition
}

// parseFilter parse lindb sql filter
func parseFilter(identFilterContext *parser.IdentFilterContext) *proto.Condition {
	if identFilterContext == nil {
		return nil
	}
	fieldCondition := new(proto.Condition)
	expr := identFilterContext.TagBooleanExpr().(*parser.TagBooleanExprContext)
	parseTagExpr(expr, fieldCondition)
	if proto.LogicOperator_UNKNOWN == fieldCondition.Operator {
		fieldCondition.Operator = proto.LogicOperator_AND
	}
	return fieldCondition
}

// parseCallParam parse lindb sql call expr param
func (qs *QueryStatement) parseCallParam(ctx *parser.FuncParamContext) {
	if nil != ctx.TagBooleanExpr() {
		conditionBuilder := new(proto.Condition)
		exprContext := ctx.TagBooleanExpr().(*parser.TagBooleanExprContext)
		parseTagExpr(exprContext, conditionBuilder)
	} else {
		exprContext := ctx.Expr().(*parser.ExprContext)
		qs.parseExpr(exprContext)
	}

}

// parseTagExpr parse lindb sql tag
func parseTagExpr(ctx *parser.TagBooleanExprContext, condition *proto.Condition) {
	if nil == ctx {
		return
	}
	tagKeyCtx := ctx.TagKey()
	if nil != tagKeyCtx {
		filter := parseTag(ctx)
		condition.TagFilters = append(condition.TagFilters, filter)
	} else {
		tagCtxList := ctx.AllTagBooleanExpr()
		size := len(tagCtxList)
		tagCtx0 := tagCtxList[0].(*parser.TagBooleanExprContext)
		if 2 == size {
			var logicOp proto.LogicOperator
			if nil != ctx.T_AND() {
				logicOp = proto.LogicOperator_AND
			} else if nil != ctx.T_OR() {
				logicOp = proto.LogicOperator_OR
			}

			tagCtx1 := tagCtxList[1].(*parser.TagBooleanExprContext)
			if proto.LogicOperator_OR == logicOp {
				orCondition := new(proto.Condition)
				orCondition.Operator = logicOp
				parseTagExpr(tagCtx0, orCondition)
				parseTagExpr(tagCtx1, orCondition)
				condition.Condition = append(condition.Condition, orCondition)
			} else {
				if nil != tagCtx0.TagKey() {
					filter := parseTag(tagCtx0)
					condition.TagFilters = append(condition.TagFilters, filter)
				} else {
					parseTagExpr(tagCtx0, condition)
				}

				if nil != tagCtx1.TagKey() {
					filter := parseTag(tagCtx1)
					condition.TagFilters = append(condition.TagFilters, filter)
				} else {
					parseTagExpr(tagCtx1, condition)
				}
			}
			if proto.LogicOperator_UNKNOWN == condition.Operator {
				condition.Operator = logicOp
			}
		} else if 1 == size {
			condition2 := new(proto.Condition)
			parseTagExpr(tagCtx0, condition2)
			if proto.LogicOperator_AND == condition2.Operator &&
				0 == len(condition2.Condition) &&
				proto.LogicOperator_AND == condition.Operator {
				for i := range condition2.TagFilters {
					tagFilter := condition2.TagFilters[i]
					condition.TagFilters = append(condition.TagFilters, tagFilter)
				}
			} else {
				condition.Condition = append(condition.Condition, condition2)
			}
		}

	}
}

// parseTag parse lindb sql tag filter
func parseTag(ctx *parser.TagBooleanExprContext) *proto.TagFilter {
	tagKeyCtx := ctx.TagKey()
	var op proto.Operator
	switch {
	case nil != ctx.T_EQUAL():
		op = proto.Operator_EQUAL
	case nil != ctx.T_LIKE() || nil != ctx.T_REGEXP():
		op = proto.Operator_LIKE
	case nil != ctx.T_NOTEQUAL() || nil != ctx.T_NOTEQUAL2():
		op = proto.Operator_NOT_EQUAL
	case nil != ctx.T_NOT() && nil != ctx.T_IN():
		op = proto.Operator_NOT_IN
	case nil == ctx.T_NOT() && nil != ctx.T_IN():
		op = proto.Operator_IN
	default:
		panic(fmt.Sprintf("unknown filter operator"))
	}

	filter := new(proto.TagFilter)
	tagKey := util.GetStringValue(tagKeyCtx.GetText())
	filter.TagKey = tagKey
	filter.Op = op
	ctx.TagValue()
	if proto.Operator_IN == op || proto.Operator_NOT_IN == op {
		if ctx.TagValueList() != nil {
			tagValueListContext, _ := ctx.TagValueList().(*parser.TagValueListContext)
			tagValueContexts := tagValueListContext.AllTagValue()
			if len(tagValueContexts) > 0 {
				var tagValues = make([]string, 0)
				for i := range tagValueContexts {
					item := tagValueContexts[i].(*parser.TagValueContext)
					tagValue := item.Ident().GetText()
					tagValues = append(tagValues, tagValue)
				}
				filter.TagValueItems = tagValues
			} else {
				panic(fmt.Sprintf(" the tag values for operator %s must be empty", op))
			}
		} else {
			panic(fmt.Sprintf(" the tag values for operator %s must be list", op))
		}
	} else {
		tagValueContext := ctx.TagValue().(*parser.TagValueContext)
		tagValue := util.GetStringValue(tagValueContext.Ident().GetText())
		filter.TagValue = tagValue
	}
	return filter
}

// parseDuration lindb sql parse a  time duration from a string
func parseDuration(ctx *parser.DurationLitContext) int64 {
	duration, _ := strconv.ParseInt(ctx.IntNumber().GetText(), 10, 64)
	var result int64
	if ctx.IntervalItem() == nil {
		return result
	}
	unit, _ := ctx.IntervalItem().(*parser.IntervalItemContext)
	switch {
	case nil != unit.T_SECOND():
		result = duration * util.OneSeconds
	case nil != unit.T_MINUTE():
		result = duration * util.OneMinute
	case nil != unit.T_HOUR():
		result = duration * util.OneHour
	case nil != unit.T_DAY():
		result = duration * util.OneDay
	case nil != unit.T_WEEK():
		result = duration * util.OneWeek
	case nil != unit.T_MONTH():
		result = duration * util.OneMonth
	case nil != unit.T_YEAR():
		result = duration * util.OneYear
	}
	return result
}

// parseGroupBy parse lindb sql group by
func (qs *QueryStatement) parseGroupBy(ctx *parser.GroupByClauseContext) {
	if nil == ctx {
		return
	}
	if ctx.Dimensions() == nil {
		return
	}
	dimensions, _ := ctx.Dimensions().(*parser.DimensionsContext)
	var dimensionContextList = dimensions.AllDimension()
	for i := range dimensionContextList {
		dimension := dimensionContextList[i].(*parser.DimensionContext)
		tagKey := dimension.Ident()
		if nil != tagKey {
			qs.groupBy[util.GetStringValue(tagKey.GetText())] = true
		}
		durationContext := dimension.DurationLit()
		if durationContext != nil {
			qs.interval = parseDuration(durationContext.(*parser.DurationLitContext))
		}
	}
}
