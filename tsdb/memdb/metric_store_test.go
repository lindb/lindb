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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	gomock "go.uber.org/mock/gomock"

	"github.com/lindb/lindb/pkg/imap"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestMetricStore_SetTimestamp(t *testing.T) {
	mStoreInterface := newMetricStore()
	mStoreInterface.SetSlot(10)
	slotRange := mStoreInterface.GetSlotRange()
	assert.Equal(t, uint16(10), slotRange.Start)
	assert.Equal(t, uint16(10), slotRange.End)
	mStoreInterface.SetSlot(5)
	slotRange = mStoreInterface.GetSlotRange()
	assert.Equal(t, uint16(5), slotRange.Start)
	assert.Equal(t, uint16(10), slotRange.End)
	mStoreInterface.SetSlot(50)
	slotRange = mStoreInterface.GetSlotRange()
	assert.Equal(t, uint16(5), slotRange.Start)
	assert.Equal(t, uint16(50), slotRange.End)
}

func TestMetricStore_Flush_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	flusher := metricsdata.NewMockFlusher(ctrl)
	flusher.EXPECT().PrepareMetric(gomock.Any(), gomock.Any()).AnyTimes()
	ids := imap.NewIntMap[uint32]()
	ids.Put(1, 1)
	ms := &metricStore{
		slotRange: &timeutil.SlotRange{Start: 0},
		ids:       ids,
	}
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "no field",
			prepare: func() {
				ms.fields = nil
			},
		},
		{
			name: "flush field error",
			prepare: func() {
				ms.fields = append(ms.fields, field.Meta{ID: 1, Persisted: true})
			},
			wantErr: true,
		},
	}
	for i := range cases {
		tt := cases[i]
		t.Run(tt.name, func(t *testing.T) {
			tt.prepare()
			err := ms.FlushMetricsDataTo(flusher,
				&flushContext{},
				func(memSeriesID uint32, fields field.Metas) error {
					return fmt.Errorf("err")
				},
			)
			if (err != nil) != tt.wantErr {
				t.Fatal(tt.name)
			}
		})
	}
}
