package infoschema

import (
	"fmt"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/spi"
)

func init() {
	// register table handle
	encoding.RegisterNodeType(TableHandle{})
	spi.RegisterCreateTableFn(spi.InfoSchema, func(db, ns, name string) spi.TableHandle {
		return &TableHandle{
			Table: name,
		}
	})
}

type TableHandle struct {
	Table string `json:"table"`
}

func (t *TableHandle) SetTimeRange(timeRange timeutil.TimeRange) {}

func (t *TableHandle) GetTimeRange() timeutil.TimeRange {
	panic(constants.ErrNotSupportOperation)
}

// Kind returns the datasource kind.
func (t *TableHandle) Kind() spi.DatasourceKind {
	return spi.InfoSchema
}

// String returns the table info of information schema.
func (t *TableHandle) String() string {
	return fmt.Sprintf("%s.%s", constants.InformationSchema, t.Table)
}

type InfoSplit struct {
	table   string
	colIdxs []int
}
