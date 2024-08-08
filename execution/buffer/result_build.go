package buffer

import (
	"fmt"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/execution/model"
	"github.com/lindb/lindb/spi"
)

type ResultSetBuild struct {
	pages     chan *spi.Page
	completed chan struct{}
	resultSet *model.ResultSet
}

func CreateResultSetBuild() *ResultSetBuild {
	return &ResultSetBuild{
		pages:     make(chan *spi.Page),
		completed: make(chan struct{}),
		resultSet: model.NewResultSet(),
	}
}

func (rsb *ResultSetBuild) AddPage(page *spi.Page) {
	fmt.Println("add result page")
	rsb.pages <- page
}

func (rsb *ResultSetBuild) Process() {
	defer func() {
		close(rsb.completed)
	}()
	// TODO: need close when timeout
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
		fmt.Println(string(encoding.JSONMarshal(rsb.resultSet)))
		fmt.Println("merge result page")
	}
}

func (rsb *ResultSetBuild) Complete() {
	fmt.Println("close result page")
	close(rsb.pages)
	// waiting process result page completed
	<-rsb.completed
}

func (rsb *ResultSetBuild) ResultSet() *model.ResultSet {
	fmt.Println("result.....")
	fmt.Println(string(encoding.JSONMarshal(rsb.resultSet)))
	return rsb.resultSet
}
