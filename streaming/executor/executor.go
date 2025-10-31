package executor

type Executor interface {
	Process(event any)
	ResultSet(fn func(result any))
}
