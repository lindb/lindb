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
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/flow"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

const testDBPath = "test_db"

var cfg = MemoryDatabaseCfg{
	TempPath: testDBPath,
}

func TestMemoryDatabase_New(t *testing.T) {
	defer func() {
		mkdirFunc = fileutil.MkDirIfNotExist
	}()

	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, mdINTF)
	err = mdINTF.Close()
	assert.NoError(t, err)

	mkdirFunc = func(path string) error {
		return fmt.Errorf("err")
	}
	mdINTF, err = NewMemoryDatabase(cfg)
	assert.Error(t, err)
	assert.Nil(t, mdINTF)
}

func TestMemoryDatabase_AcquireWrite(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, mdINTF)
	mdINTF.AcquireWrite()
	a := time.After(100 * time.Millisecond)
	go func() {
		<-a
		mdINTF.CompleteWrite()
	}()
	flusher := metricsdata.NewMockFlusher(ctrl)
	flusher.EXPECT().Commit().Return(nil)
	err = mdINTF.FlushFamilyTo(flusher)
	assert.NoError(t, err)
}

func TestMemoryDatabase_Write(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		defer ctrl.Finish()
	}()
	// mock
	mockMetadata := metadb.NewMockMetadata(ctrl)
	mockMetadataDatabase := metadb.NewMockMetadataDatabase(ctrl)
	mockMetadata.EXPECT().MetadataDatabase().Return(mockMetadataDatabase).AnyTimes()
	mockMStore := NewMockmStoreINTF(ctrl)
	tStore := NewMocktStoreINTF(ctrl)
	fStore := NewMockfStoreINTF(ctrl)
	mockMStore.EXPECT().GetOrCreateTStore(uint32(10)).Return(tStore, 10).AnyTimes()
	// build memory-database
	cfg.Metadata = mockMetadata
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	md := mdINTF.(*memoryDatabase)
	assert.Zero(t, md.MemSize())

	// load mock
	md.mStores.Put(uint32(1), mockMStore)
	// case 1: write ok
	gomock.InOrder(
		mockMetadataDatabase.EXPECT().GenFieldID("ns", "test1", field.Name("f1"), field.SumField).Return(field.ID(1), nil),
		tStore.EXPECT().GetFStore(gomock.Any()).Return(fStore, true),
		fStore.EXPECT().Write(gomock.Any(), gomock.Any(), gomock.Any()).Return(10),
		mockMStore.EXPECT().AddField(gomock.Any(), gomock.Any()),
		mockMStore.EXPECT().SetSlot(gomock.Any()),
	)
	err = md.Write("ns", "test1", uint32(1), uint32(10), 1, []*pb.Field{{
		Name:  "f1",
		Type:  pb.FieldType_Sum,
		Value: 10.0,
	}})
	assert.NoError(t, err)
	// case 2: field type unknown
	err = md.Write("ns", "test1", uint32(1), uint32(10), 1, []*pb.Field{{
		Name: "f1",
	}})
	assert.NoError(t, err)
	// case 3: generate field err
	mockMetadataDatabase.EXPECT().GenFieldID("ns", "test1", field.Name("f1-err"), field.SumField).Return(field.ID(0), fmt.Errorf("err"))
	err = md.Write("ns", "test1", uint32(1), uint32(10), 1, []*pb.Field{{
		Name:  "f1-err",
		Type:  pb.FieldType_Sum,
		Value: 10.0,
	}})
	assert.NoError(t, err)
	// case 5: new metric store
	err = md.Write("ns", "test1", uint32(20), uint32(20), 1, []*pb.Field{{
		Name: "f1",
	}})
	assert.NoError(t, err)
	// case 6: create new field store
	gomock.InOrder(
		mockMetadataDatabase.EXPECT().GenFieldID("ns", "test1", field.Name("f4"), field.SumField).Return(field.ID(1), nil),
		tStore.EXPECT().GetFStore(gomock.Any()).Return(nil, false),
		tStore.EXPECT().InsertFStore(gomock.Any()),
		mockMStore.EXPECT().AddField(gomock.Any(), gomock.Any()),
		mockMStore.EXPECT().SetSlot(gomock.Any()),
	)
	err = md.Write("ns", "test1", uint32(1), uint32(10), 15, []*pb.Field{{
		Name:  "f4",
		Type:  pb.FieldType_Sum,
		Value: 10.0,
	}})
	assert.NoError(t, err)
	assert.True(t, md.MemSize() > 0)

	err = md.Close()
	assert.NoError(t, err)
}

