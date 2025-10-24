// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package infoschema

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	commonConstants "github.com/lindb/common/constants"
	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/coordinator/storage"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/meta"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
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

func (r *reader) ReadData(ctx context.Context, table string, expr tree.Expression) (rows [][]*types.Datum, err error) {
	predicate := newPredicate(ctx)
	if expr != nil {
		_ = expr.Accept(nil, predicate)
	}
	switch strings.ToLower(table) {
	case constants.TableEnv:
		rows, err = r.readEnv(predicate)
	case constants.TableMaster:
		rows = r.readMaster()
	case constants.TableBrokers:
		rows = r.readBroker()
	case constants.TableStorages:
		rows = r.readStorage()
	case constants.TableEngines:
		rows = r.readEngines()
	case constants.TableSchemata:
		rows = r.readSchemata()
	case constants.TableMetadataTypes:
		rows = r.readMetadataTypes()
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
	case constants.TableFunctions:
		rows, err = r.readFunctions()
	case constants.TableSnippets:
		rows, err = r.readSnippets()
	}
	return rows, err
}

func (r *reader) readEnv(predicate *predicate) (rows [][]*types.Datum, err error) {
	instance := predicate.getColumnValue(envSchema.Columns[0].Name) // instance
	if instance == "" {
		currentNode := r.metadataMgr.GetCurrentNode()
		instance = fmt.Sprintf("%s:%d", currentNode.HostIP, currentNode.HTTPPort)
	}
	keys := predicate.getColumnValues(envSchema.Columns[1].Name) // key
	envs, err := r.env(instance, keys)
	if err != nil {
		return nil, err
	}
	for _, env := range envs {
		rows = append(rows, types.MakeDatums(
			instance,    // instance
			env.Key,     // key
			env.Value,   // value
			env.Default, // default
		))
	}
	return
}

func (r *reader) readMaster() (rows [][]*types.Datum) {
	masterNode := r.metadataMgr.GetMaster()
	rows = append(rows, types.MakeDatums(
		masterNode.Node.HostIP,     // host_ip
		masterNode.Node.HostName,   // host_name
		masterNode.Node.HTTPPort,   // http
		masterNode.Node.Version,    // version
		masterNode.Node.OnlineTime, // online_time
		masterNode.ElectTime,       // elect_time
	))
	return
}

func (r *reader) readBroker() (rows [][]*types.Datum) {
	nodes := r.metadataMgr.GetBrokerNodes()
	now := time.Now().UnixMilli()
	for _, node := range nodes {
		rows = append(rows, types.MakeDatums(
			node.HostIP,     // host_ip
			node.HostName,   // host_name
			node.Version,    // version
			node.OnlineTime, // online_time
			time.Duration((now-node.OnlineTime)*1000_000), // uptime
			node.GRPCPort, // grpc
			node.HTTPPort, // http
		))
	}
	return
}

func (r *reader) readStorage() (rows [][]*types.Datum) {
	nodes := r.metadataMgr.GetStorageNodes()
	now := time.Now().UnixMilli()
	for _, node := range nodes {
		rows = append(rows, types.MakeDatums(
			node.ID,         // id
			node.HostIP,     // host_ip
			node.HostName,   // host_name
			node.Version,    // version
			node.OnlineTime, // online_time
			time.Duration((now-node.OnlineTime)*1000_000), // uptime
			node.GRPCPort, // grpc
			node.HTTPPort, // http
		))
	}
	return
}

func (r *reader) readEngines() (rows [][]*types.Datum) {
	// TODO: read from config file
	rows = [][]*types.Datum{
		types.MakeDatums(option.Metric, "DEFAULT"), // engine/support
		types.MakeDatums(option.Log, "NO"),
		types.MakeDatums(option.Trace, "NO"),
	}
	return
}

func (r *reader) readSchemata() (rows [][]*types.Datum) {
	databases := r.metadataMgr.GetDatabases()
	for _, database := range databases {
		rows = append(rows, types.MakeDatums(
			database.Name,          // schema_name
			database.Option.Engine, // engine
			database.String(),      // statement
		))
	}
	return
}

func (r *reader) readMetadataTypes() (rows [][]*types.Datum) {
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
	info, err0 := r.getStateMachineInfo(role, metadataType)
	if err0 != nil {
		return nil, err0
	}
	var data []byte
	switch strings.ToLower(source) {
	case "repo":
		rs, err0 := r.exploreStateRepoData(ctx, info)
		if err0 != nil {
			return nil, err0
		}
		data, _ = json.MarshalIndent(rs, "", "  ")
	case "state_machine":
		rs := r.exploreStateMachineDate(role, metadataType)
		data, _ = json.MarshalIndent(rs, "", "  ")
	}
	rows = append(rows, types.MakeDatums(
		role,                    // role
		metadataType,            // type
		strings.ToLower(source), // source
		string(data),            // data
	))
	return rows, err
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
		return inputRole == "" || strings.EqualFold(item.role, inputRole)
	})
	for _, role := range roleList {
		nodes := role.nodes
		if len(nodes) == 0 {
			continue
		}
		metrics, err0 := metricCli.FetchMetricData(nodes, names)
		if err0 != nil {
			return nil, err0
		}
		for name, metricList := range metrics {
			for _, metric := range metricList {
				for _, field := range metric.Fields {
					rows = append(rows, types.MakeDatums(
						role.role, // role
						name,      // name
						string(encoding.JSONMarshal(metric.Tags)), // tags
						field.Name,  // field name
						field.Type,  // field type
						field.Value, // field value
					))
				}
			}
		}
	}
	return rows, err
}

