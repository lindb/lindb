package buffer

import (
	"fmt"

	"github.com/lindb/common/pkg/encoding"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/model"
)

type ResultSetBuild struct {
	pages     chan *types.Page
	completed chan struct{}
	resultSet *model.ResultSet
}

func CreateResultSetBuild() *ResultSetBuild {
	return &ResultSetBuild{
		pages:     make(chan *types.Page),
		completed: make(chan struct{}),
		resultSet: model.NewResultSet(),
	}
}

func (rsb *ResultSetBuild) AddPage(page *types.Page) {
	rsb.pages <- page
}

func (rsb *ResultSetBuild) Process() {
	defer func() {
		close(rsb.completed)
	}()
	// TODO: need close when timeout
	for page := range rsb.pages {
		if len(rsb.resultSet.Schema.Columns) == 0 {
			rsb.resultSet.Schema.Columns = page.Layout
		}
		it := page.Iterator()
		for row := it.Begin(); row != it.End(); row = it.Next() {
			columns := make([]any, len(page.Layout))
			for i, meta := range page.Layout {
				// TODO: add more type
				switch meta.DataType {
				case types.DTString:
					columns[i] = row.GetString(i)
				case types.DTInt:
					columns[i] = row.GetInt(i)
				case types.DTFloat:
					columns[i] = row.GetFloat(i)
				case types.DTTimeSeries:
					columns[i] = row.GetTimeSeries(i)
				default:
					panic(fmt.Sprintf("unknown data type:%v", meta.DataType))
				}
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
