package function

import "github.com/lindb/lindb/pkg/collections"

// FuncCall calls the function calc by function type and params
func FuncCall(funcType FuncType, params ...collections.FloatArray) collections.FloatArray {
	switch funcType {
	case Sum, Min, Max, Count:
		if len(params) == 0 {
			return nil
		}
		return params[0]
	case Avg:
		if len(params) < 2 {
			return nil
		}
		//FIXME need get params by pFieldID
		result := collections.NewFloatArray(params[0].Capacity())
		it := params[0].Iterator()
		for it.HasNext() {
			idx, value := it.Next()
			result.SetValue(idx, value/params[1].GetValue(idx))
		}
		return result
	default:
		return nil
	}
}