func (r *reader) readReplications(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.getColumnValue(replicationSchema.Columns[0].Name)
	if schema == "" {
		return nil, errors.New("table_schema not found in where clause")
	}
	state, err0 := r.getStateFromStorage("/state/replica", map[string]string{"db": schema}, func() any {
		var state []models.FamilyLogReplicaState
		return &state
	})
	if err0 != nil {
		return nil, err0
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
	return rows, err
}

func (r *reader) readMemoryDatabases(predicate *predicate) (rows [][]*types.Datum, err error) {
	schema := predicate.getColumnValue(memoryDatabaseSchema.Columns[0].Name)
	if schema == "" {
		return nil, errors.New("table_schema not found in where clause")
	}
	state, err0 := r.getStateFromStorage("/state/tsdb/memory", map[string]string{"db": schema}, func() any {
		var state []models.DataSegmentState
		return &state
	})
	if err0 != nil {
		return nil, err0
	}
	result := state.(map[string]any)
	for node, state := range result {
		familyStates := *(state.(*[]models.DataSegmentState))
		for _, familyState := range familyStates {
			for _, replicator := range familyState.MemoryDatabases {
				rows = append(rows, types.MakeDatums(
					schema,                  // table_schema
					node,                    // node
					familyState.ShardID,     // shard_id
					familyState.SegmentTime, // family_time
					replicator.State,        // state
					replicator.Uptime,       // uptime
					replicator.MemSize,      // mem_size
					replicator.NumOfSeries,  // num_of_series
				))
			}
		}
	}
	return rows, err
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

func (r *reader) readFunctions() (rows [][]*types.Datum, err error) {
	data, err := os.ReadFile(filepath.Join(getCurrentDir(), "functions.yaml"))
	if err != nil {
		return nil, err
	}

	var functions []Function
	err = yaml.Unmarshal(data, &functions)
	if err != nil {
		return nil, err
	}
	for _, f := range functions {
		rows = append(rows, types.MakeDatums(
			f.Name,     // name
			f.Template, // template
		))
	}
	return
}

func (r *reader) readSnippets() (rows [][]*types.Datum, err error) {
	data, err := os.ReadFile(filepath.Join(getCurrentDir(), "snippets.yaml"))
	if err != nil {
		return nil, err
	}

	var snippets []Snippet
	err = yaml.Unmarshal(data, &snippets)
	if err != nil {
		return nil, err
	}
	for _, s := range snippets {
		rows = append(rows, types.MakeDatums(
			s.Name,     // name
			s.Template, // template
		))
	}
	return
}

func getCurrentDir() string {
	_, dir, _, _ := runtime.Caller(0) //nolint
	return filepath.Dir(dir)
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

func (v *predicate) getColumnName(column tree.Expression) string {
	columnSymbols := plan.ExtractSymbolsFromExpression(column)
	if len(columnSymbols) != 1 {
		panic(fmt.Sprintf("column values lookup error, column: %s, symbol size: %d",
			tree.FormatExpression(column), len(columnSymbols)))
	}
	return columnSymbols[0].Name
}

func (v *predicate) Visit(ctx any, n tree.Node) (rs any) {
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		// TODO: check err
		columnName := v.getColumnName(node.Left)
		columnValue, _ := expression.EvalString(v.evalCtx, node.Right)
		v.addColumnValue(columnName, columnValue)
	case *tree.LikePredicate:
		// TODO: check err
		columnName := v.getColumnName(node.Value)
		columnValue, _ := expression.EvalString(v.evalCtx, node.Pattern)
		v.addColumnValue(columnName, columnValue)
	case *tree.LogicalExpression:
		for _, term := range node.Terms {
			_ = term.Accept(ctx, v)
		}
	case *tree.InPredicate:
		columnName := v.getColumnName(node.Value)
		if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
			lo.ForEach(inListExpression.Values, func(item tree.Expression, index int) {
				columnValue, _ := expression.EvalString(v.evalCtx, item)
				v.addColumnValue(columnName, columnValue)
			})
		}
	case *tree.Cast:
		_ = node.Expression.Accept(ctx, v)
	default:
		panic(fmt.Sprintf("infoschema predicate visit error, not support node type: %T", n))
	}
	return
}
