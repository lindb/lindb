package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestTimeSeriesStore_GetOrCreateFStore(t *testing.T) {
	tStore := newTimeSeriesStore()
	f, ok := tStore.GetFStore(1, 10, 11)
	assert.Nil(t, f)
	assert.False(t, ok)
	tStore.InsertFStore(newFieldStore(make([]byte, pageSize), 1, 10, 11))
	// get field store
	f, ok = tStore.GetFStore(1, 10, 11)
	assert.NotNil(t, f)
	assert.True(t, ok)
	// field store not exist
	f, ok = tStore.GetFStore(1, 10, 10)
	assert.Nil(t, f)
	assert.False(t, ok)
	for i := 1; i < 10; i++ {
		tStore.InsertFStore(newFieldStore(make([]byte, pageSize), familyID(1*i), field.ID(10*i), field.PrimitiveID(11*i)))
		tStore.InsertFStore(newFieldStore(make([]byte, pageSize), 1, 10, 11))
		f, ok = tStore.GetFStore(1, 10, 11)
		assert.NotNil(t, f)
		assert.True(t, ok)
	}
}

func TestTimeSeriesStore_FlushSeriesTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	flusher := metricsdata.NewMockFlusher(ctrl)
	tStore := newTimeSeriesStore()
	s := tStore.(*timeSeriesStore)
	fStore := NewMockfStoreINTF(ctrl)
	fStore.EXPECT().GetFamilyID().Return(familyID(10))
	s.InsertFStore(fStore)

	// case 1: not match family id
	tStore.FlushSeriesTo(flusher, flushContext{familyID: 20})

	// case 2: flush by family id
	gomock.InOrder(
		fStore.EXPECT().GetFamilyID().Return(familyID(20)),
		fStore.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any()),
	)
	tStore.FlushSeriesTo(flusher, flushContext{familyID: 20})
}

func TestTimeSeriesStore_scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore()
	tStore := tStoreInterface.(*timeSeriesStore)

	for i := 0; i < 10; i++ {
		fStore := newFieldStore(make([]byte, pageSize), familyID(i), field.ID(i*10), field.PrimitiveID(i))
		tStore.InsertFStore(fStore)
		fStore.Write(field.SumField, uint16(i), 10.1)
	}
	pAgg := aggregation.NewMockPrimitiveAggregator(ctrl)
	pAgg.EXPECT().FieldID().Return(field.PrimitiveID(5)).Times(2)

	// case 1: family time not match
	tStore.scan(&memScanContext{
		fieldAggs: []*fieldAggregator{
			newFieldAggregator(familyID(11), field.Meta{
				ID:   10,
				Type: field.SumField,
			}, pAgg),
		},
	})
	// case 2: field id not match
	tStore.scan(&memScanContext{
		fieldAggs: []*fieldAggregator{
			newFieldAggregator(familyID(5), field.Meta{
				ID:   200,
				Type: field.SumField,
			}, pAgg),
		},
	})
	// case 3: primitive field id not match
	pAgg.EXPECT().FieldID().Return(field.PrimitiveID(50))
	tStore.scan(&memScanContext{
		fieldAggs: []*fieldAggregator{
			newFieldAggregator(familyID(5), field.Meta{
				ID:   80,
				Type: field.SumField,
			}, pAgg),
		},
	})
	// case 4: field key not match
	pAgg.EXPECT().FieldID().Return(field.PrimitiveID(5))
	tStore.scan(&memScanContext{
		fieldAggs: []*fieldAggregator{
			newFieldAggregator(familyID(5), field.Meta{
				ID:   80,
				Type: field.SumField,
			}, pAgg),
		},
	})
	// case 4: match one field
	pAgg.EXPECT().FieldID().Return(field.PrimitiveID(5))
	pAgg.EXPECT().Aggregate(5, 10.1)
	tStore.scan(&memScanContext{
		fieldAggs: []*fieldAggregator{
			newFieldAggregator(familyID(5), field.Meta{
				ID:   50,
				Type: field.SumField,
			}, pAgg),
		},
	})
	// case 4: match two fields
	pAgg2 := aggregation.NewMockPrimitiveAggregator(ctrl)
	gomock.InOrder(
		pAgg.EXPECT().FieldID().Return(field.PrimitiveID(5)),
		pAgg2.EXPECT().FieldID().Return(field.PrimitiveID(8)),
		pAgg.EXPECT().Aggregate(5, 10.1),
		pAgg2.EXPECT().Aggregate(8, 10.1),
	)
	tStore.scan(&memScanContext{
		fieldAggs: []*fieldAggregator{
			newFieldAggregator(familyID(5), field.Meta{
				ID:   50,
				Type: field.SumField,
			}, pAgg),
			newFieldAggregator(familyID(8), field.Meta{
				ID:   80,
				Type: field.SumField,
			}, pAgg2),
		},
	})
}
