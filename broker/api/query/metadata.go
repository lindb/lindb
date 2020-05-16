package query

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/lindb/lindb/broker/api"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/service"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
)

// for testing
var (
	parseSQLFunc = parseSQL
)

var errWrongQueryStmt = errors.New("can't parse metadata query ql")
var errUnknownMetadataStmt = errors.New("unknown metadata statement")

// MetadataAPI represents metadata query api
type MetadataAPI struct {
	databaseService     service.DatabaseService
	replicaStateMachine replica.StatusStateMachine
	nodeStateMachine    broker.NodeStateMachine
	executorFactory     parallel.ExecutorFactory
	jobManager          parallel.JobManager
}

// NewDatabaseAPI creates database api instance
func NewMetadataAPI(databaseService service.DatabaseService,
	replicaStateMachine replica.StatusStateMachine, nodeStateMachine broker.NodeStateMachine,
	executorFactory parallel.ExecutorFactory, jobManager parallel.JobManager,
) *MetadataAPI {
	return &MetadataAPI{
		databaseService:     databaseService,
		replicaStateMachine: replicaStateMachine,
		nodeStateMachine:    nodeStateMachine,
		executorFactory:     executorFactory,
		jobManager:          jobManager,
	}
}

// Handle handles metadata query by LinQL
func (d *MetadataAPI) Handle(w http.ResponseWriter, r *http.Request) {
	ql, err := api.GetParamsFromRequest("sql", r, "", true)
	if err != nil {
		api.Error(w, err)
		return
	}
	metaQuery, err := parseSQLFunc(ql)
	if err != nil {
		api.Error(w, err)
		return
	}
	switch metaQuery.Type {
	case stmt.Database:
		d.showDatabases(w)
	case stmt.Namespace, stmt.Metric, stmt.Field, stmt.TagKey, stmt.TagValue:
		db, err := api.GetParamsFromRequest("db", r, "", true)
		if err != nil {
			api.Error(w, err)
			return
		}
		d.suggest(w, db, metaQuery)
	default:
		api.Error(w, errUnknownMetadataStmt)
	}
}

// showDatabases shows all database names
func (d *MetadataAPI) showDatabases(w http.ResponseWriter) {
	databases, err := d.databaseService.List()
	if err != nil {
		api.Error(w, err)
		return
	}
	var databaseNames []string
	for _, db := range databases {
		databaseNames = append(databaseNames, db.Name)
	}
	api.OK(w, databaseNames)
}

// suggest executes the suggest query
func (d *MetadataAPI) suggest(w http.ResponseWriter, database string, request *stmt.Metadata) {
	//TODO add timeout cfg
	ctx, cancel := context.WithTimeout(context.TODO(), time.Minute)
	defer cancel()

	exec := d.executorFactory.NewMetadataBrokerExecutor(ctx, database, request,
		d.replicaStateMachine, d.nodeStateMachine, d.jobManager)
	values, err := exec.Execute()
	if err != nil {
		api.Error(w, err)
		return
	}
	switch request.Type {
	case stmt.Field:
		// build field result model
		result := make(map[string]string)
		fields := field.Metas{}
		for _, value := range values {
			err = encoding.JSONUnmarshal([]byte(value), &fields)
			if err != nil {
				api.Error(w, err)
				return
			}
			for _, f := range fields {
				result[f.Name] = f.Type.String()
			}
		}
		api.OK(w, result)
	default:
		api.OK(w, values)
	}
}

// parseSQL parses metadata query sql
func parseSQL(ql string) (*stmt.Metadata, error) {
	query, err := sql.Parse(ql)
	if err != nil {
		return nil, err
	}
	metaQuery, ok := query.(*stmt.Metadata)
	if !ok {
		return nil, errWrongQueryStmt
	}
	return metaQuery, nil
}
