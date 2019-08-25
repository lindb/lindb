package kv

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv/table"
)

func TestFlusher_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	family := NewMockFamily(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(10),
		family.EXPECT().newTableBuilder().Return(nil, fmt.Errorf("err")),
	)
	flusher := newStoreFlusher(family)
	err := flusher.Add(uint32(10), []byte("value10"))
	assert.NotNil(t, err)

	builder := table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(10),
		family.EXPECT().newTableBuilder().Return(builder, nil),
		builder.EXPECT().FileNumber().Return(int64(100)),
		family.EXPECT().addPendingOutput(int64(100)),
		builder.EXPECT().Add(uint32(10), []byte("value10")).Return(fmt.Errorf("err")),
	)
	flusher = newStoreFlusher(family)
	err = flusher.Add(uint32(10), []byte("value10"))
	assert.NotNil(t, err)

	builder = table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(10),
		family.EXPECT().newTableBuilder().Return(builder, nil),
		builder.EXPECT().FileNumber().Return(int64(100)),
		family.EXPECT().addPendingOutput(int64(100)),
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
		family.EXPECT().ID().Return(10),
		family.EXPECT().commitEditLog(gomock.Any()).Return(false),
	)
	flusher := newStoreFlusher(family)
	err := flusher.Commit()
	assert.NotNil(t, err)

	// empty commit edit log success
	family = NewMockFamily(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(10),
		family.EXPECT().commitEditLog(gomock.Any()).Return(true),
	)
	flusher = newStoreFlusher(family)
	err = flusher.Commit()
	if err != nil {
		t.Fatal(err)
	}

	builder := table.NewMockBuilder(ctrl)
	gomock.InOrder(
		family.EXPECT().ID().Return(10),
		builder.EXPECT().Close().Return(fmt.Errorf("err")),
		builder.EXPECT().FileNumber().Return(int64(10)),
		family.EXPECT().removePendingOutput(int64(10)),
	)
	flusher = newStoreFlusher(family)
	f := flusher.(*storeFlusher)
	f.builder = builder
	err = flusher.Commit()
	assert.NotNil(t, err)

	gomock.InOrder(
		family.EXPECT().ID().Return(10),
		builder.EXPECT().Close().Return(nil),
		builder.EXPECT().FileNumber().Return(int64(10)),
		builder.EXPECT().MinKey().Return(uint32(1)),
		builder.EXPECT().MaxKey().Return(uint32(10)),
		builder.EXPECT().Size().Return(int32(100)),
		family.EXPECT().commitEditLog(gomock.Any()).Return(false),
		builder.EXPECT().FileNumber().Return(int64(10)),
		family.EXPECT().removePendingOutput(int64(10)),
	)
	flusher = newStoreFlusher(family)
	f = flusher.(*storeFlusher)
	f.builder = builder
	err = flusher.Commit()
	assert.NotNil(t, err)

	gomock.InOrder(
		family.EXPECT().ID().Return(10),
		builder.EXPECT().Close().Return(nil),
		builder.EXPECT().FileNumber().Return(int64(10)),
		builder.EXPECT().MinKey().Return(uint32(1)),
		builder.EXPECT().MaxKey().Return(uint32(10)),
		builder.EXPECT().Size().Return(int32(100)),
		family.EXPECT().commitEditLog(gomock.Any()).Return(true),
		builder.EXPECT().FileNumber().Return(int64(10)),
		family.EXPECT().removePendingOutput(int64(10)),
	)
	flusher = newStoreFlusher(family)
	f = flusher.(*storeFlusher)
	f.builder = builder
	err = flusher.Commit()
	if err != nil {
		t.Fatal(err)
	}
}
