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

package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/linmetric"
)

func TestNewNetworkStatistics(t *testing.T) {
	assert.NotNil(t, NewConnStatistics(linmetric.BrokerRegistry, "1.1.1.1:8080"))
	assert.NotNil(t, NewGRPCUnaryClientStatistics(linmetric.BrokerRegistry))
	assert.NotNil(t, NewGRPCStreamClientStatistics(linmetric.BrokerRegistry, "t", "s", "m"))
	assert.NotNil(t, NewGRPCUnaryServerStatistics(linmetric.BrokerRegistry))
	assert.NotNil(t, NewGRPCStreamServerStatistics(linmetric.BrokerRegistry, "t", "s", "m"))
}