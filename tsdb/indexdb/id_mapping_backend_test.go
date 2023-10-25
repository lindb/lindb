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

package indexdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/common/pkg/fileutil"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/unique"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/series/metric"
)

func TestIDMappingBackend_New(t *testing.T) {
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "make dir failure",
			prepare: func() {
				mkDir = func(path string) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "new id store failure",
			prepare: func() {
				mkDir = func(path string) error {
					return nil
				}
				newIDStoreFn = func(path string) (unique.IDStore, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "create id mapping backend successfully",
			prepare: func() {
				mkDir = func(path string) error {
					return nil
				}
				newIDStoreFn = func(path string) (unique.IDStore, error) {
					return nil, nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				mkDir = fileutil.MkDirIfNotExist
				newIDStoreFn = unique.NewIDStore
			}()

			if tt.prepare != nil {
				tt.prepare()
			}

			backend, err := newIDMappingBackend(t.TempDir())

			if ((err != nil) != tt.wantErr && backend == nil) || (!tt.wantErr && backend == nil) {
				t.Errorf("newIDMappingBackend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIDMappingBackend_Close_sync(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idStore := unique.NewMockIDStore(ctrl)
	backend := &idMappingBackend{
		db: idStore,
	}

	idStore.EXPECT().Flush().Return(fmt.Errorf("err"))
	assert.Error(t, backend.sync())
	idStore.EXPECT().Close().Return(nil)
	assert.NoError(t, backend.Close())
}

func TestIDMappingBackend_getSeriesID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idStore := unique.NewMockIDStore(ctrl)
	cases := []struct {
		name    string
		prepare func(idStore *unique.MockIDStore)
		out     struct {
			seriesID uint32
			err      error
		}
	}{
		{
			name: "get series id failure from backend",
			prepare: func(idStore *unique.MockIDStore) {
				idStore.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			out: struct {
				seriesID uint32
				err      error
			}{
				seriesID: series.EmptySeriesID,
				err:      fmt.Errorf("err"),
			},
		},
		{
			name: "series not found from backend",
			prepare: func(idStore *unique.MockIDStore) {
				idStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
			},
			out: struct {
				seriesID uint32
				err      error
			}{
				seriesID: series.EmptySeriesID,
				err: fmt.Errorf("%w, metricID: %d, tagsHash: %d",
					constants.ErrSeriesIDNotFound, 2, 123),
			},
		},
		{
			name: "get series successfully",
			prepare: func(idStore *unique.MockIDStore) {
				idStore.EXPECT().Get(gomock.Any()).Return([]byte{2, 0, 0, 0}, true, nil)
			},
			out: struct {
				seriesID uint32
				err      error
			}{seriesID: uint32(2), err: nil},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			backend := &idMappingBackend{
				db: idStore,
			}
			if tt.prepare != nil {
				tt.prepare(idStore)
			}

			seriesID, err := backend.getSeriesID(metric.ID(2), 123)
			assert.Equal(t, tt.out.seriesID, seriesID)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestIDMappingBackend_genSeriesID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idStore := unique.NewMockIDStore(ctrl)
	backend := &idMappingBackend{db: idStore}

	idStore.EXPECT().Put(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	assert.Error(t, backend.genSeriesID(metric.ID(2), 123, 2))
}

func TestIDMappingBackend_loadMetricIDMapping(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	idStore := unique.NewMockIDStore(ctrl)

	cases := []struct {
		name    string
		prepare func(idStore *unique.MockIDStore)
		wantErr bool
	}{
		{
			name: "load mapping failure from backend",
			prepare: func(idStore *unique.MockIDStore) {
				idStore.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "load mapping not exist",
			prepare: func(idStore *unique.MockIDStore) {
				idStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
				idStore.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "load mapping, init sequence, persist failure",
			prepare: func(idStore *unique.MockIDStore) {
				idStore.EXPECT().Get(gomock.Any()).Return([]byte{2, 0, 0, 0}, true, nil)
				idStore.EXPECT().Put([]byte{2, 0, 0, 0}, gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "load mapping, init sequence, persist successfully",
			prepare: func(idStore *unique.MockIDStore) {
				idStore.EXPECT().Get(gomock.Any()).Return([]byte{2, 0, 0, 0}, true, nil)
				idStore.EXPECT().Put([]byte{2, 0, 0, 0}, gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			backend := &idMappingBackend{db: idStore}
			if tt.prepare != nil {
				tt.prepare(idStore)
			}
			mapping, err := backend.loadMetricIDMapping(metric.ID(2))
			if ((err != nil) != tt.wantErr && backend == nil) || (!tt.wantErr && mapping == nil) {
				t.Errorf("loadMetricIDMapping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestIdMappingBackend_saveSeriesSequence(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := unique.NewMockIDStore(ctrl)
	backend := idMappingBackend{
		db: store,
	}
	store.EXPECT().Put([]byte{1, 0, 0, 0}, []byte{10, 0, 0, 0}).Return(nil)
	assert.NoError(t, backend.saveSeriesSequence(metric.ID(1), uint32(10)))
}
