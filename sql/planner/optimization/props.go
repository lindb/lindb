package optimization

import (
	"github.com/lindb/lindb/sql/planner/plan"
)

type StreamDistribution string

var (
	Single   StreamDistribution = "Single"
	Multiple StreamDistribution = "Multiple"
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
	distribution StreamDistribution
	ordered      bool
}

type StreamPreferredProps struct {
	distribution   StreamDistribution
	orderSensitive bool
}

func (p *StreamPreferredProps) isSatisfiedBy(actualProps *StreamProps) bool {
	if p.distribution == "" {
		// is there a specific preference
		return true
	}
	if p.orderSensitive && actualProps.ordered {
		return true
	}
	return false
}

func (p *StreamPreferredProps) isSingleStreamPreferred() bool {
	return p.distribution != "" && p.distribution == Single
}

func (p *StreamPreferredProps) isParallelPreferred() bool {
	return p.distribution != "" && p.distribution == Single
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

func (p *StreamPreferredProps) withOrderSensitivity() *StreamPreferredProps {
	return &StreamPreferredProps{
		distribution:   p.distribution,
		orderSensitive: true,
	}
}

func (p *StreamPreferredProps) constrainTo(symbols []*plan.Symbol) *StreamPreferredProps {
	// FIXME:>>>>>>>>>>>>
	return p
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
