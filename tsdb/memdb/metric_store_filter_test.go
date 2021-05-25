package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/series/field"
)

func TestMetricStore_Filter(t *testing.T) {
	metricStore := mockMetricStore()

	// case 1: field not found
	rs, err := metricStore.Filter([]field.ID{1, 2}, nil)
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Nil(t, rs)
	// case 3: series ids not found
	rs, err = metricStore.Filter([]field.ID{1, 20}, roaring.BitmapOf(1, 2))
	assert.Equal(t, constants.ErrNotFound, err)
	assert.Nil(t, rs)
	// case 3: found data
	rs, err = metricStore.Filter([]field.ID{1, 20}, roaring.BitmapOf(1, 100, 200))
	assert.NoError(t, err)
	assert.NotNil(t, rs)
	mrs := rs[0].(*memFilterResultSet)
	assert.EqualValues(t, roaring.BitmapOf(100, 200).ToArray(), mrs.SeriesIDs().ToArray())
	assert.Equal(t,
		field.Metas{{
			ID:   20,
			Type: field.SumField,
		}}, mrs.fields)
	assert.Equal(t, "memory", rs[0].Identifier())
}

func TestMemFilterResultSet_Load(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mStore := mockMetricStore()

	rs, err := mStore.Filter([]field.ID{1, 20}, roaring.BitmapOf(1, 100, 200))
	assert.NoError(t, err)
	// case 1: load data success
	scanner := rs[0].Load(0, roaring.BitmapOf(100, 200).GetContainer(0), []field.ID{20, 30})
	assert.NotNil(t, scanner)
	scanner.Scan(100)
	scanner.Scan(200)
	// case 2: series ids not found
	rs, _ = mStore.Filter([]field.ID{1, 20}, roaring.BitmapOf(1, 100, 200))
	scanner = rs[0].Load(0, roaring.BitmapOf(1, 2).GetContainer(0), []field.ID{20, 30})
	assert.Nil(t, scanner)
	// case 3: high key not exist
	rs, _ = mStore.Filter([]field.ID{1, 20}, roaring.BitmapOf(1, 100, 200))
	scanner = rs[0].Load(10, roaring.BitmapOf(1, 2).GetContainer(0), []field.ID{20, 30})
	assert.Nil(t, scanner)
	// case 4: field not exist
	rs, _ = mStore.Filter([]field.ID{1, 20}, roaring.BitmapOf(1, 100, 200))
	scanner = rs[0].Load(0, roaring.BitmapOf(100, 200).GetContainer(0), []field.ID{21, 30})
	assert.Nil(t, scanner)
}

func mockMetricStore() *metricStore {
	mStore := newMetricStore()
	mStore.AddField(field.ID(10), field.SumField)
	mStore.AddField(field.ID(20), field.SumField)
	mStore.SetSlot(10)
	mStore.SetSlot(20)
	mStore.GetOrCreateTStore(100)
	mStore.GetOrCreateTStore(120)
	mStore.GetOrCreateTStore(200)
	return mStore.(*metricStore)
}
