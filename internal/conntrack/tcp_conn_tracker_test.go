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

package conntrack

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

func TestTrackedConn(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	conn := NewMockConn(ctrl)
	connTracker := &TrackedConn{
		Conn:       conn,
		statistics: metrics.NewConnStatistics(linmetric.BrokerRegistry, "1.1.1.1:8080"),
	}

	conn.EXPECT().Write(gomock.Any()).Return(0, fmt.Errorf("err"))
	conn.EXPECT().Read(gomock.Any()).Return(0, fmt.Errorf("err"))
	conn.EXPECT().Close().Return(fmt.Errorf("err"))

	n, err := connTracker.Read(nil)
	assert.Error(t, err)
	assert.Zero(t, n)
	n, err = connTracker.Write(nil)
	assert.Error(t, err)
	assert.Zero(t, n)
	err = connTracker.Close()
	assert.Error(t, err)
}
