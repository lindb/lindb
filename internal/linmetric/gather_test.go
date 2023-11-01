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

package linmetric

import (
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/series/tag"
)

func Test_Gather(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	observer := NewMockObserver(ctrl)
	observer.EXPECT().Observe()
	gather := BrokerRegistry.NewGather(
		WithGlobalKeyValueOption(tag.TagsFromMap(map[string]string{
			"host": "alpha",
			"ip":   "1.1.1.1",
		})),
		WithReadRuntimeOption(observer),
		WithNamespaceOption("default-ns"),
	)
	_, _ = gather.Gather()
}
