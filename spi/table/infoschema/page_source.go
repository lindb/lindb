package infoschema

import (
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/spi/types"
)

type PageSourceProvider struct {
	reader Reader
}

func NewPageSourceProvider(reader Reader) spi.PageSourceProvider {
	return &PageSourceProvider{
		reader: reader,
	}
}

// CreatePageSource implements spi.PageSourceProvider.
func (p *PageSourceProvider) CreatePageSource(table spi.TableHandle,
	outputs []types.ColumnMetadata, assignments []*spi.ColumnAssignment,
) spi.PageSource {
	return &PageSource{
		outputs: outputs,
		reader:  p.reader,
	}
}

type PageSource struct {
	reader  Reader
	split   *InfoSplit
	outputs []types.ColumnMetadata
}

// AddSplit implements spi.PageSource.
func (p *PageSource) AddSplit(split spi.Split) {
	if info, ok := split.(*InfoSplit); ok {
		p.split = info
	}
}

// GetNextPage implements spi.PageSource.
func (p *PageSource) GetNextPage() *types.Page {
	var (
		rows [][]*types.Datum
		err  error
	)
	switch p.split.table {
	case constants.TableMaster:
		rows, err = p.reader.ReadMaster()
	case constants.TableSchemata:
		rows, err = p.reader.ReadSchemata()
	case constants.TableMetrics:
		rows, err = p.reader.ReadMetrics()
	}
	if err != nil {
		panic(err)
	}
	page := types.NewPage()
	var columns []*types.Column
	for _, output := range p.outputs {
		column := types.NewColumn()
		page.AppendColumn(output, column)
		columns = append(columns, column)
	}
	fmt.Printf("ouputs=%v,index=%v\n", p.outputs, p.split.colIdxs)
	for _, row := range rows {
		for idx, col := range columns {
			switch p.outputs[idx].DataType {
			case types.DTString:
				col.AppendString(row[p.split.colIdxs[idx]].String())
			case types.DTFloat:
				col.AppendFloat(row[p.split.colIdxs[idx]].Float())
			case types.DTInt:
				col.AppendInt(row[p.split.colIdxs[idx]].Int())
			}
		}
	}
	return page
}
