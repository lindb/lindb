package optimization

import "github.com/lindb/lindb/sql/planner/plan"

type StreamDistribution string

var (
	Single   StreamDistribution = "Single"
	Multiple StreamDistribution = "Multiple"
)

type ActualProps struct {
}

type PreferredProps struct {
}

func undistributed() *PreferredProps {
	return &PreferredProps{}
}

func partitioned() *PreferredProps {
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
	//FIXME:>>>>>>>>>>>>
	return p
}

func (p *StreamPreferredProps) withOrderSensitivity() *StreamPreferredProps {
	return &StreamPreferredProps{
		distribution:   p.distribution,
		orderSensitive: true,
	}
}
func (p *StreamPreferredProps) constrainTo(symbols []*plan.Symbol) *StreamPreferredProps {
	//FIXME:>>>>>>>>>>>>
	return p
}

func empty() *StreamPreferredProps {
	return &StreamPreferredProps{}
}

func defaultParallelism() *StreamPreferredProps {
	//FIXME:>>>>>>>>>>>>
	return empty()
}
func singleStream() *StreamPreferredProps {
	return &StreamPreferredProps{
		distribution: Single,
	}
}
