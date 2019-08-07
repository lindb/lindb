package function

import "github.com/lindb/lindb/pkg/collections"

// FuncCall calls the function calc by function type and params
func FuncCall(funcType FuncType, params ...collections.FloatArray) collections.FloatArray {
	switch funcType {
	case Sum, Min, Max:
		if len(params) == 0 {
			return nil
		}
		return params[0]
	default:
		return nil
	}
}
