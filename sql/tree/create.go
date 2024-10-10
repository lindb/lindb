package tree

import "github.com/lindb/lindb/models"

type CreateOption interface{}

type EngineOption struct {
	Type models.EngineType
}

type CreateDatabase struct {
	BaseNode
	Name          string
	CreateOptions []CreateOption
	Props         map[string]any
	Rollup        []RollupOption
}

func (n *CreateDatabase) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type CreateBroker struct {
	BaseNode
	Options map[string]any
	Name    string
}

func (n *CreateBroker) Accept(context any, visitor Visitor) any {
	return visitor.Visit(context, n)
}

type RollupOption struct {
	Options map[string]any
}
