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

package kv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestFlusher_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := NewMockFamily(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		family.EXPECT().newTableBuilder().Return(nil, fmt.Errorf("err")),
	)
	flusher := newStoreFlusher(family, func() {})
	defer flusher.Release()
	err := flusher.Add(uint32(10), []byte("value10"))
	assert.Error(t, err)

	builder := table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		family.EXPECT().newTableBuilder().Return(builder, nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(100)),
		builder.EXPECT().Add(uint32(10), []byte("value10")).Return(fmt.Errorf("err")),
	)
	flusher = newStoreFlusher(family, func() {})
	defer flusher.Release()
	err = flusher.Add(uint32(10), []byte("value10"))
	assert.Error(t, err)

	builder = table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		family.EXPECT().newTableBuilder().Return(builder, nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(100)),
		builder.EXPECT().Add(uint32(10), []byte("value10")).Return(nil),
	)
	flusher = newStoreFlusher(family, func() {})
	defer flusher.Release()
	err = flusher.Add(uint32(10), []byte("value10"))
	assert.NoError(t, err)
}

func TestStoreFlusher_Commit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// empty but commit edit log fail
	family := NewMockFamily(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		family.EXPECT().commitEditLog(gomock.Any()).Return(false),
	)
	flusher := newStoreFlusher(family, func() {})
	defer flusher.Release()
	err := flusher.Commit()
	assert.Error(t, err)

	// empty commit edit log success
	family = NewMockFamily(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		family.EXPECT().commitEditLog(gomock.Any()).Return(true),
	)
	flusher = newStoreFlusher(family, func() {})
	defer flusher.Release()
	flusher.Sequence(1, 10)
	flusher.Sequence(2, 20)
	err = flusher.Commit()
	assert.NoError(t, err)

	builder := table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		builder.EXPECT().Close().Return(fmt.Errorf("err")),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		family.EXPECT().removePendingOutput(table.FileNumber(10)),
	)
	flusher = newStoreFlusher(family, func() {})
	defer flusher.Release()
	f := flusher.(*storeFlusher)
	f.builder = builder
	err = flusher.Commit()
	assert.Error(t, err)

	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		builder.EXPECT().Close().Return(nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		builder.EXPECT().MinKey().Return(uint32(1)),
		builder.EXPECT().MaxKey().Return(uint32(10)),
		builder.EXPECT().Size().Return(uint32(100)),
		family.EXPECT().commitEditLog(gomock.Any()).Return(false),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		family.EXPECT().removePendingOutput(table.FileNumber(10)),
	)
	flusher = newStoreFlusher(family, func() {})
	defer flusher.Release()
	f = flusher.(*storeFlusher)
	f.builder = builder
	err = flusher.Commit()
	assert.Error(t, err)

	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		builder.EXPECT().Close().Return(nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		builder.EXPECT().MinKey().Return(uint32(1)),
		builder.EXPECT().MaxKey().Return(uint32(10)),
		builder.EXPECT().Size().Return(uint32(100)),
		family.EXPECT().commitEditLog(gomock.Any()).Return(true),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		family.EXPECT().removePendingOutput(table.FileNumber(10)),
	)
	flusher = newStoreFlusher(family, func() {})
	defer flusher.Release()
	store := NewMockStore(ctrl)
	family.EXPECT().getStore().Return(store)
	opt := StoreOption{Rollup: []timeutil.Interval{10, 20}}
	store.EXPECT().Option().Return(opt)
	flusher1 := flusher.(*storeFlusher)
	flusher1.outputs = []table.FileNumber{1, 2, 3}
	f = flusher.(*storeFlusher)
	f.builder = builder
	err = flusher.Commit()
	assert.NoError(t, err)
}

func TestStoreFlusher_StreamWriter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := NewMockFamily(ctrl)
	family.EXPECT().ID().Return(version.FamilyID(10)).AnyTimes()
	flusher := newStoreFlusher(family, func() {})
	cases := []struct {
		name    string
		prepare func()
		wantErr bool
	}{
		{
			name: "create stream writer failure",
			prepare: func() {
				family.EXPECT().newTableBuilder().Return(nil, fmt.Errorf("err"))
			},
			wantErr: true,
		},
		{
			name: "create stream writer successfully",
			prepare: func() {
				builder := table.NewMockBuilder(ctrl)
				builder.EXPECT().FileNumber().Return(table.FileNumber(10))
				builder.EXPECT().StreamWriter().Return(&nopStreamWriter{})
				family.EXPECT().newTableBuilder().Return(builder, nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare()
			}
			sw, err := flusher.StreamWriter()
			if ((err != nil) != tt.wantErr && sw == nil) || (!tt.wantErr && sw == nil) {
				t.Errorf("StreamWriter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_NopFlusher(t *testing.T) {
	nf := NewNopFlusher()
	nf.Sequence(1, 10)
	assert.Nil(t, nf.Commit())
	assert.Nil(t, nf.Add(1, nil))
	assert.Nil(t, nf.Bytes())

	_ = nf.Add(2, []byte{1, 2, 3})
	assert.NotNil(t, nf.Bytes())

	writer, _ := nf.StreamWriter()
	writer.Prepare(0)
	assert.Equal(t, uint32(0), writer.CRC32CheckSum())
	assert.Zero(t, writer.Size())
	_, _ = writer.Write([]byte{1, 2, 3})
	assert.Equal(t, uint32(3), writer.Size())
	_, _ = writer.Write([]byte{4, 5, 6})
	assert.Equal(t, uint32(6), writer.Size())
	_, _ = writer.Write(nil)
	assert.Equal(t, uint32(6), writer.Size())
	_ = writer.Commit()
	assert.Equal(t, uint32(2180413220), writer.CRC32CheckSum())

	writer.Prepare(2)
	assert.Zero(t, writer.Size())
	assert.Equal(t, uint32(0), writer.CRC32CheckSum())
	_, _ = writer.Write([]byte{1, 2})
	_, _ = writer.Write([]byte{3, 4})
	_, _ = writer.Write([]byte{5, 6})
	assert.Equal(t, uint32(6), writer.Size())
	assert.Equal(t, uint32(2180413220), writer.CRC32CheckSum())
	nf.Release()

	nopSW := &nopStreamWriter{}
	nopSW.Release()
}
