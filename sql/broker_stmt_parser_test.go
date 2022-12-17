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

package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/stmt"
)

func TestShowBrokers(t *testing.T) {
	q, err := Parse("show brokers")
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Broker{Type: stmt.BrokerOpShow}, q)
}

func TestCreateBroker(t *testing.T) {
	cfg := `{\"config\":{\"namespace\":\"test\",\"timeout\":10,\"dialTimeout\":10,\"leaseTTL\":10,\"endpoints\":[\"http://localhost:2379\"]}}`
	sql := `create broker ` + cfg
	q, err := Parse(sql)
	assert.NoError(t, err)
	assert.Equal(t, &stmt.Broker{
		Type:  stmt.BrokerOpCreate,
		Value: `{"config":{"namespace":"test","timeout":10,"dialTimeout":10,"leaseTTL":10,"endpoints":["http://localhost:2379"]}}`,
	}, q)
}
