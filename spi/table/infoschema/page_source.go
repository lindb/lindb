package infoschema

import (
	"context"
	"fmt"

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
func (p *PageSourceProvider) CreatePageSource(ctx context.Context, table spi.TableHandle,
	outputs []types.ColumnMetadata, assignments []*spi.ColumnAssignment,
) spi.PageSource {
	return &PageSource{
		ctx:     ctx,
		outputs: outputs,
		reader:  p.reader,
	}
}

type PageSource struct {
	ctx     context.Context
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
	rows, err := p.reader.ReadData(p.ctx, p.split.table, p.split.predicate)
	if err != nil {
		panic(err)
	}
	fmt.Printf("info schema: rows=%v\n", rows)
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
