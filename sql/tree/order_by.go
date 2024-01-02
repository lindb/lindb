package tree

type Ordering string

var (
	ASCENDING  Ordering = "ASCENDING"
	DESCENDING Ordering = "DESCENDING"
)

type OrderBy struct {
	BaseNode
	SortItems []*SortItem
}

type SortItem struct {
	BaseNode
	SortKey  Expression
	Ordering Ordering
}
