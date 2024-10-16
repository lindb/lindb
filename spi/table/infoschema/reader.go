package infoschema

import (
	"context"
	"errors"
	"fmt"
	"strings"

	commonConstants "github.com/lindb/common/constants"
	"github.com/lindb/common/pkg/logger"
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

	logger logger.Logger
}

func NewReader(metadataMgr meta.MetadataManager) Reader {
	return &reader{metadataMgr: metadataMgr, logger: logger.GetLogger("Infoschema", "Reader")}
}

func (r *reader) ReadData(ctx context.Context, table string, expression tree.Expression) (rows [][]*types.Datum, err error) {
	predicate := newPredicate(ctx)
	if expression != nil {
		_ = expression.Accept(nil, predicate)
	}
	fmt.Printf("read meta data table=%v\n", table)

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
	case constants.TableReplications:
		rows, err = r.readReplications(predicate)
	case constants.TableMemoryDatabases:
		rows, err = r.readMemoryDatabases(predicate)
	case constants.TableNamespaces:
		rows, err = r.readNamespaces(predicate)
	case constants.TableTableNames:
		rows, err = r.readTableNames(predicate)
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

func (r *reader) readReplications(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.columns[replicationSchema.Columns[0].Name]
	if schema == "" {
		return nil, errors.New("table_schema not found in where clause")
	}
	state, err := r.getStateFromStorage("/state/replica", map[string]string{"db": schema}, func() any {
		var state []models.FamilyLogReplicaState
		return &state
	})
	if err != nil {
		return nil, err
	}
	result := state.(map[string]any)
	for node, state := range result {
		familyWALLogs := *(state.(*[]models.FamilyLogReplicaState))
		for _, familyWALLog := range familyWALLogs {
			for _, replicator := range familyWALLog.Replicators {
				rows = append(rows, types.MakeDatums(
					schema,                    // table_schema
					node,                      // node
					familyWALLog.ShardID,      // shard_id
					familyWALLog.FamilyTime,   // family_time
					familyWALLog.Leader,       // leader
					replicator.Replicator,     // replicator
					replicator.ReplicatorType, // replicator_type
					familyWALLog.Append,       // append
					replicator.Consume,        // consume
					replicator.ACK,            // ack
					replicator.Pending,        // pending
					replicator.State.String(), // state
					replicator.StateErrMsg,    // err_msg
				))
			}
		}
	}
	return
}

func (r *reader) readMemoryDatabases(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.columns[replicationSchema.Columns[0].Name]
	if schema == "" {
		return nil, errors.New("table_schema not found in where clause")
	}
	state, err := r.getStateFromStorage("/state/tsdb/memory", map[string]string{"db": schema}, func() any {
		var state []models.DataFamilyState
		return &state
	})
	if err != nil {
		return nil, err
	}
	result := state.(map[string]any)
	for node, state := range result {
		familyStates := *(state.(*[]models.DataFamilyState))
		for _, familyState := range familyStates {
			for _, replicator := range familyState.MemoryDatabases {
				rows = append(rows, types.MakeDatums(
					schema,                 // table_schema
					node,                   // node
					familyState.ShardID,    // shard_id
					familyState.FamilyTime, // family_time
					replicator.State,       // state
					replicator.Uptime,      // uptime
					replicator.MemSize,     // mem_size
					replicator.NumOfSeries, // num_of_series
				))
			}
		}
	}
	return
}

func (r *reader) readNamespaces(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.columns[columnsSchema.Columns[0].Name]
	if schema == "" {
		return nil, errors.New("table_schema not found in where clause")
	}
	namespace := predicate.columns[columnsSchema.Columns[1].Name]
	namespaces, err := r.suggestNamespaces(schema, namespace, 10)
	if err != nil {
		return nil, err
	}
	for _, ns := range namespaces {
		rows = append(rows, types.MakeDatums(
			schema, // table_schema
			ns,     // namespace
		))
	}
	return
}

func (r *reader) readTableNames(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.columns[columnsSchema.Columns[0].Name]
	namespace := predicate.columns[columnsSchema.Columns[1].Name]
	if namespace == "" {
		namespace = commonConstants.DefaultNamespace
	}
	tableName := predicate.columns[columnsSchema.Columns[2].Name]
	if schema == "" {
		return nil, errors.New("table_schema not found in where clause")
	}
	tableNames, err := r.suggestTables(schema, namespace, tableName, 10) // FIXME: set limit
	if err != nil {
		return nil, err
	}
	for _, name := range tableNames {
		rows = append(rows, types.MakeDatums(
			schema,    // table_schema
			namespace, // namespace
			name,      // table_name
		))
	}
	return
}

func (r *reader) readColumns(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.columns[columnsSchema.Columns[0].Name]
	namespace := predicate.columns[columnsSchema.Columns[1].Name]
	if namespace == "" {
		namespace = commonConstants.DefaultNamespace
	}
	tableName := predicate.columns[columnsSchema.Columns[2].Name]
	if schema == "" || tableName == "" {
		return nil, errors.New("table_schema/table_name not found in where clause")
	}
	table, err := r.metadataMgr.GetTableMetadata(schema, namespace, tableName)
	if err != nil {
		return nil, err
	}
	// add time(reserved column)
	table.Schema.AddColumn(types.ColumnMetadata{Name: "time", DataType: types.DTTimestamp})
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
	case *tree.LikePredicate:
		// TODO: check err
		columnName, _ := expression.EvalString(v.evalCtx, node.Value)
		columnValue, _ := expression.EvalString(v.evalCtx, node.Pattern)
		v.columns[columnName] = columnValue
	case *tree.LogicalExpression:
		for _, term := range node.Terms {
			_ = term.Accept(context, v)
		}
	case *tree.Cast:
		_ = node.Expression.Accept(context, v)
	default:
		panic(fmt.Sprintf("infoschema predicate visit error, not support node type: %T", n))
	}
	return
}
