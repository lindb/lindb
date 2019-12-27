package indexdb

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/kv/version"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore/invertedindex"
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
	defer func() {
		newReader = invertedindex.NewReader
	}()
	mockedDB := mockIndexDatabase(ctrl)
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(0), fmt.Errorf("err"))
	// case1: get tag key err
	g, err := mockedDB.indexDatabase.GetGroupingContext(1, []string{"host"}, series.NewVersion())
	assert.Error(t, err)
	assert.Nil(t, g)
	// case2: get reader err
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(10), nil).AnyTimes()
	mockedDB.snapShot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("rer"))
	g, err = mockedDB.indexDatabase.GetGroupingContext(1, []string{"host"}, series.NewVersion())
	assert.Error(t, err)
	assert.Nil(t, g)
	// case3: index reader walk tag value err
	indexReader := invertedindex.NewMockReader(ctrl)
	newReader = func(readers []table.Reader) invertedindex.Reader {
		return indexReader
	}
	mockedDB.snapShot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{mockedDB.reader}, nil).AnyTimes()
	mockedDB.idGetter.EXPECT().GetTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(10), nil).AnyTimes()
	indexReader.EXPECT().WalkTagValues(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	g, err = mockedDB.indexDatabase.GetGroupingContext(1, []string{"host"}, series.NewVersion())
	assert.Error(t, err)
	assert.Nil(t, g)
	newReader = invertedindex.NewReader
	// case4: unmarshal series ids err
	ipBlock := buildInvertedIndexBlock()
	ipBlock[908] = 99
	mockedDB.reader.EXPECT().Get(gomock.Any()).Return(ipBlock)
	g, err = mockedDB.indexDatabase.GetGroupingContext(1, []string{"host"}, series.Version(1500000000000))
	assert.Error(t, err)
	assert.Nil(t, g)
	// case4: normal
	ipBlock = buildInvertedIndexBlock()
	mockedDB.reader.EXPECT().Get(gomock.Any()).Return(ipBlock)
	g, err = mockedDB.indexDatabase.GetGroupingContext(1, []string{"host"}, series.Version(1500000000000))
	assert.NoError(t, err)
	assert.NotNil(t, g)
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

func TestIndexDatabase_GetSeriesIDsForMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockedDB := mockIndexDatabase(ctrl)
	// case1: get tag key ids err
	mockedDB.idGetter.EXPECT().GetTagKeyIDs(gomock.Any()).Return(nil, fmt.Errorf("err"))
	set, err := mockedDB.indexDatabase.GetSeriesIDsForMetric(0, timeutil.TimeRange{})
	assert.Error(t, err)
	assert.Nil(t, set)
	// case2: snapshot get reader err
	mockedDB.idGetter.EXPECT().GetTagKeyIDs(gomock.Any()).Return([]uint32{1, 2}, nil).AnyTimes()
	mockedDB.snapShot.EXPECT().FindReaders(gomock.Any()).Return(nil, fmt.Errorf("err"))
	set, err = mockedDB.indexDatabase.GetSeriesIDsForMetric(0, timeutil.TimeRange{})
	assert.Error(t, err)
	assert.Nil(t, set)
	// case3: snapshot get reader nil
	mockedDB.snapShot.EXPECT().FindReaders(gomock.Any()).Return(nil, nil).Times(2)
	set, err = mockedDB.indexDatabase.GetSeriesIDsForMetric(0, timeutil.TimeRange{})
	assert.Error(t, err)
	assert.Nil(t, set)
	// case4: get value
	defer func() {
		newReader = invertedindex.NewReader
	}()
	reader := invertedindex.NewMockReader(ctrl)
	newReader = func(readers []table.Reader) invertedindex.Reader {
		return reader
	}
	mockedDB.snapShot.EXPECT().FindReaders(gomock.Any()).Return([]table.Reader{mockedDB.reader}, nil).Times(2)
	reader.EXPECT().GetTagKVEntries(gomock.Any(), gomock.Any()).Return(nil)
	tagKVEntry := invertedindex.NewMockTagKVEntrySetINTF(ctrl)
	tagKVEntry.EXPECT().TagValuesCount().Return(10)
	pIt := invertedindex.NewMockpositionIteratorINTF(ctrl)
	pIt.EXPECT().HasNext().Return(false)
	tagKVEntry.EXPECT().PositionIterator().Return(pIt)
	reader.EXPECT().GetTagKVEntries(gomock.Any(), gomock.Any()).Return(invertedindex.TagKVEntries{tagKVEntry})
	set, err = mockedDB.indexDatabase.GetSeriesIDsForMetric(0, timeutil.TimeRange{})
	assert.Error(t, err)
	assert.Nil(t, set)

}

func buildInvertedIndexBlock() (ipBlock []byte) {
	nopKVFlusher := kv.NewNopFlusher()
	seriesFlusher := invertedindex.NewFlusher(nopKVFlusher)
	// disable auto reset to pick the entrySetBuffer
	/////////////////////////
	// seriesID mapping relation
	/////////////////////////
	ipMapping := map[uint32]string{
		1: "192.168.1.1",
		2: "192.168.1.2",
		3: "192.168.1.3",
		4: "192.168.2.4",
		5: "192.168.2.5",
		6: "192.168.2.6",
		7: "192.168.3.7",
		8: "192.168.3.8",
		9: "192.168.3.9"}
	/////////////////////////
	// flush ip tag, tagID: 21
	/////////////////////////
	for seriesID, ip := range ipMapping {
		for v := series.Version(1500000000000); v < 1800000000000; v += 100000000000 {
			bitmap := roaring.New()
			bitmap.Add(seriesID)
			seriesFlusher.FlushVersion(v, timeutil.TimeRange{
				Start: v.Int64() + 10000*1000, End: v.Int64() + 20000*10000}, bitmap)
		}
		seriesFlusher.FlushTagValue(ip)
	}
	// pick the ipBlock buffer
	_ = seriesFlusher.FlushTagKeyID(21)
	ipBlock = append(ipBlock, nopKVFlusher.Bytes()...)

	return ipBlock
}
