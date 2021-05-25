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

package query

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestExecutorFactory_NewExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := NewExecutorFactory()
	mockDatabase := tsdb.NewMockDatabase(ctrl)
	assert.NotNil(t, factory.NewStorageExecutor(nil, mockDatabase, newStorageExecuteContext(nil, &stmt.Query{})))
	assert.NotNil(t, factory.NewBrokerExecutor(
		context.TODO(), "db", "sql", nil, nil, nil, nil))
	assert.NotNil(t, factory.NewMetadataStorageExecutor(nil, nil, nil))
	assert.NotNil(t, factory.NewMetadataBrokerExecutor(
		context.TODO(), "db", nil, nil, nil, nil))
}

func TestNewExecutorFactory_NewContext(t *testing.T) {
	factory := NewExecutorFactory()
	assert.NotNil(t, factory.NewStorageExecuteContext(nil, &stmt.Query{}))
}
