package spi

type Page struct {
	Layout   []ColumnMetadata `json:"layout,omitempty"`
	Grouping []int            `json:"grouping,omitempty"` // grouping column indexes
	Columns  []*Column        `json:"columns,omitempty"`
}

func NewPage() *Page {
	return &Page{}
}

func (p *Page) SetGrouping(columnIndexes []int) {
	p.Grouping = columnIndexes
}

func (p *Page) AppendColumn(info ColumnMetadata, column *Column) {
	p.Layout = append(p.Layout, info)
	p.Columns = append(p.Columns, column)
}

// GetRow gets the Row in the page with the row index.
func (p *Page) GetRow(idx int) Row {
	// TODO: select rows?
	return Row{p: p, idx: idx}
}

// NumRows returns the number of rows in the page.
func (p *Page) NumRows() int {
	if len(p.Columns) == 0 {
		return 0
	}
	// TODO: select rows/no column
	return p.Columns[0].NumOfRows
}

func (p *Page) Iterator() *Iterator4Page {
	return NewIterator4Page(p)
}
