package infoschema

import (
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/spi/types"
)

type Reader interface {
	ReadData(table string) (rows [][]*types.Datum, err error)
}

// reader implements Reader interface.
// schema of rows returned ref to: tables.go
type reader struct {
	metadataMgr meta.MetadataManager
}

func NewReader(metadataMgr meta.MetadataManager) Reader {
	return &reader{metadataMgr: metadataMgr}
}

func (r *reader) ReadData(table string) (rows [][]*types.Datum, err error) {
	switch table {
	case constants.TableMaster:
		rows, err = r.readMaster()
	case constants.TableBroker:
		rows, err = r.readBroker()
	case constants.TableStorage:
		rows, err = r.readStorage()
	case constants.TableSchemata:
		rows, err = r.readSchemata()
	case constants.TableMetrics:
		rows, err = r.readMetrics()
	}
	return
}

func (r *reader) readMaster() (rows [][]*types.Datum, err error) {
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

func (r *reader) readBroker() (rows [][]*types.Datum, err error) {
	nodes := r.metadataMgr.GetBrokerNodes()
	for _, node := range nodes {
		rows = append(rows, types.MakeDatums(
			node.HostIP,     // host_ip
			node.HostName,   // host_name
			node.Version,    // version
			node.OnlineTime, // online_time
			node.GRPCPort,   // grpc
			node.HTTPPort,   // http
		))
	}
	return
}

func (r *reader) readStorage() (rows [][]*types.Datum, err error) {
	nodes := r.metadataMgr.GetStorageNodes()
	for _, node := range nodes {
		rows = append(rows, types.MakeDatums(
			node.ID,         // id
			node.HostIP,     // host_ip
			node.HostName,   // host_name
			node.Version,    // version
			node.OnlineTime, // online_time
			node.GRPCPort,   // grpc
			node.HTTPPort,   // http
		))
	}
	return
}

func (r *reader) readSchemata() (rows [][]*types.Datum, err error) {
	databases := r.metadataMgr.GetDatabases()
	for _, database := range databases {
		rows = append(rows, types.MakeDatums(
			database.Name, // schema_name
			"METRIC",      // FIXME: engine
		))
	}
	return
}

func (r reader) readMetrics() (rows [][]*types.Datum, err error) {
	rows = append(rows, types.MakeDatums(
		"cpu",         // metrics_name
		"1.1.1.1",     // instance
		float64(10.3), // value
	))
	return
}
