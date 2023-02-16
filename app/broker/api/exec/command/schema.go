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
	"fmt"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/validate"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

// SchemaCommand executes database schema statement.
func SchemaCommand(ctx context.Context, deps *depspkg.HTTPDeps, _ *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	schemaStmt := stmt.(*stmtpkg.Schema)
	switch schemaStmt.Type {
	case stmtpkg.DatabaseSchemaType:
		return listDataBases(ctx, deps)
	case stmtpkg.CreateDatabaseSchemaType:
		return saveDataBase(ctx, deps, schemaStmt)
	case stmtpkg.DropDatabaseSchemaType:
		return dropDatabase(ctx, deps, schemaStmt)
	case stmtpkg.DatabaseNameSchemaType:
		dbs, err := listDataBases(ctx, deps)
		if err != nil {
			return nil, err
		}
		var databaseNames []interface{}
		databases := dbs.([]*models.Database)
		for _, db := range databases {
			databaseNames = append(databaseNames, db.Name)
		}
		return databaseNames, nil
	}
	return nil, nil
}

// dropDatabase drops database config.
func dropDatabase(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Schema) (interface{}, error) {
	databaseName := stmt.Value
	log.Info("drop database", logger.String("name", databaseName))
	if err := deps.Repo.Delete(ctx, constants.GetDatabaseConfigPath(databaseName)); err != nil {
		return nil, err
	}
	if err := deps.Repo.Delete(ctx, constants.GetDatabaseAssignPath(databaseName)); err != nil {
		return nil, err
	}
	// TODO: remove limits
	rs := fmt.Sprintf("Drop database[%s] ok", stmt.Value)
	return &rs, nil
}

// listDataBases returns database list in cluster.
func listDataBases(ctx context.Context, deps *depspkg.HTTPDeps) (interface{}, error) {
	data, err := deps.Repo.List(ctx, constants.DatabaseConfigPath)
	if err != nil {
		return nil, err
	}
	var dbs []*models.Database
	for _, val := range data {
		db := &models.Database{}
		err = encoding.JSONUnmarshal(val.Value, db)
		if err != nil {
			log.Warn("unmarshal data error",
				logger.String("data", string(val.Value)))
			continue
		}
		db.Desc = db.String()
		dbs = append(dbs, db)
	}
	return dbs, nil
}

// saveDataBase creates the database config if there is no database
// config with the name database.Name, otherwise update the config.
func saveDataBase(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Schema) (interface{}, error) {
	data := []byte(stmt.Value)
	database := &models.Database{}
	err := encoding.JSONUnmarshal(data, database)
	if err != nil {
		return nil, err
	}
	err = validate.Validator.Struct(database)
	if err != nil {
		return nil, err
	}

	opt := database.Option
	// validate time series engine option
	if err := opt.Validate(); err != nil {
		return nil, err
	}
	// set default value
	opt.Default()
	database.Option = opt // reset option after set default value

	log.Info("Saving Database", logger.String("config", stmt.Value))
	if err := deps.Repo.Put(ctx, constants.GetDatabaseConfigPath(database.Name), data); err != nil {
		return nil, err
	}
	rs := "Create database ok"
	return &rs, nil
}
