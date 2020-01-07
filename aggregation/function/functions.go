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
		// params: 0=>sum, 1=>count
		if len(params) < 2 {
			return nil
		}
		result := collections.NewFloatArray(params[0].Capacity())
		it := params[0].Iterator()
		for it.HasNext() {
			idx, sum := it.Next()
			if params[1].HasValue(idx) {
				count := params[1].GetValue(idx)
				if count != 0 {
					result.SetValue(idx, sum/count)
				}
			}
		}
		return result
	default:
		return nil
	}
}
