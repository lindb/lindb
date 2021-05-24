package memdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/tblstore/metricsdata"
)

func TestTimeSeriesStore_GetOrCreateFStore(t *testing.T) {
	tStore := newTimeSeriesStore()
	f, ok := tStore.GetFStore(1, 10)
	assert.Nil(t, f)
	assert.False(t, ok)
	tStore.InsertFStore(newFieldStore(make([]byte, pageSize), 1, 10))
	// get field store
	f, ok = tStore.GetFStore(1, 10)
	assert.NotNil(t, f)
	assert.True(t, ok)
	// field store not exist
	f, ok = tStore.GetFStore(1, 100)
	assert.Nil(t, f)
	assert.False(t, ok)
	for i := 1; i < 10; i++ {
		tStore.InsertFStore(newFieldStore(make([]byte, pageSize), familyID(1*i), field.ID(10*i)))
		tStore.InsertFStore(newFieldStore(make([]byte, pageSize), 1, 10))
		f, ok = tStore.GetFStore(1, 10)
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
		flusher.EXPECT().GetFieldMetas().Return(field.Metas{{ID: 1}, {ID: 2}, {ID: 3}}),
		fStore.EXPECT().GetFieldID().Return(field.ID(2)),
		flusher.EXPECT().FlushField(nil),
		fStore.EXPECT().GetFieldID().Return(field.ID(2)),
		fStore.EXPECT().FlushFieldTo(gomock.Any(), gomock.Any(), gomock.Any()),
		flusher.EXPECT().FlushField(nil),
	)
	tStore.FlushSeriesTo(flusher, flushContext{familyID: 20})
}

func TestTimeSeriesStore_scan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tStoreInterface := newTimeSeriesStore()
	tStore := tStoreInterface.(*timeSeriesStore)

	for i := 0; i < 10; i++ {
		fStore := newFieldStore(make([]byte, pageSize), familyID(i), field.ID(i*10))
		tStore.InsertFStore(fStore)
		fStore.Write(field.SumField, uint16(i), 10.1)
	}

	// case 1: family time not match
	tStore.scan([]FieldKey{buildFieldKey(11, 10)},
		field.Metas{{
			ID:   10,
			Type: field.SumField,
		}})
	// case 2: field id not match
	tStore.scan([]FieldKey{buildFieldKey(5, 100)},
		field.Metas{{
			ID:   200,
			Type: field.SumField,
		}})
	// case 3: field id not match
	tStore.scan([]FieldKey{buildFieldKey(5, 80)},
		field.Metas{{
			ID:   80,
			Type: field.SumField,
		}})
	// case 4: field key not match
	tStore.scan([]FieldKey{buildFieldKey(5, 80)},
		field.Metas{{
			ID:   80,
			Type: field.SumField,
		}})
	// case 4: match one field
	tStore.scan([]FieldKey{buildFieldKey(5, 50)},
		field.Metas{{
			ID:   50,
			Type: field.SumField,
		}})
	// case 4: match two fields
	tStore.scan([]FieldKey{buildFieldKey(5, 50), buildFieldKey(8, 80)},
		field.Metas{{
			ID:   50,
			Type: field.SumField,
		}, {
			ID:   80,
			Type: field.SumField,
		}})
}
