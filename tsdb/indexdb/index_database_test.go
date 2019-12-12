package indexdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/metadb"
)

////////////////////////////////
// helper methods
////////////////////////////////

type mockedIndexDatabase struct {
	indexDatabase *indexDatabase
	family        *kv.MockFamily
	snapShot      *version.MockSnapshot
	reader        *table.MockReader
	idGetter      *metadb.MockIDGetter
}

func (db *mockedIndexDatabase) WithFindReadersError() {
	db.snapShot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("error"))
}

func (db *mockedIndexDatabase) WithFindReadersOK() {
	db.snapShot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{db.reader}, nil)
}

func (db *mockedIndexDatabase) WithFindReadersEmpty() {
	db.snapShot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil)
}

func mockIndexDatabase(ctrl *gomock.Controller) *mockedIndexDatabase {
	mockReader := table.NewMockReader(ctrl)

	mockSnapShot := version.NewMockSnapshot(ctrl)
	mockSnapShot.EXPECT().Close().Return().AnyTimes()

	mockFamily := kv.NewMockFamily(ctrl)
	mockFamily.EXPECT().GetSnapshot().Return(mockSnapShot).AnyTimes()

	mockIDGetter := metadb.NewMockIDGetter(ctrl)
	return &mockedIndexDatabase{
		indexDatabase: NewIndexDatabase(mockIDGetter, mockFamily).(*indexDatabase),
		family:        mockFamily,
		snapShot:      mockSnapShot,
		reader:        mockReader,
		idGetter:      mockIDGetter}
}

func Test_IndexDatabase_GetGroupingContext(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedDB := mockIndexDatabase(ctrl)
	g, err := mockedDB.indexDatabase.GetGroupingContext(1, nil, series.NewVersion())
	assert.NoError(t, err)
	assert.Nil(t, g)
}

func Test_IndexDatabase_SuggestTagValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedDB := mockIndexDatabase(ctrl)

	// case1: invalid limit
	assert.Nil(t, mockedDB.indexDatabase.SuggestTagValues("", "", "", 0))
	// case2: limit>max, GetMetricID failed
	mockedDB.idGetter.EXPECT().GetMetricID(gomock.Any()).Return(uint32(0), fmt.Errorf("error"))
	assert.Nil(t, mockedDB.indexDatabase.SuggestTagValues("", "", "", 100000000))
	// case3: GetTagKeyID failed
	mockedDB.idGetter.EXPECT().GetMetricID(gomock.Any()).Return(uint32(1), nil)
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("error"))
	assert.Nil(t, mockedDB.indexDatabase.SuggestTagValues("", "", "", 10000))
	// case4: snapshot FindReaders error
	mockedDB.WithFindReadersError()
	mockedDB.idGetter.EXPECT().GetMetricID(gomock.Any()).Return(uint32(1), nil)
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
	assert.Nil(t, mockedDB.indexDatabase.SuggestTagValues("", "", "", 10000))
	// case4: snapshot FindReaders ok
	mockedDB.WithFindReadersOK()
	mockedDB.reader.EXPECT().Get(gomock.Any()).Return(nil)
	mockedDB.idGetter.EXPECT().GetMetricID(gomock.Any()).Return(uint32(1), nil)
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
	assert.Nil(t, mockedDB.indexDatabase.SuggestTagValues("", "", "", 10000))
}

type mockTagKey struct {
	key string
}

func (k mockTagKey) TagKey() string {
	return k.key
}

func Test_IndexDatabase_FindSeriesIDsByExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedDB := mockIndexDatabase(ctrl)

	// case1: GetTagKeyID failed
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("error"))
	set, err := mockedDB.indexDatabase.FindSeriesIDsByExpr(0, &mockTagKey{key: ""}, timeutil.TimeRange{})
	assert.Nil(t, set)
	assert.NotNil(t, err)
	// case2: snapshot FindReaders error
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
	mockedDB.WithFindReadersError()
	_, err = mockedDB.indexDatabase.FindSeriesIDsByExpr(0, &mockTagKey{key: ""}, timeutil.TimeRange{})
	assert.NotNil(t, err)
	// case3: snapshot FindReaders ok
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
	mockedDB.WithFindReadersOK()
	mockedDB.reader.EXPECT().Get(gomock.Any()).Return(nil)
	_, err = mockedDB.indexDatabase.FindSeriesIDsByExpr(0, &mockTagKey{key: ""}, timeutil.TimeRange{})
	assert.NotNil(t, err)

	// case4: snapshot FindReaders is nil
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
	mockedDB.WithFindReadersEmpty()
	_, err = mockedDB.indexDatabase.FindSeriesIDsByExpr(0, &mockTagKey{key: ""}, timeutil.TimeRange{})
	assert.NotNil(t, err)
}

func Test_IndexDatabase_GetSeriesIDsForTag(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedDB := mockIndexDatabase(ctrl)

	// case1: GetTagKeyID failed
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("error"))
	set, err := mockedDB.indexDatabase.GetSeriesIDsForTag(0, "", timeutil.TimeRange{})
	assert.Nil(t, set)
	assert.NotNil(t, err)
	// case2: snapshot FindReaders error
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
	mockedDB.WithFindReadersError()
	_, err = mockedDB.indexDatabase.GetSeriesIDsForTag(0, "", timeutil.TimeRange{})
	assert.NotNil(t, err)
	// case3: snapshot FindReaders ok
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(1), nil)
	mockedDB.WithFindReadersOK()
	mockedDB.reader.EXPECT().Get(gomock.Any()).Return(nil)
	_, err = mockedDB.indexDatabase.GetSeriesIDsForTag(0, "", timeutil.TimeRange{})
	assert.NotNil(t, err)
}
