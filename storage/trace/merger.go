package trace

type TraceMergeOperator struct{}

func (op *TraceMergeOperator) Name() string {
	return "TraceMergeOperator"
}

func (op *TraceMergeOperator) FullMerge(key, existingValue []byte, operands [][]byte) ([]byte, bool) {
	total := len(existingValue)
	for _, v := range operands {
		total += len(v)
	}
	dest := make([]byte, total)
	offset := copy(dest, existingValue)
	for _, operand := range operands {
		offset += copy(dest[offset:], operand)
	}
	return dest, true
}

func (op *TraceMergeOperator) PartialMerge(key, leftOperand, rightOperand []byte) ([]byte, bool) {
	dest := make([]byte, (len(rightOperand) + len(leftOperand)))
	copy(dest, leftOperand)
	copy(dest[len(leftOperand):], rightOperand)
	return dest, true
}
