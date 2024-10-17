package infoschema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	commonConstants "github.com/lindb/common/constants"
	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/common/pkg/timeutil"
	"github.com/samber/lo"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/tree"
)

var (
	metricCli     = client.NewMetricCli()
	metadataPaths = map[string]map[string]models.StateMachineInfo{
		strings.ToLower(constants.BrokerRole):  broker.StateMachinePaths,
		strings.ToLower(constants.MasterRole):  master.StateMachinePaths,
		strings.ToLower(constants.StorageRole): storage.StateMachinePaths,
	}
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
	switch strings.ToLower(table) {
	case constants.TableEnv:
		rows, err = r.readEnv(predicate)
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
	case constants.TableMetadataTypes:
		rows, err = r.readMetadataTypes()
	case constants.TableMetadatas:
		rows, err = r.readMetadatas(ctx, predicate)
	case constants.TableMetrics:
		rows, err = r.readMetrics(predicate)
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

func (r *reader) readEnv(predicate *predicate) (rows [][]*types.Datum, err error) {
	instance := predicate.getColumnValue(envSchema.Columns[0].Name) // instance
	if instance == "" {
		return nil, errors.New("instance not found in where clause(ip:port)")
	}
	key := predicate.getColumnValue(envSchema.Columns[1].Name) // key
	envs, err := r.env(instance)
	if err != nil {
		return nil, err
	}
	for _, env := range envs {
		if key == "" || env.Key == key {
			rows = append(rows, types.MakeDatums(
				instance,    // instance
				env.Key,     // key
				env.Value,   // value
				env.Default, // default
			))
		}
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

func (r *reader) readMetadataTypes() (rows [][]*types.Datum, err error) {
	for role, paths := range metadataPaths {
		for key, info := range paths {
			rows = append(rows, types.MakeDatums(
				role,         // role
				key,          // type
				info.Comment, // comment
			))
		}
	}
	return
}

func (r *reader) getStateMachineInfo(role, metadataType string) (models.StateMachineInfo, error) {
	paths := metadataPaths[strings.ToLower(role)]
	info, ok := paths[metadataType]
	if !ok {
		return models.StateMachineInfo{}, errors.New("metadata type not found")
	}
	return info, nil
}

func (r *reader) readMetadatas(ctx context.Context, predicate *predicate) (rows [][]*types.Datum, err error) {
	role := predicate.getColumnValue(metadatasSchema.Columns[0].Name) // role
	if role == "" {
		return nil, errors.New("role not found in where clause(broker/master/storage)")
	}
	metadataType := predicate.getColumnValue(metadatasSchema.Columns[1].Name) // type
	if metadataType == "" {
		return nil, errors.New("type not found in where clause")
	}
	source := predicate.getColumnValue(metadatasSchema.Columns[2].Name) // source
	if source == "" {
		return nil, errors.New("source not found in where clause")
	}
	info, err := r.getStateMachineInfo(role, metadataType)
	if err != nil {
		return nil, err
	}
	var data []byte
	switch strings.ToLower(source) {
	case "repo":
		rs, err := r.exploreStateRepoData(ctx, info)
		if err != nil {
			return nil, err
		}
		data, _ = json.MarshalIndent(rs, "", "  ")
	case "state_machine":
		rs, err := r.exploreStateMachineDate(role, metadataType)
		if err != nil {
			return nil, err
		}
		data, _ = json.MarshalIndent(rs, "", "  ")
	}
	rows = append(rows, types.MakeDatums(
		role,                    // role
		metadataType,            // type
		strings.ToLower(source), // source
		string(data),            // data
	))
	return
}

func (r *reader) readMetrics(predicate *predicate) (rows [][]*types.Datum, err error) {
	inputRole := predicate.getColumnValue(metricsSchema.Columns[0].Name)
	names := predicate.getColumnValues(metricsSchema.Columns[1].Name)
	roles := []struct {
		role  string
		nodes []models.Node
	}{
		{
			constants.BrokerRole,
			lo.Map(r.metadataMgr.GetBrokerNodes(), func(item models.StatelessNode, index int) models.Node {
				return &item
			}),
		},
		{
			constants.StorageRole,
			lo.Map(r.metadataMgr.GetStorageNodes(), func(item models.StatefulNode, index int) models.Node {
				return &item
			}),
		},
	}

	roleList := lo.Filter(roles, func(item struct {
		role  string
		nodes []models.Node
	}, index int,
	) bool {
		return inputRole == "" || strings.ToUpper(item.role) == strings.ToUpper(inputRole)
	})
	for _, role := range roleList {
		nodes := role.nodes
		if len(nodes) == 0 {
			continue
		}
		metrics, err := metricCli.FetchMetricData(nodes, names)
		if err != nil {
			return nil, err
		}
		for name, metricList := range metrics {
			for _, metric := range metricList {
				for _, field := range metric.Fields {
					rows = append(rows, types.MakeDatums(
						role.role, // role
						name,      // name
						string(encoding.JSONMarshal(metric.Tags)), // tags
						field.Type,  // type
						field.Value, // value
					))
				}
			}
		}
	}
	return
}

func (r *reader) readReplications(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.getColumnValue(replicationSchema.Columns[0].Name)
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
					replicator.StateErrMsg,    // error
				))
			}
		}
	}
	return
}

