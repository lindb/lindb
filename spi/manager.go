package spi

import (
	"fmt"
	"reflect"
)

var (
	splitSourceProviders = make(map[reflect.Type]SplitSourceProvider)
	pageSourceProviders  = make(map[reflect.Type]PageSourceProvider)
)

func RegisterSplitSourceProvider(table TableHandle, provider SplitSourceProvider) {
	splitSourceProviders[reflect.TypeOf(table)] = provider
}

func RegisterPageSourceProvider(table TableHandle, provider PageSourceProvider) {
	pageSourceProviders[reflect.TypeOf(table)] = provider
}

func GetSplitSourceProvider(table TableHandle) SplitSourceProvider {
	if prodiver, ok := splitSourceProviders[reflect.TypeOf(table)]; ok {
		return prodiver
	}
	panic(fmt.Sprintf("split source provider not found by table handle type for '%s'", reflect.TypeOf(table)))
}

func GetPageSourceProvider(table TableHandle) PageSourceProvider {
	if prodiver, ok := pageSourceProviders[reflect.TypeOf(table)]; ok {
		return prodiver
	}
	panic(fmt.Sprintf("page source provider not found by table handle type for '%s'", reflect.TypeOf(table)))
}