func TestMemoryDatabase_Write_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		defer ctrl.Finish()
	}()

	// mock
	mockMetadata := metadb.NewMockMetadata(ctrl)
	mockMetadataDatabase := metadb.NewMockMetadataDatabase(ctrl)
	mockMetadata.EXPECT().MetadataDatabase().Return(mockMetadataDatabase).AnyTimes()
	mockMStore := NewMockmStoreINTF(ctrl)
	tStore := NewMocktStoreINTF(ctrl)
	mockMStore.EXPECT().GetOrCreateTStore(uint32(10)).Return(tStore, 10).AnyTimes()
	// build memory-database
	cfg.Metadata = mockMetadata
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	buf := NewMockDataPointBuffer(ctrl)
	buf.EXPECT().AllocPage().Return(nil, fmt.Errorf("err"))
	md := mdINTF.(*memoryDatabase)
	md.buf = buf

	// load mock
	md.mStores.Put(uint32(1), mockMStore)
	// case 1: write ok
	gomock.InOrder(
		mockMetadataDatabase.EXPECT().GenFieldID("ns", "test1", field.Name("f1"), field.SumField).Return(field.ID(1), nil),
		tStore.EXPECT().GetFStore(gomock.Any()).Return(nil, false),
	)
	err = md.Write("ns", "test1", uint32(1), uint32(10), 1, []*pb.Field{{
		Name:  "f1",
		Type:  pb.FieldType_Sum,
		Value: 10.0,
	}})
	assert.Error(t, err)

	buf.EXPECT().Close().Return(nil)
	err = md.Close()
	assert.NoError(t, err)
}

func TestMemoryDatabase_FlushFamilyTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	md := mdINTF.(*memoryDatabase)
	flusher := metricsdata.NewMockFlusher(ctrl)
	flusher.EXPECT().Commit().Return(nil).AnyTimes()
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	md.mStores.Put(uint32(3333), mockMStore)

	// case 1: flusher ok
	mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(nil)
	err = md.FlushFamilyTo(flusher)
	assert.NoError(t, err)
	// case 2: flusher err
	mockMStore.EXPECT().FlushMetricsDataTo(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = md.FlushFamilyTo(flusher)
	assert.Error(t, err)

	err = md.Close()
	assert.NoError(t, err)
}

func TestMemoryDatabase_Filter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mdINTF, err := NewMemoryDatabase(cfg)
	assert.NoError(t, err)
	md := mdINTF.(*memoryDatabase)

	// case 1: family not found
	rs, err := md.Filter(uint32(3333), nil, timeutil.TimeRange{}, field.Metas{{ID: 1}})
	assert.NoError(t, err)
	assert.Nil(t, rs)
	now := timeutil.Now()
	// case 2: metric store not found
	rs, err = md.Filter(0, nil, timeutil.TimeRange{Start: now - 10, End: now + 20}, field.Metas{{ID: 1}})
	assert.NoError(t, err)
	assert.Nil(t, rs)
	// case 3: filter success
	// mock mStore
	mockMStore := NewMockmStoreINTF(ctrl)
	mockMStore.EXPECT().Filter(gomock.Any(), gomock.Any(), gomock.Any()).Return([]flow.FilterResultSet{}, nil)
	md.mStores.Put(uint32(3333), mockMStore)
	rs, err = md.Filter(uint32(3333), nil, timeutil.TimeRange{Start: now - 10, End: now + 20}, field.Metas{{ID: 1}})
	assert.NoError(t, err)
	assert.NotNil(t, rs)

	err = md.Close()
	assert.NoError(t, err)
}
