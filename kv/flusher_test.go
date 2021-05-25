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

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
)

func TestFlusher_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := NewMockFamily(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		family.EXPECT().newTableBuilder().Return(nil, fmt.Errorf("err")),
	)
	flusher := newStoreFlusher(family)
	err := flusher.Add(uint32(10), []byte("value10"))
	assert.NotNil(t, err)

	builder := table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		family.EXPECT().newTableBuilder().Return(builder, nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(100)),
		family.EXPECT().addPendingOutput(table.FileNumber(100)),
		builder.EXPECT().Add(uint32(10), []byte("value10")).Return(fmt.Errorf("err")),
	)
	flusher = newStoreFlusher(family)
	err = flusher.Add(uint32(10), []byte("value10"))
	assert.NotNil(t, err)

	builder = table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		family.EXPECT().newTableBuilder().Return(builder, nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(100)),
		family.EXPECT().addPendingOutput(table.FileNumber(100)),
		builder.EXPECT().Add(uint32(10), []byte("value10")).Return(nil),
	)
	flusher = newStoreFlusher(family)
	err = flusher.Add(uint32(10), []byte("value10"))
	if err != nil {
		t.Fatal(err)
	}
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
	flusher := newStoreFlusher(family)
	err := flusher.Commit()
	assert.NotNil(t, err)

	// empty commit edit log success
	family = NewMockFamily(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		family.EXPECT().commitEditLog(gomock.Any()).Return(true),
	)
	flusher = newStoreFlusher(family)
	err = flusher.Commit()
	if err != nil {
		t.Fatal(err)
	}

	builder := table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		builder.EXPECT().Close().Return(fmt.Errorf("err")),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		family.EXPECT().removePendingOutput(table.FileNumber(10)),
	)
	flusher = newStoreFlusher(family)
	f := flusher.(*storeFlusher)
	f.builder = builder
	err = flusher.Commit()
	assert.NotNil(t, err)

	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		builder.EXPECT().Close().Return(nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		builder.EXPECT().MinKey().Return(uint32(1)),
		builder.EXPECT().MaxKey().Return(uint32(10)),
		builder.EXPECT().Size().Return(int32(100)),
		family.EXPECT().commitEditLog(gomock.Any()).Return(false),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		family.EXPECT().removePendingOutput(table.FileNumber(10)),
	)
	flusher = newStoreFlusher(family)
	f = flusher.(*storeFlusher)
	f.builder = builder
	err = flusher.Commit()
	assert.NotNil(t, err)

	gomock.InOrder(
		family.EXPECT().ID().Return(version.FamilyID(10)),
		builder.EXPECT().Close().Return(nil),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		builder.EXPECT().MinKey().Return(uint32(1)),
		builder.EXPECT().MaxKey().Return(uint32(10)),
		builder.EXPECT().Size().Return(int32(100)),
		family.EXPECT().commitEditLog(gomock.Any()).Return(true),
		builder.EXPECT().FileNumber().Return(table.FileNumber(10)),
		family.EXPECT().removePendingOutput(table.FileNumber(10)),
	)
	flusher = newStoreFlusher(family)
	f = flusher.(*storeFlusher)
	f.builder = builder
	err = flusher.Commit()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_NopFlusher(t *testing.T) {
	nf := NewNopFlusher()
	assert.Nil(t, nf.Commit())
	assert.Nil(t, nf.Add(1, nil))
	assert.Nil(t, nf.Bytes())

	_ = nf.Add(2, []byte{1, 2, 3})
	assert.NotNil(t, nf.Bytes())
}
