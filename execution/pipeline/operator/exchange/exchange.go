package exchange

import (
	"context"
	"fmt"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/execution/pipeline/operator"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/planner/plan"
)

type mergeKey struct {
	columns []string
}
type ExchangeOperatorFactory struct {
	sourceID   plan.PlanNodeID
	numOfChild int
}

func NewExchangeOperatorFactory(sourceID plan.PlanNodeID, numOfChild int) operator.SourceOperatorFactory {
	return &ExchangeOperatorFactory{
		sourceID:   sourceID,
		numOfChild: numOfChild,
	}
}

func (fct *ExchangeOperatorFactory) CreateOperator() operator.Operator {
	return NewExchangeOperator(fct.sourceID, fct.numOfChild)
}

type ExchangeOperator struct {
	sourceID            plan.PlanNodeID
	noMoreSplitsTracker *atomic.Int32

	splits chan *spi.BinarySplit

	ctx context.Context

	completed chan struct{}

	mergedPage *spi.Page
	mergedRows map[*mergeKey]int // merge key => row index
}

func NewExchangeOperator(sourceID plan.PlanNodeID, numOfChild int) operator.SourceOperator {
	return &ExchangeOperator{
		sourceID:            sourceID,
		noMoreSplitsTracker: atomic.NewInt32(int32(numOfChild)),
		splits:              make(chan *spi.BinarySplit, 10),
		mergedPage:          spi.NewPage(),
		mergedRows:          make(map[*mergeKey]int),
		completed:           make(chan struct{}, 1),
	}
}

func (op *ExchangeOperator) GetSourceID() plan.PlanNodeID {
	return op.sourceID
}

func (op *ExchangeOperator) NoMoreSplits() {
	newVal := op.noMoreSplitsTracker.Dec()
	if newVal == 0 {
		close(op.splits)
	}
}

// AddSplit implements operator.SourceOperator
func (op *ExchangeOperator) AddSplit(split spi.Split) {
	if data, ok := split.(*spi.BinarySplit); ok {
		op.splits <- data
	}
}

// AddInput implements Operator
func (op *ExchangeOperator) AddInput(page *spi.Page) {
	panic(fmt.Errorf("%w: exchange cannot take input", constants.ErrNotSupportOperation))
}

// GetOutput implements Operator
func (op *ExchangeOperator) GetOutput() *spi.Page {
	for split := range op.splits {
		page := split.Page
		if page == nil {
			continue
		}
		// it := page.Iterator()
		// groupingColumns := page.Grouping
		// for row := it.Begin(); row != it.End(); row = it.Next() {
		// 	fmt.Println("kkkk......")
		// 	op.mergedPage.AppendColumn(page.Layout[], page.Columns[])
		// }
		op.mergedPage.Layout = page.Layout
		op.mergedPage.Columns = page.Columns
	}
	return op.mergedPage
}

func (op *ExchangeOperator) Finish() {
	close(op.completed)
}

func (op *ExchangeOperator) IsFinished() bool {
	// <-op.completed
	return true
}
