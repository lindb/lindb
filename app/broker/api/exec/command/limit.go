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
	"strings"

	"github.com/BurntSushi/toml"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

var (
	tomlDecodeFn = toml.Decode
)

// LimitCommand executes database limit statement.
func LimitCommand(ctx context.Context, deps *depspkg.HTTPDeps, param *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	db := strings.TrimSpace(param.Database)
	if db == "" {
		return nil, constants.ErrDatabaseNameRequired
	}
	limitStmt := stmt.(*stmtpkg.Limit)
	switch limitStmt.Type {
	case stmtpkg.SetLimit:
		return setLimit(ctx, db, deps, limitStmt)
	case stmtpkg.ShowLimit:
		return showLimit(ctx, db, deps)
	}
	return nil, nil
}

// showLimit returns database's limits.
func showLimit(ctx context.Context, db string, deps *depspkg.HTTPDeps) (interface{}, error) {
	data, err := deps.Repo.Get(ctx, constants.GetDatabaseLimitPath(db))
	if err == state.ErrNotExist {
		limit := models.NewDefaultLimits().TOML()
		return &limit, nil
	}
	if err != nil {
		return nil, err
	}
	limit := string(data)
	return &limit, nil
}

// setLimit set the database's limits.
func setLimit(ctx context.Context, db string, deps *depspkg.HTTPDeps, stmt *stmtpkg.Limit) (interface{}, error) {
	data := []byte(stmt.Limit)
	limits := &models.Limits{}
	// check limit if valid
	_, err := tomlDecodeFn(string(data), limits)
	if err != nil {
		return nil, err
	}
	if err := deps.Repo.Put(ctx, constants.GetDatabaseLimitPath(db), data); err != nil {
		return nil, err
	}
	rs := "set limit ok"
	return &rs, nil
}
