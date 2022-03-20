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

package brokerquery

import (
	"context"

	"github.com/lindb/lindb/coordinator/broker"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

type queryFactory struct {
	stateMgr    broker.StateManager
	taskManager TaskManager
}

func NewQueryFactory(
	stateMgr broker.StateManager,
	taskManager TaskManager,
) Factory {
	return &queryFactory{
		stateMgr:    stateMgr,
		taskManager: taskManager,
	}
}

func (qh *queryFactory) NewMetricQuery(
	ctx context.Context,
	databaseName string,
	sql *stmtpkg.Query,
) MetricQuery {
	return newMetricQuery(ctx, databaseName, sql, qh)
}

func (qh *queryFactory) NewMetadataQuery(
	ctx context.Context,
	database string,
	stmt *stmtpkg.MetricMetadata,
) MetaDataQuery {
	return newMetadataQuery(ctx, database, stmt, qh)
}
