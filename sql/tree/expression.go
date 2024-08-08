package tree

type (
	LogicalOperator    string
	ComparisonOperator string
)

var (
	LogicalAND LogicalOperator = "AND"
	LogicalOR  LogicalOperator = "OR"

	ComparisonEqual ComparisonOperator = "="
)

type Expression interface {
	Node
}

type Cast struct {
	BaseNode

	Expression Expression `json:"expression"`
}

func (n *Cast) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type FieldReference struct {
	BaseNode

	FieldIndex int
}

func (n *FieldReference) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type DereferenceExpression struct {
	BaseNode
	Base  Expression
	Field *Identifier
}

func (n *DereferenceExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

func (n *DereferenceExpression) ToQualifiedName() (name *QualifiedName) {
	if n.Field == nil {
		return
	}
	if base, ok := n.Base.(*Identifier); ok {
		name = NewQualifiedName([]*Identifier{base, n.Field})
	} else if base, ok := n.Base.(*DereferenceExpression); ok {
		baseQualifiedName := base.ToQualifiedName()
		if baseQualifiedName != nil {
			parts := baseQualifiedName.OriginalParts
			parts = append(parts, n.Field)
			name = NewQualifiedName(parts)
		}
	}
	return
}

type ArithmeticBinaryExpression struct {
	Left     Expression
	Right    Expression
	Operator string
}

type ComparisonExpression struct {
	BaseNode

	Left     Expression         `json:"left"`
	Right    Expression         `json:"right"`
	Operator ComparisonOperator `json:"operator"`
}

func (n *ComparisonExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type LogicalExpression struct {
	BaseNode
	Operator LogicalOperator
	Terms    []Expression
}

func (n *LogicalExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type InListExpression struct {
	BaseNode
	Values []Expression
}

func (n *InListExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}

type NotExpression struct {
	BaseNode
	Value Expression
}

func (n *NotExpression) Accept(context any, visitor Visitor) (r any) {
	return visitor.Visit(context, n)
}
