package strutil

import (
	"reflect"
	"strings"
	"unsafe"
)

// GetStringValue aggregation format function name
func GetStringValue(rawString string) string {
	if len(rawString) > 0 {
		if (strings.HasPrefix(rawString, "'") && strings.HasSuffix(rawString, "'")) ||
			(strings.HasPrefix(rawString, "\"") && strings.HasSuffix(rawString, "\"")) {
			return rawString[1 : len(rawString)-1]
		}
		return rawString
	}
	return ""
}

func ByteSlice2String(bytes []byte) string {
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: sliceHeader.Data,
		Len:  sliceHeader.Len,
	}))
}
