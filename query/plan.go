package query

// Plan represents an execute plan of a query language with computing and storage
type Plan interface {
	// Plan plans the query language, then generate an execute plan
	Plan() error
}
