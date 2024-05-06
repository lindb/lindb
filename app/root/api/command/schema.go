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

package command

import (
	"context"

	"github.com/lindb/common/pkg/encoding"
	"github.com/lindb/common/pkg/logger"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/validate"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// schemaCommandFn represents schema command function define.
type schemaCommandFn = func(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Schema) (interface{}, error)

// schemaCommands registers all schema related commands.
var schemaCommands = map[stmtpkg.SchemaType]schemaCommandFn{
	stmtpkg.CreateDatabaseSchemaType: saveDatabase,
	stmtpkg.DatabaseSchemaType:       listDatabases,
	stmtpkg.DatabaseNameSchemaType:   listDatabaseNames,
}

// SchemaCommand executes lin query language for broker related.
func SchemaCommand(ctx context.Context, deps *depspkg.HTTPDeps, _ *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	schemaStmt := stmt.(*stmtpkg.Schema)
	if commandFn, ok := schemaCommands[schemaStmt.Type]; ok {
		return commandFn(ctx, deps, schemaStmt)
	}
	return nil, nil
}

// listDatabaseNames returns database name list in cluster.
func listDatabaseNames(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Schema) (interface{}, error) {
	dbs, err := listDatabases(ctx, deps, stmt)
	if err != nil {
		return nil, err
	}
	var databaseNames []interface{}
	databases := dbs.([]*models.LogicDatabase)
	for _, db := range databases {
		databaseNames = append(databaseNames, db.Name)
	}
	return databaseNames, nil
}

// listDatabases returns database list in cluster.
func listDatabases(ctx context.Context, deps *depspkg.HTTPDeps, _ *stmtpkg.Schema) (interface{}, error) {
	data, err := deps.Repo.List(ctx, constants.DatabaseConfigPath)
	if err != nil {
		return nil, err
	}
	var dbs []*models.LogicDatabase
	for _, val := range data {
		db := &models.LogicDatabase{}
		err = encoding.JSONUnmarshal(val.Value, db)
		if err != nil {
			log.Warn("unmarshal data error",
				logger.String("data", string(val.Value)))
			continue
		}
		dbs = append(dbs, db)
	}
	return dbs, nil
}

// saveDatabase creates the database config if there is no database
// config with the name database.Name, otherwise update the config.
func saveDatabase(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Schema) (interface{}, error) {
	data := []byte(stmt.Value)
	database := &models.LogicDatabase{}
	err := encoding.JSONUnmarshal(data, database)
	if err != nil {
		return nil, err
	}
	err = validate.Validator.Struct(database)
	if err != nil {
		return nil, err
	}

	log.Info("Saving Database", logger.String("config", stmt.Value))
	if err := deps.Repo.Put(ctx, constants.GetDatabaseConfigPath(database.Name), data); err != nil {
		return nil, err
	}
	rs := "Create database ok"
	return &rs, nil
}
