package optimization

import (
	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
)

type StreamDistribution string

var (
	Single   StreamDistribution = "Single"
	Multiple StreamDistribution = "Multiple"
	Fixed    StreamDistribution = "Fixed"
)

type ActualProps struct {
	global *ActualPropsGlobal
}

// isSingleNode checks if the plan will only execute on a single node.
func (props *ActualProps) isSingleNode() bool {
	// FIXME:  impl check single node
	return props.global.isSingleNode()
}

func (props *ActualProps) translate(translator func(symbol *plan.Symbol) *plan.Symbol) *ActualProps {
	return &ActualProps{
		global: props.global.translate(translator),
	}
}

type ActualPropsBuilder struct {
	global *ActualPropsGlobal
}

func NewActualPropsBuilder(global *ActualPropsGlobal) *ActualPropsBuilder {
	return &ActualPropsBuilder{
		global: global,
	}
}

func BuilderFrom(props *ActualProps) *ActualPropsBuilder {
	return &ActualPropsBuilder{
		global: props.global,
	}
}

func (b *ActualPropsBuilder) Build() *ActualProps {
	return &ActualProps{
		global: b.global,
	}
}

type PreferredProps struct {
	global *PreferredPropsGlobal
}

func Undistributed() *PreferredProps {
	return &PreferredProps{
		// global:Undistributed()
	}
}

func Any() *PreferredProps {
	return &PreferredProps{}
}

func Partitioned() *PreferredProps {
	return &PreferredProps{}
}

func PartitionedWithLocal(columns []*plan.Symbol) *PreferredProps {
	return &PreferredProps{}
}

type PlanProps struct {
	node  plan.PlanNode
	props *StreamProps
}

type StreamProps struct {
	distribution        StreamDistribution
	partitioningColumns []*plan.Symbol
	ordered             bool
}

func FixedStreams() *StreamProps {
	return &StreamProps{
		distribution: Fixed,
	}
}

func SingleStream() *StreamProps {
	return &StreamProps{
		distribution: Single,
	}
}

func (p *StreamProps) isPartitionedOn(columns []*plan.Symbol) bool {
	if len(p.partitioningColumns) == 0 {
		return false
	}
	for _, column := range p.partitioningColumns {
		if !lo.ContainsBy(columns, func(item *plan.Symbol) bool {
			return item.Name == column.Name
		}) {
			return false
		}
	}
	// columns contains all partitioning columns
	return true
}

func (p *StreamProps) translate(translator func(column *plan.Symbol) *plan.Symbol) *StreamProps {
	var newPartitioningColumns []*plan.Symbol
	for _, column := range p.partitioningColumns {
		translated := translator(column)
		if translated != nil {
			newPartitioningColumns = append(newPartitioningColumns, translated)
		}
	}
	return &StreamProps{
		distribution:        p.distribution,
		partitioningColumns: newPartitioningColumns,
	}
}

type StreamPreferredProps struct {
	distribution        StreamDistribution
	partitioningColumns []*plan.Symbol
	orderSensitive      bool
}

func (p *StreamPreferredProps) isSatisfiedBy(actualProps *StreamProps) bool {
	if p.distribution == "" && len(p.partitioningColumns) == 0 {
		// is there a specific preference
		return true
	}
	if p.orderSensitive && actualProps.ordered {
		// TODO: add check
		return true
	}
	if p.distribution != "" {
		switch {
		case p.distribution == Single && actualProps.distribution != Single:
			return false
		case p.distribution == Fixed && actualProps.distribution != Fixed:
			return false
		case p.distribution == Multiple && actualProps.distribution != Fixed && actualProps.distribution != Multiple:
			return false
		}
	} else if actualProps.distribution == Single {
		return true
	}
	if len(p.partitioningColumns) > 0 {
		return actualProps.isPartitionedOn(p.partitioningColumns)
	}
	return true
}

func (p *StreamPreferredProps) isSingleStreamPreferred() bool {
	return p.distribution != "" && p.distribution == Single
}

func (p *StreamPreferredProps) isParallelPreferred() bool {
	return p.distribution != "" && p.distribution != Single
}

func (p *StreamPreferredProps) withoutPreference() *StreamPreferredProps {
	return &StreamPreferredProps{
		orderSensitive: p.orderSensitive,
	}
}

func (p *StreamPreferredProps) withDefaultParallelism() *StreamPreferredProps {
	// FIXME:>>>>>>>>>>>>
	return p
}

func (p *StreamPreferredProps) withPartitioning(partitionSymbols []*plan.Symbol) *StreamPreferredProps {
	if len(partitionSymbols) == 0 {
		return singleStream()
	}
	desiredPartitioning := partitionSymbols
	if len(p.partitioningColumns) > 0 {
		// TODO: check exact column order?
		common := lo.Intersect(desiredPartitioning, p.partitioningColumns)
		if len(common) > 0 {
			desiredPartitioning = common
		}
	}
	return &StreamPreferredProps{
		distribution:        p.distribution,
		partitioningColumns: desiredPartitioning,
	}
}

func (p *StreamPreferredProps) withOrderSensitivity() *StreamPreferredProps {
	return &StreamPreferredProps{
		distribution:   p.distribution,
		orderSensitive: true,
	}
}

func (p *StreamPreferredProps) constrainTo(symbols []*plan.Symbol) *StreamPreferredProps {
	if len(p.partitioningColumns) == 0 {
		return p
	}
	// FIXME: add available symbols
	common := lo.Filter(p.partitioningColumns, func(column *plan.Symbol, index int) bool {
		return lo.ContainsBy(symbols, func(availableSymbol *plan.Symbol) bool {
			return availableSymbol.Name == column.Name
		})
	})
	if len(common) == 0 {
		return empty()
	}
	return &StreamPreferredProps{
		distribution:        p.distribution,
		partitioningColumns: common,
	}
}

func empty() *StreamPreferredProps {
	return &StreamPreferredProps{}
}

func defaultParallelism() *StreamPreferredProps {
	// FIXME:>>>>>>>>>>>>
	return empty()
}

func singleStream() *StreamPreferredProps {
	return &StreamPreferredProps{
		distribution: Single,
	}
}

type PreferredPropsGlobal struct {
	// nil => partitioned with some unknown scheme
	partitioningProps *plan.PartitioningProps
}

type ActualPropsGlobal struct {
	nodePartitioning *plan.Partitioning
}

func arbitraryPartition() *ActualPropsGlobal {
	return &ActualPropsGlobal{}
}

func singlePartition() *ActualPropsGlobal {
	// FIXME: impl single partition
	return partitionedOn(&plan.Partitioning{})
}

func partitionedOn(nodePartitioning *plan.Partitioning) *ActualPropsGlobal {
	return &ActualPropsGlobal{
		nodePartitioning: nodePartitioning,
	}
}

func (g *ActualPropsGlobal) translate(translator func(symbol *plan.Symbol) *plan.Symbol) *ActualPropsGlobal {
	return &ActualPropsGlobal{}
}

func (g *ActualPropsGlobal) isSingleNode() bool {
	if g.nodePartitioning == nil {
		return false
	}
	// TODO: fixme  check partition single node
	return g.nodePartitioning.Handle.IsSingleNode()
}
