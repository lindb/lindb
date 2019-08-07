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

	UNKNOWN
)

// BinaryOPString returns the binary operator's string value
func BinaryOPString(op BinaryOP) string {
	switch op {
	case AND:
		return "and"
	case OR:
		return "or"
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	default:
		return "unknown"
	}
}
