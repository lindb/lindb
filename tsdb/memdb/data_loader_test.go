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

package memdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/pkg/timeutil"
)

func TestTimeSeriesLoader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	loader := NewTimeSeriesLoader(nil, nil, 0, timeutil.SlotRange{}, nil)
	assert.NotNil(t, loader)

	tsIndex := NewMockTimeSeriesIndex(ctrl)
	tsLoader := &timeSeriesLoader{
		timeSeriesIndex: tsIndex,
	}

	tsIndex.EXPECT().Load(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
	tsLoader.Load(nil)
}
