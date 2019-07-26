package stmt

// BinaryOP represents binary operation type
type BinaryOP int

const (
	AND BinaryOP = iota + 1
	OR

	ADD
	SUB
	MUL
	DIV
)
