package spi

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
)

func init() {
	// register json encoder/decoder for table handle
	jsoniter.RegisterTypeEncoder("spi.TableHandle", &encoding.JSONEncoder[TableHandle]{})
	jsoniter.RegisterTypeDecoder("spi.TableHandle", &encoding.JSONDecoder[TableHandle]{})
}

type DatasourceKind int

const (
	InfoSchema DatasourceKind = iota + 1
	Metric
)

// TableHandle represents a table handle that connect the storage engine.
type TableHandle interface {
	// SetTimeRange sets the time range of table handle.
	SetTimeRange(timeRange timeutil.TimeRange)
	// Kind returns the kind of data source.
	Kind() DatasourceKind
	// String returns table info, format: ${database}:${namespace}:${tableName}
	String() string
}
