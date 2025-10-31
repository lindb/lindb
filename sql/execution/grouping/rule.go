package grouping

type Rule interface {
	Map(value any, buf *Buffer)
	Unmap(buf *Buffer) any
}

type MapValuesRule struct {
	mapper *StringMapper
}

func newMapValuesRule(mapper *StringMapper) Rule {
	return &MapValuesRule{mapper: mapper}
}

func (m *MapValuesRule) Map(value any, buf *Buffer) {
	values, ok := value.(map[string]string)
	if !ok {
		return
	}
}

func (m *MapValuesRule) Unmap(buf *Buffer) any {
	values := make(map[string]string)

	// TODO:
	return values
}

type StringRule struct{}
