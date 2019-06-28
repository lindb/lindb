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

type queryStatement struct {
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
	fieldId             int
}

// NewDefaultQueryStatement build lindb default query statement
func NewDefaultQueryStatement() *queryStatement {
	return &queryStatement{
		limit:               20,
		fieldId:             1,
		condition:           new(proto.Condition),
		fieldExprList:       make([]proto.FieldExpr, 0, 10),
		groupBy:             make(map[string]bool),
		conditionAggregates: make(map[*proto.Condition]map[util.AggregatorUnit]Unit),
	}
}

// Parse lindb parse query sql
func (qs *queryStatement) Parse(ctx *parser.Query_stmtContext) {
	// parse measurement like 'table'
	if ctx.From_clause() != nil {
		clauseContext := ctx.From_clause().(*parser.From_clauseContext)
		qs.measurement = util.GetStringValue(clauseContext.Metric_name().GetText())
	}
	// parse select fields
	if ctx.Fields() != nil {
		fieldsContext := ctx.Fields().(*parser.FieldsContext)
		qs.parseFields(fieldsContext)
	}
	//parse where clause
	if ctx.Where_clause() != nil {
		whereClause := ctx.Where_clause().(*parser.Where_clauseContext)
		qs.parseClause(whereClause)
	}
	//parse group by
	if ctx.Group_by_clause() != nil {
		groupBy := ctx.Group_by_clause().(*parser.Group_by_clauseContext)
		qs.parseGroupBy(groupBy)
	}
	//parse order by
	if ctx.Order_by_clause() != nil {
		orderBy := ctx.Order_by_clause().(*parser.Order_by_clauseContext)
		qs.parseOrderBy(orderBy)
	}
	//parse limit
	if ctx.Limit_clause() != nil {
		limitClause := ctx.Limit_clause().(*parser.Limit_clauseContext)
		qs.parseLimit(limitClause)
	}
}

// parseLimit parse lindb sql limit
func (qs *queryStatement) parseLimit(limitClause *parser.Limit_clauseContext) {
	if nil != limitClause {
		limit, _ := strconv.ParseInt(limitClause.L_INT().GetText(), 10, 32)
		qs.limit = int32(limit)
	}
}

