package infoschema

import (
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi/types"
)

type Reader interface {
	ReadSchemata() (rows [][]*types.Datum, err error)
}

type reader struct {
	metadataMgr meta.MetadataManager
}

func NewReader(metadataMgr meta.MetadataManager) Reader {
	return &reader{metadataMgr: metadataMgr}
}

func (r *reader) ReadSchemata() (rows [][]*types.Datum, err error) {
	databases := r.metadataMgr.GetDatabases()
	for _, database := range databases {
		rows = append(rows, types.MakeDatums(
			database.Name, // schema_name
			"METRIC",      // FIXME: engine
		))
	}
	return
}
