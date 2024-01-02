package function

type OperatorType *operatorType

var (
	Equal *operatorType = NewOperatorType("=", 2)
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