// build build lindb query statement
func (qs *queryStatement) build() *proto.Query {
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
func (qs *queryStatement) parseClause(whereClauseContext *parser.Where_clauseContext) {
	if nil == whereClauseContext {
		return
	}
	clause := whereClauseContext.Clause_boolean_expr().(*parser.Clause_boolean_exprContext)
	allClause := clause.AllClause_boolean_expr()
	if nil != allClause && len(allClause) > 0 {
		var exprs = make([]parser.Clause_boolean_exprContext, 0, len(allClause))
		for i := range allClause {
			expr := allClause[i].(*parser.Clause_boolean_exprContext)
			exprs = append(exprs, *expr)
		}
		qs.parseConditionClause(exprs)
	} else if nil != clause.Time_expr() {
		expr := clause.Time_expr().(*parser.Time_exprContext)
		qs.parseTimeExpr(expr)
	} else if nil != clause.Tag_boolean_expr() {
		expr := clause.Tag_boolean_expr().(*parser.Tag_boolean_exprContext)
		parseTagExpr(expr, qs.condition)
	}

	if qs.condition.Operator == proto.LogicOperator_UNKNOWN {
		qs.condition.Operator = proto.LogicOperator_AND
	}
}

// parseConditionClause parse lindb sql condition clause
func (qs *queryStatement) parseConditionClause(clauseList []parser.Clause_boolean_exprContext) {
	if clauseList == nil {
		return
	}
	if 2 == len(clauseList) {
		for i := 0; i < 2; i++ {
			clause := clauseList[i]
			if nil != clause.Time_expr() {
				timeExpr := clause.Time_expr().(*parser.Time_exprContext)
				qs.parseTimeExpr(timeExpr)
			} else if nil != clause.Tag_boolean_expr() {
				exprContext := clause.Tag_boolean_expr().(*parser.Tag_boolean_exprContext)
				parseTagExpr(exprContext, qs.condition)
			}
		}
	}
}

// parseTimeExpr lindb Parse time range condition
func (qs *queryStatement) parseTimeExpr(ctx *parser.Time_exprContext) {
	if ctx == nil {
		return
	}
	timeBooleanExprContexts := ctx.AllTime_boolean_expr()
	for i := range timeBooleanExprContexts {
		timeExpr := timeBooleanExprContexts[i].(*parser.Time_boolean_exprContext)
		timestamp := timeExpr.Ident()
		var times int64 = 0
		if timestamp != nil {
			times = util.ParseTimestamp(util.GetStringValue(timestamp.GetText()))
		}
		iNowExpr := timeExpr.Now_expr()
		if iNowExpr != nil {
			now := time.Now().Unix() * 1000
			nowExpr := iNowExpr.(*parser.Now_exprContext)
			if nowExpr.Duration_lit() != nil {
				durationList := nowExpr.Duration_lit().(*parser.Duration_litContext)
				duration := parseDuration(durationList)
				times = now + duration
			} else {
				times = now
			}
		}
		iOpCtx := timeExpr.Bool_expr_binary_operator()
		opCtx := iOpCtx.(*parser.Bool_expr_binary_operatorContext)
		if opCtx.T_GREATER() != nil || opCtx.T_GREATEREQUAL() != nil {
			qs.startTime = times
		}
		if opCtx.T_LESS() != nil || opCtx.T_LESSEQUAL() != nil {
			qs.endTime = times
		}

	}
}

// parseFields lindb sql parse a list of one or more fields
func (qs *queryStatement) parseFields(ctx *parser.FieldsContext) {
	if ctx == nil {
		return
	}
	fieldContexts := ctx.AllField()
	if fieldContexts == nil || len(fieldContexts) == 0 {
		return
	}
	for i := range fieldContexts {
		fieldContext := fieldContexts[i].(*parser.FieldContext)
		qs.parseField(fieldContext)
	}
}

// parseOrderBy lindb sql parse order by
func (qs *queryStatement) parseOrderBy(orderByClause *parser.Order_by_clauseContext) {
	if orderByClause != nil {
		iSortFieldsContext := orderByClause.Sort_fields()
		if iSortFieldsContext != nil {
			sortFieldsContext := iSortFieldsContext.(*parser.Sort_fieldsContext)
			sortFieldContextList := sortFieldsContext.AllSort_field()
			if sortFieldContextList != nil && len(sortFieldContextList) > 0 {
				sortFieldContext := sortFieldContextList[0].(*parser.Sort_fieldContext)
				if len(sortFieldContext.AllT_ASC()) > 0 {
					qs.desc = false
				} else {
					qs.desc = true
				}

				expr := new(proto.Expr)
				exprContext := sortFieldContext.Expr().(*parser.ExprContext)
				qs.parseExpr(exprContext, expr)
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
func (qs *queryStatement) parseField(ctx *parser.FieldContext) {
	exprCtx := ctx.Expr().(*parser.ExprContext)
	field := new(proto.FieldExpr)
	expr := new(proto.Expr)
	qs.parseExpr(exprCtx, expr)
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

func (qs *queryStatement) parseExpr(ctx *parser.ExprContext, expr *proto.Expr) {
	if nil == ctx {
		return
	}
	if nil != ctx.Expr_atom() {
		atom := ctx.Expr_atom().(*parser.Expr_atomContext)
		qs.parseExprAtom(atom, expr)
	} else if nil != ctx.Expr_func() {
		fun := ctx.Expr_func().(*parser.Expr_funcContext)
		qs.parseCall(fun, expr)
	} else if nil != ctx.Duration_lit() {
		lin := ctx.Duration_lit().(*parser.Duration_litContext)
		interval := parseDuration(lin)
		longVal := new(proto.LongExpr)
		longVal.Value = interval
		//todo set long value
	} else if nil != ctx.AllExpr() {
		exprContexts := ctx.AllExpr()
		size := len(exprContexts)
		if 2 == size {
			binaryExp := new(proto.BinaryExpr)
			left := new(proto.Expr)
			exprContext0 := exprContexts[0].(*parser.ExprContext)
			qs.parseExpr(exprContext0, left)
			binaryExp.Left = left
			var op proto.Operator
			if ctx.T_ADD() != nil {
				op = proto.Operator_ADD
			} else if ctx.T_SUB() != nil {
				op = proto.Operator_SUB
			} else if ctx.T_MUL() != nil {
				op = proto.Operator_MUL
			} else if ctx.T_DIV() != nil {
				op = proto.Operator_DIV
			} else {
				panic(fmt.Sprintln("Unknown operator type in expr"))
			}
			binaryExp.Op = op
			right := new(proto.Expr)
			exprContext1 := exprContexts[1].(*parser.ExprContext)
			qs.parseExpr(exprContext1, right)
			binaryExp.Right = right
			//todo set binary value
		} else if 1 == size {
			exprContext := exprContexts[0].(*parser.ExprContext)
			qs.parseExpr(exprContext, expr)
		} else {
			panic(fmt.Sprintln("Unknown expr type"))
		}
	}

}

func (qs *queryStatement) parseExprAtom(ctx *parser.Expr_atomContext, expr *proto.Expr) {
	if ctx.Ident() != nil {
		valRef := new(proto.ValRefExpr)
		field := util.GetStringValue(ctx.Ident().GetText())
		filterContext, e := ctx.Ident_filter().(*parser.Ident_filterContext)
		var currentFieldId int32
		if true == e {
			condition := parseFilter(filterContext)
			currentFieldId = qs.addField(condition, *util.NewAggregatorUnit(field, nil, nil))
		} else {
			currentFieldId = qs.addField(new(proto.Condition), *util.NewAggregatorUnit(field, nil, nil))
		}
		valRef.RefName = strconv.Itoa(int(currentFieldId))
		//todo set ref value
	} else if ctx.Int_number() != nil {
		longVal := new(proto.LongExpr)
		dec, _ := strconv.ParseInt(ctx.Int_number().GetText(), 10, 64)
		longVal.Value = dec
		//todo set long value
	} else if ctx.Dec_number() != nil {
		doubleVal := new(proto.DoubleExpr)
		dec, _ := strconv.ParseFloat(ctx.Dec_number().GetText(), 64)
		doubleVal.Value = dec
		//todo set double value
	} else {
		panic(fmt.Sprintf("Unknown expr atom type"))
	}
}

// addField parse lindb sql field
func (qs *queryStatement) addField(condition *proto.Condition, unit util.AggregatorUnit) int32 {
	aggregatorUnits := qs.conditionAggregates[condition]
	if aggregatorUnits == nil {
		aggregatorUnits = make(map[util.AggregatorUnit]Unit)
	}
	id := aggregatorUnits[unit]
	if id.unit <= 0 {
		id.unit = int32(qs.fieldId)
		aggregatorUnits[unit] = id
		qs.fieldId++
	}
	qs.conditionAggregates[condition] = aggregatorUnits
	return id.unit
}

// parseCall lindb sql parse a function call
func (qs *queryStatement) parseCall(ctx *parser.Expr_funcContext, expr *proto.Expr) {
	name := util.GetStringValue(ctx.Ident().GetText())
	if !util.IsDownSamplingOrAggregator(name) {
		callExpr := new(proto.CallExpr)
		callExpr.Name = name
		if ctx.Expr_func_params() != nil {
			funcParams := ctx.Expr_func_params().(*parser.Expr_func_paramsContext)
			funcParam := funcParams.AllFunc_param()
			for i := range funcParam {
				param := funcParam[i].(*parser.Func_paramContext)
				callParam := new(proto.Expr)
				qs.parseCallParam(param, callParam)
				callExpr.Args = append(callExpr.Args, callParam)
			}
		}
		//todo expr call
		return
	}
	varRefBuilder := new(proto.ValRefExpr)
	currentFieldId := qs.parseDownSamplingOrAggregator(ctx)
	varRefBuilder.RefName = strconv.Itoa(int(currentFieldId))
	//todo set ref value
}

// parseDownSamplingOrAggregator parse lindb down sampling or aggregator
func (qs *queryStatement) parseDownSamplingOrAggregator(ctx *parser.Expr_funcContext) int32 {
	var values = make(map[int]string)
	filterCondition := getCallNames(ctx, values)
	size := len(values)
	var currentFieldId int32 = 0
	var aggregatorUnit util.AggregatorUnit
	if 2 == size {
		fun := values[0]
		function := util.ValueOf(fun)
		if strings.HasPrefix(fun, util.DownSampling) {
			aggregatorUnit = *util.NewAggregatorUnit(values[1], &function, nil)
		} else {
			aggregatorUnit = *util.NewAggregatorUnit(values[1], nil, &function)
		}
	} else if 3 == size {
		fieldName := values[2]
		fun1 := util.ValueOf(values[1])
		fun0 := util.ValueOf(values[0])
		aggregatorUnit = *util.NewAggregatorUnit(fieldName, &fun1, &fun0)
	} else if 1 == size {
		aggregatorUnit = *util.NewAggregatorUnit(values[0], nil, nil)
	} else {
		panic(fmt.Sprintln("sql is not valid"))
	}
	if filterCondition == nil {
		currentFieldId = qs.addField(new(proto.Condition), aggregatorUnit)
	} else {
		currentFieldId = qs.addField(filterCondition, aggregatorUnit)
	}
	return currentFieldId
}

func getCallNames(ctx *parser.Expr_funcContext, values map[int]string) *proto.Condition {
	var filterCondition *proto.Condition
	callName := util.GetStringValue(ctx.Ident().GetText())
	valuesLen := len(values)
	values[valuesLen] = callName
	if ctx.Expr_func_params() != nil {
		paramsContext := ctx.Expr_func_params().(*parser.Expr_func_paramsContext)
		funcParams := paramsContext.AllFunc_param()
		for i := range funcParams {
			param := funcParams[i].(*parser.Func_paramContext)
			if param.Expr() == nil {
				continue
			}
			exprContext := param.Expr().(*parser.ExprContext)
			if exprContext.Expr_func() != nil {
				exprFuncContext := exprContext.Expr_func().(*parser.Expr_funcContext)
				getCallNames(exprFuncContext, values)
			} else if exprContext.Expr_atom() != nil {
				atomContext := exprContext.Expr_atom().(*parser.Expr_atomContext)
				identContext := atomContext.Ident()
				if identContext != nil {
					valuesLen = len(values)
					values[valuesLen] = util.GetStringValue(identContext.GetText())
				}
				if atomContext.Ident_filter() != nil {
					filterContext := atomContext.Ident_filter().(*parser.Ident_filterContext)
					filterCondition = parseFilter(filterContext)
				}
			}
		}
	}
	return filterCondition
}

// parseFilter parse lindb sql filter
func parseFilter(identFilterContext *parser.Ident_filterContext) *proto.Condition {
	if identFilterContext == nil {
		return nil
	}
	fieldCondition := new(proto.Condition)
	expr := identFilterContext.Tag_boolean_expr().(*parser.Tag_boolean_exprContext)
	parseTagExpr(expr, fieldCondition)
	if proto.LogicOperator_UNKNOWN == fieldCondition.Operator {
		fieldCondition.Operator = proto.LogicOperator_AND
	}
	return fieldCondition
}

// parseCallParam parse lindb sql call expr param
func (qs *queryStatement) parseCallParam(ctx *parser.Func_paramContext, callParam *proto.Expr) {
	if nil != ctx.Tag_boolean_expr() {
		conditionBuilder := new(proto.Condition)
		exprContext := ctx.Tag_boolean_expr().(*parser.Tag_boolean_exprContext)
		parseTagExpr(exprContext, conditionBuilder)
	} else {
		exprContext := ctx.Expr().(*parser.ExprContext)
		qs.parseExpr(exprContext, callParam)
	}

}

// parseTagExpr parse lindb sql tag
func parseTagExpr(ctx *parser.Tag_boolean_exprContext, condition *proto.Condition) {
	if nil == ctx {
		return
	}
	tagKeyCtx := ctx.Tag_key()
	if nil != tagKeyCtx {
		filter := parseTag(ctx)
		condition.TagFilters = append(condition.TagFilters, filter)
	} else {
		tagCtxList := ctx.AllTag_boolean_expr()
		size := len(tagCtxList)
		tagCtx0 := tagCtxList[0].(*parser.Tag_boolean_exprContext)
		if 2 == size {
			var logicOp proto.LogicOperator
			if nil != ctx.T_AND() {
				logicOp = proto.LogicOperator_AND
			} else if nil != ctx.T_OR() {
				logicOp = proto.LogicOperator_OR
			}

			tagCtx1 := tagCtxList[1].(*parser.Tag_boolean_exprContext)
			if proto.LogicOperator_OR == logicOp {
				orCondition := new(proto.Condition)
				orCondition.Operator = logicOp
				parseTagExpr(tagCtx0, orCondition)
				parseTagExpr(tagCtx1, orCondition)
				condition.Condition = append(condition.Condition, orCondition)
			} else {
				if nil != tagCtx0.Tag_key() {
					filter := parseTag(tagCtx0)
					condition.TagFilters = append(condition.TagFilters, filter)
				} else {
					parseTagExpr(tagCtx0, condition)
				}

				if nil != tagCtx1.Tag_key() {
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
			if proto.LogicOperator_AND == condition2.Operator && 0 == len(condition2.Condition) && proto.LogicOperator_AND == condition.Operator {
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
func parseTag(ctx *parser.Tag_boolean_exprContext) *proto.TagFilter {
	tagKeyCtx := ctx.Tag_key()
	var op proto.Operator
	if nil != ctx.T_EQUAL() {
		op = proto.Operator_EQUAL
	} else if nil != ctx.T_LIKE() || nil != ctx.T_REGEXP() {
		op = proto.Operator_LIKE
	} else if nil != ctx.T_NOTEQUAL() || nil != ctx.T_NOTEQUAL2() {
		op = proto.Operator_NOT_EQUAL
	} else if nil != ctx.T_NOT() && nil != ctx.T_IN() {
		op = proto.Operator_NOT_IN
	} else if nil != ctx.T_NOT() && nil != ctx.T_IN() {
		op = proto.Operator_IN
	} else {
		panic(fmt.Sprintf("unknown filter operator"))
	}

	filter := new(proto.TagFilter)
	tagKey := util.GetStringValue(tagKeyCtx.GetText())
	filter.TagKey = tagKey
	filter.Op = op
	ctx.Tag_value()
	if proto.Operator_IN == op || proto.Operator_NOT_IN == op {
		tagValueListContext, b := ctx.Tag_value_list().(*parser.Tag_value_listContext)
		if true == b {
			tagValueContexts := tagValueListContext.AllTag_value()
			if len(tagValueContexts) > 0 {
				var tagValues = make([]string, 0)
				for i := range tagValueContexts {
					item := tagValueContexts[i].(*parser.Tag_valueContext)
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
		tagValueContext := ctx.Tag_value().(*parser.Tag_valueContext)
		tagValue := util.GetStringValue(tagValueContext.Ident().GetText())
		filter.TagValue = tagValue
	}
	return filter
}

// parseDuration lindb sql parse a  time duration from a string
func parseDuration(ctx *parser.Duration_litContext) int64 {
	duration, _ := strconv.ParseInt(ctx.Int_number().GetText(), 10, 64)
	unit, b := ctx.Interval_item().(*parser.Interval_itemContext)
	var result int64 = 0
	if false == b {
		return result
	}
	if nil != unit.T_SECOND() {
		result = duration * util.OneSeconds
	} else if nil != unit.T_MINUTE() {
		result = duration * util.OneMinute
	} else if nil != unit.T_HOUR() {
		result = duration * util.OneHour
	} else if nil != unit.T_DAY() {
		result = duration * util.OneDay
	} else if nil != unit.T_WEEK() {
		result = duration * util.OneWeek
	} else if nil != unit.T_MONTH() {
		result = duration * util.OneMonth
	} else if nil != unit.T_YEAR() {
		result = duration * util.OneYear
	}
	return result
}

// parseGroupBy parse lindb sql group by
func (qs *queryStatement) parseGroupBy(ctx *parser.Group_by_clauseContext) {
	if nil == ctx {
		return
	}
	dimensions, err := ctx.Dimensions().(*parser.DimensionsContext)
	if false == err {
		return
	}
	var dimensionContextList = dimensions.AllDimension()
	for i := range dimensionContextList {
		dimension := dimensionContextList[i].(*parser.DimensionContext)
		tagKey := dimension.Ident()
		if nil != tagKey {
			qs.groupBy[util.GetStringValue(tagKey.GetText())] = true
		}
		durationContext := dimension.Duration_lit()
		if durationContext != nil {
			qs.interval = parseDuration(durationContext.(*parser.Duration_litContext))
		}
	}
}
