package indexdb

import (
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/tsdb/indextbl"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_NewIndexDatabase_recover(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mock nameIds reader
	mockReader := table.NewMockReader(ctrl)
	mockSnapShot := kv.NewMockSnapshot(ctrl)
	mockSnapShot.EXPECT().Readers().Return([]table.Reader{mockReader}).AnyTimes()
	// mock read ns ok
	mockReader.EXPECT().Get(gomock.Any()).Return([]byte{1, 2, 3, 4, 5, 6, 7, 8}).AnyTimes()
	db, err := NewIndexDatabase(mockSnapShot, nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, db)
	db, err = NewIndexDatabase(mockSnapShot, nil, nil)
	assert.Nil(t, err)
	assert.NotNil(t, db)
}

func emptyDatabase() *indexDatabase {
	return &indexDatabase{
		tree:               newArtTree(),
		youngMetricNameIDs: make(map[string]uint32),
		youngTagKeyIDs:     make(map[uint32][]tagKeyAndID),
		youngFieldIDs:      make(map[uint32][]fieldIDAndType)}
}

func Test_IndexDatabase_GetMetricID(t *testing.T) {
	db := emptyDatabase()
	// neither in the map or on the tree
	metricID, err := db.GetMetricID("docker")
	assert.Zero(t, metricID)
	assert.NotNil(t, err)
	// in map
	db.youngMetricNameIDs["docker"] = 2
	metricID, err = db.GetMetricID("docker")
	assert.Equal(t, uint32(2), metricID)
	assert.Nil(t, err)
	// on the tree
	db.tree.Insert([]byte("cpu"), uint32(1))
	metricID, err = db.GetMetricID("cpu")
	assert.Equal(t, uint32(1), metricID)
	assert.Nil(t, err)
}

func Test_IndexDatabase_GenMetricID(t *testing.T) {
	db := emptyDatabase()
	// newly created
	db.metricIDSequence = 2
	db.youngMetricNameIDs["docker"] = 2
	assert.Equal(t, uint32(2), db.GenMetricID("docker"))
	// metricID sequence
	assert.Equal(t, uint32(3), db.GenMetricID("cpu"))
	assert.Equal(t, uint32(3), db.GenMetricID("cpu"))
	assert.Equal(t, uint32(4), db.GenMetricID("cpu1"))
	assert.Equal(t, uint32(5), db.GenMetricID("cpu2"))
}

func Test_IndexDatabase_GenTagID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := emptyDatabase()
	mockMetaReader := indextbl.NewMockMetricsMetaReader(ctrl)
	db.metaReader = mockMetaReader

	// data on disk
	mockMetaReader.EXPECT().ReadTagID(gomock.Any(), gomock.Any()).
		Return(uint32(1), true).Times(2)
	assert.Equal(t, uint32(1), db.GenTagID(1, "host1"))
	assert.Equal(t, uint32(1), db.GenTagID(1, "host1"))

	// new tagKey
	mockMetaReader.EXPECT().ReadTagID(gomock.Any(), gomock.Any()).
		Return(uint32(0), false).Times(1).AnyTimes()
	assert.Equal(t, uint32(1), db.GenTagID(1, "host2"))
	assert.Equal(t, uint32(1), db.GenTagID(1, "host2"))

	// newTagKey of same metricID
	assert.Equal(t, uint32(2), db.GenTagID(1, "host3"))
	assert.Equal(t, uint32(2), db.GenTagID(1, "host3"))
	assert.Equal(t, uint32(3), db.GenTagID(1, "host4"))
}

func Test_IndexDatabase_GetFieldID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaReader := indextbl.NewMockMetricsMetaReader(ctrl)
	db := emptyDatabase()
	db.metaReader = metaReader

	// mock not exist
	metaReader.EXPECT().ReadFieldID(gomock.Any(), gomock.Any()).
		Return(uint16(0), field.Type(0), false)
	fieldID, fieldType, err := db.GetFieldID(1, "cpu")
	assert.NotNil(t, err)
	assert.Zero(t, fieldID)
	assert.Zero(t, fieldType)
	// mock ok
	metaReader.EXPECT().ReadFieldID(gomock.Any(), gomock.Any()).
		Return(uint16(1), field.SumField, true)
	fieldID, fieldType, err = db.GetFieldID(1, "cpu")
	assert.Nil(t, err)
	assert.Equal(t, uint16(1), fieldID)
	assert.Equal(t, field.SumField, fieldType)
}

func Test_IndexDatabase_GenFieldID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	metaReader := indextbl.NewMockMetricsMetaReader(ctrl)
	db := emptyDatabase()
	db.metaReader = metaReader

	// - case0: hit memory, type match
	db.youngFieldIDs[1] = append(db.youngFieldIDs[1], fieldIDAndType{
		fieldID: 1, fieldType: field.SumField, fieldName: "sum"})
	fieldID, err := db.GenFieldID(1, "sum", field.SumField)
	assert.Equal(t, uint16(1), fieldID)
	assert.Nil(t, err)
	// - case1: hit memory, type mismatch
	fieldID, err = db.GenFieldID(1, "sum", field.MinField)
	assert.Equal(t, uint16(0), fieldID)
	assert.NotNil(t, err)

	// - case2: hit disk, type match
	metaReader.EXPECT().ReadFieldID(gomock.Any(), gomock.Any()).
		Return(uint16(2), field.MinField, true).Times(2)
	fieldID, err = db.GenFieldID(1, "min", field.MinField)
	assert.Equal(t, uint16(2), fieldID)
	assert.Nil(t, err)
	// - case3: hit disk, type mismatch
	fieldID, err = db.GenFieldID(1, "min", field.MaxField)
	assert.Zero(t, fieldID)
	assert.NotNil(t, err)

	// case4: new field, create ok
	metaReader.EXPECT().ReadMaxFieldID(gomock.Any()).Return(uint16(2)).Times(2)
	metaReader.EXPECT().ReadFieldID(gomock.Any(), gomock.Any()).
		Return(uint16(0), field.Type(0), false).AnyTimes()
	fieldID, err = db.GenFieldID(3, "sum", field.SumField)
	assert.Equal(t, uint16(3), fieldID)
	assert.Nil(t, err)
	fieldID, err = db.GenFieldID(3, "sum1", field.SumField)
	assert.Equal(t, uint16(4), fieldID)
	assert.Nil(t, err)
	// case5: new field, too many fields
	metaReader.EXPECT().ReadMaxFieldID(gomock.Any()).Return(uint16(2000)).Times(1)
	fieldID, err = db.GenFieldID(3, "sum2", field.SumField)
	assert.Zero(t, fieldID)
	assert.NotNil(t, err)
}

func Test_IndexDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	db := emptyDatabase()
	mockSeriesReader := indextbl.NewMockSeriesIndexReader(ctrl)
	db.seriesReader = mockSeriesReader

	mockSeriesReader.EXPECT().GetTagValues(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)
	tagValues, err := db.GetTagValues(1, nil, 1)
	assert.Nil(t, tagValues)
	assert.Nil(t, err)

	mockSeriesReader.EXPECT().FindSeriesIDsByExpr(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)
	set, err := db.FindSeriesIDsByExpr(1, nil, timeutil.TimeRange{})
	assert.Nil(t, set)
	assert.Nil(t, err)

	mockSeriesReader.EXPECT().GetSeriesIDsForTag(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)
	set, err = db.GetSeriesIDsForTag(1, "", timeutil.TimeRange{})
	assert.Nil(t, set)
	assert.Nil(t, err)
}
