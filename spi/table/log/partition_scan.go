package log

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/storage/log"
	"github.com/lindb/lindb/storage/store"
)

type Partition struct {
	tableScan  *TableScan
	shard      store.Shard
	paritition store.Partition
	segments   []store.Segment
}

type RowsLookupVisitor struct {
	tableScan *TableScan
	segment   *log.Segment

	logIDs *roaring.Bitmap
}

func NewRowLookupVisitor(tableScan *TableScan, segment *log.Segment, logIDs *roaring.Bitmap) *RowsLookupVisitor {
	return &RowsLookupVisitor{
		tableScan: tableScan,
		segment:   segment,
		logIDs:    logIDs,
	}
}

func (v *RowsLookupVisitor) Visit(context any, n tree.Node) any {
	var logIDs *roaring.Bitmap
	var fieldID uint32

	switch node := n.(type) {
	case *tree.ComparisonExpression:
		fieldID, logIDs = v.visitPredicate(node)
		if node.Operator == tree.ComparisonNEQ {
			// get all series ids for tag key
			all := v.segment.FindLogIDsByFields([]uint32{fieldID})
			// do and not got series ids not in 'a' list
			all.AndNot(logIDs)
			return all
		}
	case *tree.InPredicate, *tree.RegexPredicate, *tree.LikePredicate:
		_, logIDs = v.visitPredicate(node)
	case *tree.NullPredicate:
		fieldID, _ = v.visitPredicate(node)
		all := v.segment.FindLogIDsByFields([]uint32{fieldID})
		if node.Not {
			return all
		}
		v.logIDs.AndNot(all)
		return v.logIDs
	case *tree.NotExpression:
		// get filter series ids
		fieldID, logIDs = v.visitPredicate(node.Value)
		// TODO: cache if dup
		// get all series ids for tag key
		all := v.segment.FindLogIDsByFields([]uint32{fieldID})
		// do and not got series ids not in 'a' list
		all.AndNot(logIDs)
		return all
	case *tree.LogicalExpression:
		for _, term := range node.Terms {
			matchResult := term.Accept(context, v).(*roaring.Bitmap)
			if logIDs == nil {
				logIDs = matchResult
			} else {
				if node.Operator == tree.LogicalAND {
					logIDs.And(matchResult)
				} else {
					logIDs.Or(matchResult)
				}
			}
		}
		return logIDs
	case *tree.Cast:
		return node.Expression.Accept(context, v)
	}
	return logIDs
}

func (v *RowsLookupVisitor) visitPredicate(node tree.Node) (uint32, *roaring.Bitmap) {
	field, ok := v.tableScan.filterResult[node.GetID()]
	if !ok {
		panic(constants.ErrSeriesIDNotFound)
	}
	if _, ok := node.(*tree.NullPredicate); ok {
		return field.fieldID, nil
	}
	fmt.Printf("node=%d field=%v,tag value ids=%v\n", node.GetID(), field.fieldID, field.fieldValueIDs)
	logIDs := v.segment.FindLogIDsByFields(field.fieldValueIDs)
	return field.fieldID, logIDs
}