func (r *reader) readMemoryDatabases(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.getColumnValue(memoryDatabaseSchema.Columns[0].Name)
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
	schema := predicate.getColumnValue(namespacesSchema.Columns[0].Name)
	if schema == "" {
		return nil, errors.New("table_schema not found in where clause")
	}
	namespace := predicate.getColumnValue(namespacesSchema.Columns[1].Name)
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
	schema := predicate.getColumnValue(tableNamesSchema.Columns[0].Name)
	namespace := predicate.getColumnValue(tableNamesSchema.Columns[1].Name)
	if namespace == "" {
		namespace = commonConstants.DefaultNamespace
	}
	tableName := predicate.getColumnValue(tableNamesSchema.Columns[2].Name)
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
	schema := predicate.getColumnValue(columnsSchema.Columns[0].Name)
	namespace := predicate.getColumnValue(columnsSchema.Columns[1].Name)
	if namespace == "" {
		namespace = commonConstants.DefaultNamespace
	}
	tableName := predicate.getColumnValue(columnsSchema.Columns[2].Name)
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
	columns map[string][]string
}

func newPredicate(ctx context.Context) *predicate {
	return &predicate{
		evalCtx: expression.NewEvalContext(ctx),
		columns: make(map[string][]string),
	}
}

func (v *predicate) getColumnValue(name string) string {
	values := v.getColumnValues(name)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (v *predicate) getColumnValues(name string) []string {
	return v.columns[strings.ToUpper(name)]
}

func (v *predicate) addColumnValue(name, value string) {
	colName := strings.ToUpper(name)
	values := v.columns[colName]
	v.columns[colName] = append(values, value)
}

func (v *predicate) Visit(context any, n tree.Node) (rs any) {
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		// TODO: check err
		columnName, _ := expression.EvalString(v.evalCtx, node.Left)
		columnValue, _ := expression.EvalString(v.evalCtx, node.Right)
		v.addColumnValue(columnName, columnValue)
	case *tree.LikePredicate:
		// TODO: check err
		columnName, _ := expression.EvalString(v.evalCtx, node.Value)
		columnValue, _ := expression.EvalString(v.evalCtx, node.Pattern)
		v.addColumnValue(columnName, columnValue)
	case *tree.LogicalExpression:
		for _, term := range node.Terms {
			_ = term.Accept(context, v)
		}
	case *tree.InPredicate:
		columnName, _ := expression.EvalString(v.evalCtx, node.Value)
		if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
			lo.ForEach(inListExpression.Values, func(item tree.Expression, index int) {
				columnValue, _ := expression.EvalString(v.evalCtx, item)
				v.addColumnValue(columnName, columnValue)
			})
		}
	case *tree.Cast:
		_ = node.Expression.Accept(context, v)
	default:
		panic(fmt.Sprintf("infoschema predicate visit error, not support node type: %T", n))
	}
	return
}
