package plan

type Visitor interface {
	Visit(context any, n PlanNode) (r any)
}
