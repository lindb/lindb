package infoschema

import (
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi/types"
)

type Reader interface {
	ReadMaster() (rows [][]*types.Datum, err error)
	ReadSchemata() (rows [][]*types.Datum, err error)
	ReadMetrics() (rows [][]*types.Datum, err error)
}

// reader implements Reader interface.
// schema of rows returned ref to: tables.go
type reader struct {
	metadataMgr meta.MetadataManager
}

func NewReader(metadataMgr meta.MetadataManager) Reader {
	return &reader{metadataMgr: metadataMgr}
}

func (r *reader) ReadMaster() (rows [][]*types.Datum, err error) {
	master := r.metadataMgr.GetMaster()
	rows = append(rows, types.MakeDatums(
		master.Node.HostIP,     // host_ip
		master.Node.HostName,   // host_name
		master.Node.Version,    // version
		master.Node.OnlineTime, // online_time
		master.ElectTime,       // elect_time
	))
	return
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

func (r reader) ReadMetrics() (rows [][]*types.Datum, err error) {
	rows = append(rows, types.MakeDatums(
		"cpu",         // metrics_name
		"1.1.1.1",     // instance
		float64(10.3), // value
	))
	return
}
