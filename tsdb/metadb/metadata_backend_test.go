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

package metadb

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/unique"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/metric"
	"github.com/lindb/lindb/series/tag"
)

func TestNewMetadataBackend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "make storage path fail",
			prepare: func() {
				mkDirFn = func(path string) error {
					return fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "new id store fail",
			prepare: func() {
				mkDirFn = func(path string) error {
					return nil
				}
				newIDStoreFn = func(path string) (unique.IDStore, error) {
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "new some id store fail",
			prepare: func() {
				mkDirFn = func(path string) error {
					return nil
				}
				store := unique.NewMockIDStore(ctrl)
				// close fail
				store.EXPECT().Close().Return(fmt.Errorf("err"))

				newIDStoreFn = func(path string) (unique.IDStore, error) {
					if strings.Contains(path, NamespaceDB) {
						return store, nil
					}
					return nil, fmt.Errorf("err")
				}
			},
			wantErr: true,
		},
		{
			name: "init seq fail",
			prepare: func() {
				mkDirFn = func(path string) error {
					return nil
				}
				store := unique.NewMockIDStore(ctrl)
				store.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
				// close fail
				store.EXPECT().Close().Return(fmt.Errorf("err")).MaxTimes(4)

				newIDStoreFn = func(path string) (unique.IDStore, error) {
					return store, nil
				}
			},
			wantErr: true,
		},
		{
			name: "persist cached sequence value failure",
			prepare: func() {
				mkDirFn = func(path string) error {
					return nil
				}
				store := unique.NewMockIDStore(ctrl)
				store.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3, 4}, true, nil)
				store.EXPECT().Put(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
				// close fail
				store.EXPECT().Close().Return(fmt.Errorf("err")).MaxTimes(4)
				newIDStoreFn = func(path string) (unique.IDStore, error) {
					return store, nil
				}
			},
			wantErr: true,
		},
		{
			name: "new backend successfully, init seq from backend storage",
			prepare: func() {
				mkDirFn = func(path string) error {
					return nil
				}
				store := unique.NewMockIDStore(ctrl)
				store.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3, 4}, true, nil).MaxTimes(3)
				store.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)
				newIDStoreFn = func(path string) (unique.IDStore, error) {
					return store, nil
				}
			},
			wantErr: false,
		},
		{
			name: "new backend successfully, init seq = 0",
			prepare: func() {
				mkDirFn = func(path string) error {
					return nil
				}
				store := unique.NewMockIDStore(ctrl)
				store.EXPECT().Get(gomock.Any()).Return(nil, false, nil).MaxTimes(3)
				store.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil).MaxTimes(3)

				newIDStoreFn = func(path string) (unique.IDStore, error) {
					return store, nil
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				mkDirFn = fileutil.MkDirIfNotExist
				newIDStoreFn = unique.NewIDStore
			}()

			if tt.prepare != nil {
				tt.prepare()
			}

			backend, err := newMetadataBackend(t.TempDir())

			if ((err != nil) != tt.wantErr && backend == nil) || (!tt.wantErr && backend == nil) {
				t.Errorf("newMetadataBackend() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetadataBackend_suggestNamespace(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var cases = []struct {
		name    string
		prepare func(idStore *unique.MockIDStore)
		out     struct {
			ns  []string
			err error
		}
	}{
		{
			name: "suggest failure",
			prepare: func(idStore *unique.MockIDStore) {
				idStore.EXPECT().IterKeys(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("err"))
			},
			out: struct {
				ns  []string
				err error
			}{
				ns:  nil,
				err: fmt.Errorf("err"),
			},
		},
		{
			name: "suggest successfully",
			prepare: func(idStore *unique.MockIDStore) {
				idStore.EXPECT().IterKeys(gomock.Any(), gomock.Any()).
					Return([][]byte{[]byte("test"), []byte("ns"), namespaceIDSequenceKey}, nil)
			},
			out: struct {
				ns  []string
				err error
			}{
				ns:  []string{"test", "ns"},
				err: nil,
			},
		},
	}
	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			nsStore := unique.NewMockIDStore(ctrl)
			backend := &metadataBackend{
				namespace: nsStore,
			}
			if tt.prepare != nil {
				tt.prepare(nsStore)
			}

			ns, err := backend.suggestNamespace("ns", 10)

			assert.Equal(t, tt.out.ns, ns)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataBackend_suggestMetricName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name    string
		prepare func(ns, metric *unique.MockIDStore)
		out     struct {
			metricNames []string
			err         error
		}
	}{
		{
			name: "get ns id failure",
			prepare: func(ns, metric *unique.MockIDStore) {
				ns.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			out: struct {
				metricNames []string
				err         error
			}{
				metricNames: nil,
				err:         fmt.Errorf("err"),
			},
		},
		{
			name: "ns id not found",
			prepare: func(ns, metric *unique.MockIDStore) {
				ns.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
			},
			out: struct {
				metricNames []string
				err         error
			}{
				metricNames: nil,
				err:         nil,
			},
		},
		{
			name: "suggest metric name failure",
			prepare: func(ns, metric *unique.MockIDStore) {
				ns.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3, 4}, true, nil)
				metric.EXPECT().IterKeys(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
			},
			out: struct {
				metricNames []string
				err         error
			}{
				metricNames: nil,
				err:         fmt.Errorf("err"),
			},
		},
		{
			name: "suggest metric name successfully",
			prepare: func(ns, metric *unique.MockIDStore) {
				ns.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3, 4}, true, nil)
				metric.EXPECT().IterKeys(gomock.Any(), gomock.Any()).
					Return([][]byte{[]byte("1234name")}, nil)
			},
			out: struct {
				metricNames []string
				err         error
			}{
				metricNames: []string{"name"},
				err:         nil,
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			nsStore := unique.NewMockIDStore(ctrl)
			metricStore := unique.NewMockIDStore(ctrl)
			backend := &metadataBackend{
				namespace: nsStore,
				metric:    metricStore,
			}
			if tt.prepare != nil {
				tt.prepare(nsStore, metricStore)
			}

			metricNames, err := backend.suggestMetricName("ns", "name", 10)
			assert.Equal(t, tt.out.metricNames, metricNames)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataBackend_getMetricID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name    string
		prepare func(ns, metric *unique.MockIDStore)
		out     struct {
			metricID metric.ID
			err      error
		}
	}{
		{
			name: "get ns id failure",
			prepare: func(ns, metric *unique.MockIDStore) {
				ns.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			out: struct {
				metricID metric.ID
				err      error
			}{
				metricID: metric.ID(0),
				err:      fmt.Errorf("err"),
			},
		},
		{
			name: "ns id not found",
			prepare: func(ns, metric *unique.MockIDStore) {
				ns.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
			},
			out: struct {
				metricID metric.ID
				err      error
			}{
				metricID: metric.ID(0),
				err:      fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, "name"),
			},
		},
		{
			name: "get metric id failure",
			prepare: func(ns, metric *unique.MockIDStore) {
				ns.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3, 4}, true, nil)
				metric.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			out: struct {
				metricID metric.ID
				err      error
			}{
				metricID: metric.ID(0),
				err:      fmt.Errorf("err"),
			},
		},
		{
			name: "metric id not found",
			prepare: func(ns, metric *unique.MockIDStore) {
				ns.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3, 4}, true, nil)
				metric.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
			},
			out: struct {
				metricID metric.ID
				err      error
			}{
				metricID: metric.ID(0),
				err:      fmt.Errorf("%w, metric: %s", constants.ErrMetricIDNotFound, "name"),
			},
		},
		{
			name: "get metric id successfully",
			prepare: func(ns, metric *unique.MockIDStore) {
				ns.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3, 4}, true, nil)
				metric.EXPECT().Get(gomock.Any()).Return([]byte{2, 0, 0, 0}, true, nil)
			},
			out: struct {
				metricID metric.ID
				err      error
			}{
				metricID: metric.ID(2),
				err:      nil,
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			nsStore := unique.NewMockIDStore(ctrl)
			metricStore := unique.NewMockIDStore(ctrl)
			backend := &metadataBackend{
				namespace: nsStore,
				metric:    metricStore,
			}
			if tt.prepare != nil {
				tt.prepare(nsStore, metricStore)
			}

			metricID, err := backend.getMetricID("ns", "name")
			assert.Equal(t, tt.out.metricID, metricID)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataBackend_saveTagKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := unique.NewMockIDStore(ctrl)
	sequence := unique.NewMockSequence(ctrl)

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "get next sequence value failure",
			prepare: func() {
				sequence.EXPECT().HasNext().Return(false)
				sequence.EXPECT().Current().Return(uint32(10))
				store.EXPECT().Put(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "store tag meta failure",
			prepare: func() {
				sequence.EXPECT().HasNext().Return(true)
				sequence.EXPECT().Next().Return(uint32(10))
				store.EXPECT().Merge(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "store tag meta successfully",
			prepare: func() {
				sequence.EXPECT().HasNext().Return(true)
				sequence.EXPECT().Next().Return(uint32(10))
				store.EXPECT().Merge(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "store tag meta successfully, cache next sequence",
			prepare: func() {
				sequence.EXPECT().HasNext().Return(false)
				sequence.EXPECT().Current().Return(uint32(10)).MaxTimes(2)
				store.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
				sequence.EXPECT().Limit(gomock.Any())
				sequence.EXPECT().Next().Return(uint32(100))
				store.EXPECT().Merge(gomock.Any(), gomock.Any()).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			backend := &metadataBackend{
				tagKey:           store,
				tagKeyIDSequence: sequence,
			}
			if tt.prepare != nil {
				tt.prepare()
			}

			_, err := backend.saveTagKey(metric.ID(2), "key")
			if (err != nil) != tt.wantErr {
				t.Errorf("saveTagKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetadataBackend_getAllTagKeys(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name    string
		prepare func(tagKey *unique.MockIDStore)
		out     struct {
			tags tag.Metas
			err  error
		}
	}{
		{
			name: "get tag keys failure",
			prepare: func(tagKey *unique.MockIDStore) {
				tagKey.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			out: struct {
				tags tag.Metas
				err  error
			}{
				tags: nil,
				err:  fmt.Errorf("err"),
			},
		},
		{
			name: "tag keys not found",
			prepare: func(tagKey *unique.MockIDStore) {
				tagKey.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
			},
			out: struct {
				tags tag.Metas
				err  error
			}{
				tags: nil,
				err:  nil,
			},
		},
		{
			name: "get tag keys ok, but unmarshal tag data failure",
			prepare: func(tagKey *unique.MockIDStore) {
				tagKey.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3}, true, nil)
			},
			out: struct {
				tags tag.Metas
				err  error
			}{
				tags: nil,
				err:  fmt.Errorf("EOF"),
			},
		},
		{
			name: "get tag keys successfully",
			prepare: func(tagKey *unique.MockIDStore) {
				var buf []byte
				tags := tag.Metas{
					{
						Key: "test100",
						ID:  100,
					},
					{
						Key: "test10",
						ID:  10,
					},
				}
				for _, tag1 := range tags {
					data, err := tag1.MarshalBinary()
					assert.NoError(t, err)
					buf = append(buf, data...)
				}

				tagKey.EXPECT().Get(gomock.Any()).Return(buf, true, nil)
			},
			out: struct {
				tags tag.Metas
				err  error
			}{
				tags: tag.Metas{
					{
						Key: "test100",
						ID:  100,
					},
					{
						Key: "test10",
						ID:  10,
					},
				},
				err: nil,
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tagKeyStore := unique.NewMockIDStore(ctrl)
			backend := &metadataBackend{
				tagKey: tagKeyStore,
			}
			if tt.prepare != nil {
				tt.prepare(tagKeyStore)
			}

			tags, err := backend.getAllTagKeys(metric.ID(2))
			assert.Equal(t, tt.out.tags, tags)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataBackend_saveField(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := unique.NewMockIDStore(ctrl)
	backend := &metadataBackend{
		field: store,
	}
	f := field.Meta{
		ID:   1,
		Type: field.SumField,
		Name: "field",
	}
	v, err := f.MarshalBinary()
	assert.NoError(t, err)
	store.EXPECT().Merge([]byte{2, 0, 0, 0}, v).Return(nil)
	err = backend.saveField(metric.ID(2), f)
	assert.NoError(t, err)
}

func TestMetadataBackend_getAllFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cases := []struct {
		name    string
		prepare func(field *unique.MockIDStore)
		out     struct {
			fields field.Metas
			max    field.ID
			err    error
		}
	}{
		{
			name: "get fields failure",
			prepare: func(field *unique.MockIDStore) {
				field.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			out: struct {
				fields field.Metas
				max    field.ID
				err    error
			}{
				fields: nil,
				max:    0,
				err:    fmt.Errorf("err"),
			},
		},
		{
			name: "fields not found",
			prepare: func(field *unique.MockIDStore) {
				field.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
			},
			out: struct {
				fields field.Metas
				max    field.ID
				err    error
			}{
				fields: nil,
				max:    field.ID(0),
				err:    nil,
			},
		},
		{
			name: "get fields ok, but unmarshal fields data failure",
			prepare: func(field *unique.MockIDStore) {
				field.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3}, true, nil)
			},
			out: struct {
				fields field.Metas
				max    field.ID
				err    error
			}{
				fields: nil,
				max:    field.ID(0),
				err:    fmt.Errorf("EOF"),
			},
		},
		{
			name: "get fields successfully",
			prepare: func(fieldStore *unique.MockIDStore) {
				var buf []byte
				fields := field.Metas{
					{
						Name: "field100",
						Type: field.SumField,
						ID:   100,
					},
					{
						Name: "field10",
						Type: field.MaxField,
						ID:   10,
					},
				}
				for _, field1 := range fields {
					data, err := field1.MarshalBinary()
					assert.NoError(t, err)
					buf = append(buf, data...)
				}

				fieldStore.EXPECT().Get(gomock.Any()).Return(buf, true, nil)
			},
			out: struct {
				fields field.Metas
				max    field.ID
				err    error
			}{
				fields: field.Metas{
					{
						Name: "field100",
						Type: field.SumField,
						ID:   100,
					},
					{
						Name: "field10",
						Type: field.MaxField,
						ID:   10,
					},
				},
				max: field.ID(100),
				err: nil,
			},
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			fieldStore := unique.NewMockIDStore(ctrl)
			backend := &metadataBackend{
				field: fieldStore,
			}
			if tt.prepare != nil {
				tt.prepare(fieldStore)
			}

			fields, max, err := backend.getAllFields(metric.ID(2))
			assert.Equal(t, tt.out.fields, fields)
			assert.Equal(t, tt.out.max, max)
			assert.Equal(t, tt.out.err, err)
		})
	}
}

func TestMetadataBackend_getMetricMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tagStore := unique.NewMockIDStore(ctrl)
	fieldStore := unique.NewMockIDStore(ctrl)

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "get fields failure",
			prepare: func() {
				fieldStore.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "get tag keys failure",
			prepare: func() {
				fieldStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
				tagStore.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "get metric metadata successfully",
			prepare: func() {
				fieldStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
				tagStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			backend := &metadataBackend{
				tagKey: tagStore,
				field:  fieldStore,
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			meta, err := backend.getMetricMetadata(metric.ID(10))
			if ((err != nil) != tt.wantErr && meta == nil) || (!tt.wantErr && meta == nil) {
				t.Errorf("getMetricMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetadataBackend_getOrCreateMetricMetadata(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nsSequence := unique.NewMockSequence(ctrl)
	metricSequence := unique.NewMockSequence(ctrl)
	nsStore := unique.NewMockIDStore(ctrl)
	metricStore := unique.NewMockIDStore(ctrl)
	fieldStore := unique.NewMockIDStore(ctrl)
	tagStore := unique.NewMockIDStore(ctrl)

	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "get ns failure",
			prepare: func() {
				nsStore.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "gen ns sequence failure",
			prepare: func() {
				nsStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
				nsSequence.EXPECT().HasNext().Return(false)
				nsSequence.EXPECT().Current().Return(uint32(10))
				nsStore.EXPECT().Put(namespaceIDSequenceKey, gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "store ns meta failure",
			prepare: func() {
				nsStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
				nsSequence.EXPECT().HasNext().Return(true)
				nsSequence.EXPECT().Next().Return(uint32(10))
				nsStore.EXPECT().Put([]byte("ns"), []byte{10, 0, 0, 0}).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "store ns meta successfully, but get metric failure",
			prepare: func() {
				nsStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
				nsSequence.EXPECT().HasNext().Return(true)
				nsSequence.EXPECT().Next().Return(uint32(10))
				nsStore.EXPECT().Put([]byte("ns"), []byte{10, 0, 0, 0}).Return(nil)
				metricStore.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{

			name: "get metric failure",
			prepare: func() {
				nsStore.EXPECT().Get(gomock.Any()).Return([]byte{1, 0, 0, 0}, true, nil)
				metricStore.EXPECT().Get(gomock.Any()).Return(nil, false, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{

			name: "get metric sequence failure",
			prepare: func() {
				nsStore.EXPECT().Get(gomock.Any()).Return([]byte{1, 0, 0, 0}, true, nil)
				metricStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
				metricSequence.EXPECT().HasNext().Return(false)
				metricSequence.EXPECT().Current().Return(uint32(10))
				metricStore.EXPECT().Put(metricIDSequenceKey, gomock.Any()).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "save metric meta failure",
			prepare: func() {
				nsStore.EXPECT().Get(gomock.Any()).Return([]byte{1, 0, 0, 0}, true, nil)
				metricStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
				metricSequence.EXPECT().HasNext().Return(true)
				metricSequence.EXPECT().Next().Return(uint32(10))
				var key []byte
				key = append(key, []byte{1, 0, 0, 0}...)
				key = append(key, []byte("metric")...)
				metricStore.EXPECT().Put(key, []byte{10, 0, 0, 0}).Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "save metric meta successfully",
			prepare: func() {
				nsStore.EXPECT().Get(gomock.Any()).Return([]byte{1, 0, 0, 0}, true, nil)
				metricStore.EXPECT().Get(gomock.Any()).Return(nil, false, nil)
				metricSequence.EXPECT().HasNext().Return(true)
				metricSequence.EXPECT().Next().Return(uint32(10))
				var key []byte
				key = append(key, []byte{1, 0, 0, 0}...)
				key = append(key, []byte("metric")...)
				metricStore.EXPECT().Put(key, []byte{10, 0, 0, 0}).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "load metric meta successfully",
			prepare: func() {
				nsStore.EXPECT().Get([]byte("ns")).Return([]byte{1, 0, 0, 0}, true, nil)
				var key []byte
				key = append(key, []byte{1, 0, 0, 0}...)
				key = append(key, []byte("metric")...)
				mID := []byte{10, 0, 0, 0}
				metricStore.EXPECT().Get(key).Return(mID, true, nil)
				fieldStore.EXPECT().Get(mID).Return(nil, false, nil)
				tagStore.EXPECT().Get(mID).Return(nil, false, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			backend := &metadataBackend{
				namespaceIDSequence: nsSequence,
				metricIDSequence:    metricSequence,
				namespace:           nsStore,
				metric:              metricStore,
				field:               fieldStore,
				tagKey:              tagStore,
			}
			if tt.prepare != nil {
				tt.prepare()
			}

			meta, err := backend.getOrCreateMetricMetadata("ns", "metric")
			if ((err != nil) != tt.wantErr && meta == nil) || (!tt.wantErr && meta == nil) {
				t.Errorf("getOrCreateMetricMetadata() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMetadataBackend_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sequence := unique.NewMockSequence(ctrl)
	store := unique.NewMockIDStore(ctrl)
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "all components release failure",
			prepare: func() {
				sequence.EXPECT().Current().Return(uint32(10))
				store.EXPECT().Put(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
				store.EXPECT().Flush().Return(fmt.Errorf("err"))
				store.EXPECT().Close().Return(fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "all components release successfully",
			prepare: func() {
				sequence.EXPECT().Current().Return(uint32(10))
				store.EXPECT().Put(gomock.Any(), gomock.Any()).Return(nil)
				store.EXPECT().Flush().Return(nil)
				store.EXPECT().Close().Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			backend := &metadataBackend{
				sequences: []sequenceItem{{
					sequence: sequence,
					store:    store,
					key:      tagKeyIDSequenceKey,
				}},
				dbs: map[string]unique.IDStore{
					TagKeyDB: store,
				},
			}
			if tt.prepare != nil {
				tt.prepare()
			}
			err := backend.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
