package infoschema

import (
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
	schemata, err := p.reader.ReadSchemata()
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
	for _, row := range schemata {
		for idx, col := range columns {
			switch p.outputs[idx].DataType {
			case types.DTString:
				col.AppendString(row[p.split.colIdxs[idx]].String())
			}
		}
	}
	return page
}
