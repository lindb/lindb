package types

var (
	Equal    OperatorType = NewOperatorType("=", 2)
	Add      OperatorType = NewOperatorType("+", 2)
	Subtract OperatorType = NewOperatorType("-", 2)
	Multiply OperatorType = NewOperatorType("*", 2)
	Divide   OperatorType = NewOperatorType("/", 2)
	Modulus  OperatorType = NewOperatorType("%", 2)
)

type OperatorType struct {
	Operator      string
	ArgumentCount int
}

func NewOperatorType(operator string, argumentCount int) OperatorType {
	return OperatorType{
		Operator:      operator,
		ArgumentCount: argumentCount,
	}
}
