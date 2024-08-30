package function

type OperatorType *operatorType

var (
	Equal    *operatorType = NewOperatorType("=", 2)
	Add      *operatorType = NewOperatorType("+", 2)
	Subtract *operatorType = NewOperatorType("-", 2)
	Multiply *operatorType = NewOperatorType("*", 2)
	Divide   *operatorType = NewOperatorType("/", 2)
)

type operatorType struct {
	operator      string
	argumentCount int
}

func NewOperatorType(operator string, argumentCount int) *operatorType {
	return &operatorType{
		operator:      operator,
		argumentCount: argumentCount,
	}
}
