package spi

// Iterator4Page represents a iterator that is used to iterator rows inside a page.
type Iterator4Page struct {
	page *Page

	cursor  int
	numRows int
}

// NewIterator4Page creates a iterator for Page.
func NewIterator4Page(page *Page) *Iterator4Page {
	return &Iterator4Page{
		page: page,
	}
}

// Gegin resets the cursor of the iterator and returns the first Row.
func (it *Iterator4Page) Begin() Row {
	if it.page == nil {
		return it.End()
	}

	it.numRows = it.page.NumRows()
	if it.numRows == 0 {
		return it.End()
	}
	it.cursor = 1
	return it.page.GetRow(0)
}

// Next returns the next Row.
func (it *Iterator4Page) Next() Row {
	if it.cursor >= it.numRows {
		it.cursor = it.numRows + 1
		return it.End()
	}
	row := it.page.GetRow(it.cursor)
	it.cursor++
	return row
}

// End returns the invalid end Row.
func (it *Iterator4Page) End() Row {
	return Row{}
}
