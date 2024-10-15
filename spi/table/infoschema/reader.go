package infoschema

import (
	"context"
	"fmt"
	"strings"

	commonConstants "github.com/lindb/common/constants"
	"github.com/lindb/common/pkg/timeutil"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/tree"
)

type Reader interface {
	ReadData(ctx context.Context, table string, predicate tree.Expression) (rows [][]*types.Datum, err error)
}

// reader implements Reader interface.
// schema of rows returned ref to: tables.go
type reader struct {
	metadataMgr meta.MetadataManager
}

func NewReader(metadataMgr meta.MetadataManager) Reader {
	return &reader{metadataMgr: metadataMgr}
}

func (r *reader) ReadData(ctx context.Context, table string, expression tree.Expression) (rows [][]*types.Datum, err error) {
	predicate := newPredicate(ctx)
	if expression != nil {
		_ = expression.Accept(nil, predicate)
	}

	switch strings.ToLower(table) {
	case constants.TableMaster:
		rows, err = r.readMaster()
	case constants.TableBroker:
		rows, err = r.readBroker()
	case constants.TableStorage:
		rows, err = r.readStorage()
	case constants.TableEngines:
		rows, err = r.readEngines()
	case constants.TableSchemata:
		rows, err = r.readSchemata()
	case constants.TableMetrics:
		rows, err = r.readMetrics()
	case constants.TableColumns:
		rows, err = r.readColumns(predicate)
	}
	return
}

func (r *reader) readMaster() (rows [][]*types.Datum, err error) {
	master := r.metadataMgr.GetMaster()
	rows = append(rows, types.MakeDatums(
		master.Node.HostIP,   // host_ip
		master.Node.HostName, // host_name
		master.Node.Version,  // version
		timeutil.FormatTimestamp(master.Node.OnlineTime, timeutil.DataTimeFormat2), // online_time
		timeutil.FormatTimestamp(master.ElectTime, timeutil.DataTimeFormat2),       // elect_time
	))
	return
}

func (r *reader) readBroker() (rows [][]*types.Datum, err error) {
	nodes := r.metadataMgr.GetBrokerNodes()
	for _, node := range nodes {
		rows = append(rows, types.MakeDatums(
			node.HostIP,   // host_ip
			node.HostName, // host_name
			node.Version,  // version
			timeutil.FormatTimestamp(node.OnlineTime, timeutil.DataTimeFormat2), // online_time
			node.GRPCPort, // grpc
			node.HTTPPort, // http
		))
	}
	return
}

func (r *reader) readStorage() (rows [][]*types.Datum, err error) {
	nodes := r.metadataMgr.GetStorageNodes()
	for _, node := range nodes {
		rows = append(rows, types.MakeDatums(
			node.ID,       // id
			node.HostIP,   // host_ip
			node.HostName, // host_name
			node.Version,  // version
			timeutil.FormatTimestamp(node.OnlineTime, timeutil.DataTimeFormat2), // online_time
			node.GRPCPort, // grpc
			node.HTTPPort, // http
		))
	}
	return
}

func (r *reader) readEngines() (rows [][]*types.Datum, err error) {
	rows = [][]*types.Datum{
		types.MakeDatums(models.Metric, "DEFAULT"), // engine/support
		types.MakeDatums(models.Log, "NO"),
		types.MakeDatums(models.Trace, "NO"),
	}
	return
}

func (r *reader) readSchemata() (rows [][]*types.Datum, err error) {
	databases := r.metadataMgr.GetDatabases()
	for _, database := range databases {
		rows = append(rows, types.MakeDatums(
			database.Name,   // schema_name
			database.Engine, // engine
		))
	}
	return
}

func (r *reader) readMetrics() (rows [][]*types.Datum, err error) {
	rows = append(rows, types.MakeDatums(
		"cpu",         // metrics_name
		"1.1.1.1",     // instance
		float64(10.3), // value
	))
	return
}

func (r *reader) readColumns(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.columns[columnsSchema.Columns[0].Name]
	namespace := predicate.columns[columnsSchema.Columns[1].Name]
	if namespace == "" {
		namespace = commonConstants.DefaultNamespace
	}
	tableName := predicate.columns[columnsSchema.Columns[2].Name]
	table, err := r.metadataMgr.GetTableMetadata(schema, namespace, tableName)
	if err != nil {
		return nil, err
	}
	for _, column := range table.Schema.Columns {
		rows = append(rows, types.MakeDatums(
			schema,                   // table_schema
			namespace,                // namespace
			tableName,                // table_name
			column.Name,              // column_name
			column.DataType.String(), // data_type
			column.AggType.String(),  // agg_type
		))
	}
	return
}

type predicate struct {
	evalCtx expression.EvalContext
	columns map[string]string
}

func newPredicate(ctx context.Context) *predicate {
	return &predicate{
		evalCtx: expression.NewEvalContext(ctx),
		columns: make(map[string]string),
	}
}

func (v *predicate) Visit(context any, n tree.Node) (rs any) {
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		// TODO: check err
		columnName, _ := expression.EvalString(v.evalCtx, node.Left)
		columnValue, _ := expression.EvalString(v.evalCtx, node.Right)
		v.columns[columnName] = columnValue
	case *tree.LogicalExpression:
		for _, term := range node.Terms {
			_ = term.Accept(context, v)
		}
	default:
		panic(fmt.Sprintf("infoschema predicate visit error, not support node type: %T", n))
	}
	return
}
