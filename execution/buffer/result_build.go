package buffer

import (
	"github.com/lindb/lindb/execution/model"
	"github.com/lindb/lindb/spi"
)

type ResultSetBuild struct {
	pages     chan *spi.Page
	resultSet *model.ResultSet
}

func CreateResultSetBuild() *ResultSetBuild {
	return &ResultSetBuild{
		pages:     make(chan *spi.Page),
		resultSet: model.NewResultSet(),
	}
}

func (rsb *ResultSetBuild) AddPage(page *spi.Page) {
	rsb.pages <- page
}

func (rsb *ResultSetBuild) Process() {
	for page := range rsb.pages {
		rsb.resultSet.Schema.Columns = page.Layout
		it := page.Iterator()
		for row := it.Begin(); row != it.End(); row = it.Next() {
			columns := []any{}
			for i := range page.Layout {
				columns = append(columns, row.GetTimeSeries(i))
			}
			rsb.resultSet.Rows = append(rsb.resultSet.Rows, columns)
		}
	}
}

func (rsb *ResultSetBuild) Complete() {
	close(rsb.pages)
}

func (rsb *ResultSetBuild) ResultSet() *model.ResultSet {
	return rsb.resultSet
}
