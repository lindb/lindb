package log

import (
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/store"
)

type TableScan struct {
	db        store.Database
	timeRange timeutil.TimeRange
}
