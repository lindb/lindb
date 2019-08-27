package indexdb

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/field"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb/tblstore"

	"github.com/golang/mock/gomock"
	art "github.com/plar/go-adaptive-radix-tree"
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
	db := NewIndexDatabase(nil, nil)
	// once test
	_ = NewIndexDatabase(nil, nil)

	nameIDReader := tblstore.NewMetricsNameIDReader(mockSnapShot)
	err := db.Recover(nameIDReader)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// once test
	err = db.Recover(nameIDReader)
	assert.Nil(t, err)
	assert.NotNil(t, db)

	// mock unmarshal failure
	mockNameIDReader := tblstore.NewMockMetricsNameIDReader(ctrl)
	mockNameIDReader.EXPECT().ReadMetricNS(gomock.Any()).
		Return([][]byte{{1, 3}}, uint32(1), uint32(1), true).AnyTimes()
	assert.NotNil(t, db.Recover(mockNameIDReader))
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
	mockMetaReader := tblstore.NewMockMetricsMetaReader(ctrl)
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

	metaReader := tblstore.NewMockMetricsMetaReader(ctrl)
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

	metaReader := tblstore.NewMockMetricsMetaReader(ctrl)
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
	mockForwardIdxReader := tblstore.NewMockForwardIndexReader(ctrl)
	mockInvertedIdxReader := tblstore.NewMockInvertedIndexReader(ctrl)
	db.forwardIndexReader = mockForwardIdxReader
	db.invertedIndexReader = mockInvertedIdxReader

	mockForwardIdxReader.EXPECT().GetTagValues(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)
	tagValues, err := db.GetTagValues(1, nil, 1)
	assert.Nil(t, tagValues)
	assert.Nil(t, err)

	db.youngTagKeyIDs = map[uint32][]tagKeyAndID{1: {{tagKey: "host", tagKeyID: 2}, {tagKey: "zone", tagKeyID: 3}}}
	mockInvertedIdxReader.EXPECT().FindSeriesIDsByExprForTagID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil, nil)
	set, err := db.FindSeriesIDsByExpr(1, &stmt.EqualsExpr{Key: "host", Value: "dev"}, timeutil.TimeRange{})
	assert.Nil(t, set)
	assert.Nil(t, err)

	mockInvertedIdxReader.EXPECT().GetSeriesIDsForTagID(gomock.Any(), gomock.Any()).
		Return(nil, nil)
	set, err = db.GetSeriesIDsForTag(1, "zone", timeutil.TimeRange{})
	assert.Nil(t, set)
	assert.Nil(t, err)
}

func Test_IndexDatabase_FlushNameIDsTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := emptyDatabase()
	mockKVFlusher := kv.NewMockFlusher(ctrl)
	mockKVFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	mockFlusher := tblstore.NewMetricsNameIDFlusher(mockKVFlusher)
	assert.Nil(t, db.FlushNameIDsTo(mockFlusher))

	db.youngMetricNameIDs["1"] = 1
	db.youngMetricNameIDs["2"] = 2
	db.metricIDSequence = 10
	db.tagKeyIDSequence = 15

	assert.Nil(t, db.FlushNameIDsTo(mockFlusher))
	assert.Equal(t, 2, db.tree.Size())
}

func Test_IndexDatabase_FlushMetricsMetaTo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := emptyDatabase()
	assert.Nil(t, db.FlushMetricsMetaTo(nil))

	set := func() {
		db.youngTagKeyIDs = map[uint32][]tagKeyAndID{
			1: {{tagKey: "11", tagKeyID: 11},
				{tagKey: "12", tagKeyID: 12}},
			2: {{tagKey: "22", tagKeyID: 22},
				{tagKey: "23", tagKeyID: 23}}}
		db.youngFieldIDs = map[uint32][]fieldIDAndType{
			2: {{fieldID: 22, fieldType: field.SumField},
				{fieldID: 23, fieldType: field.MaxField}},
			3: {{fieldID: 33, fieldType: field.MinField},
				{fieldID: 34, fieldType: field.SumField}}}
	}
	mockKVFlusher := kv.NewMockFlusher(ctrl)
	set()
	mockKVFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(nil).Times(3)
	mockMetaFlusher := tblstore.NewMetricsMetaFlusher(mockKVFlusher)
	assert.Nil(t, db.FlushMetricsMetaTo(mockMetaFlusher))

	// map empty
	mockKVFlusher.EXPECT().Add(gomock.Any(), gomock.Any()).Return(fmt.Errorf("error")).Times(1)
	assert.Nil(t, db.FlushMetricsMetaTo(mockMetaFlusher))
	// flush with error
	set()
	assert.NotNil(t, db.FlushMetricsMetaTo(mockMetaFlusher))
}

func Test_IndexDatabase_SuggestMetrics(t *testing.T) {
	db := emptyDatabase()
	for i := 10000; i < 30000; i++ {
		db.tree.Insert(art.Key(strconv.Itoa(i)), i)
	}
	// invalid limit
	assert.Len(t, db.SuggestMetrics("1", -1), 0)
	// limit exceeds the limit
	assert.Len(t, db.SuggestMetrics("2", 20000), 10000)
	// smaller than limit
	assert.Len(t, db.SuggestMetrics("2000", 5000), 11)
}

func Test_IndexDatabase_SuggestTagKeysValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	db := emptyDatabase()
	// invalid limit
	assert.Nil(t, db.SuggestTagKeys("", "", 0))
	assert.Nil(t, db.SuggestTagValues("", "", "", 0))

	mockMetaReader := tblstore.NewMockMetricsMetaReader(ctrl)
	db.metaReader = mockMetaReader

	// SuggestTagKeys
	assert.Len(t, db.SuggestTagKeys("inexistent-metric", "", 100), 0)
	db.tree.Insert(art.Key("m1"), uint32(1))
	mockMetaReader.EXPECT().SuggestTagKeys(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]string{"key1"}).AnyTimes()
	assert.Len(t, db.SuggestTagKeys("m1", "", 100), 1)

	// SuggestTagValues
	mockInvertedReader := tblstore.NewMockInvertedIndexReader(ctrl)
	mockInvertedReader.EXPECT().SuggestTagValues(gomock.Any(), gomock.Any(), gomock.Any()).
		Return([]string{"v1"}).AnyTimes()
	db.invertedIndexReader = mockInvertedReader
	// metricID inexistent
	assert.Nil(t, db.SuggestTagValues("inexistent-metric", "", "", 10000000))
	gomock.InOrder(
		mockMetaReader.EXPECT().ReadTagID(uint32(1), "1").Return(uint32(0), false),
		mockMetaReader.EXPECT().ReadTagID(uint32(1), "2").Return(uint32(1), true))
	// tagID inexistent
	assert.Nil(t, db.SuggestTagValues("m1", "1", "", 10))
	// tagID exist
	assert.Len(t, db.SuggestTagValues("m1", "2", "", 10), 1)
}
