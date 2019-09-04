package query

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
	"github.com/lindb/lindb/tsdb/series"
)

func TestStorageExecute_validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().Name().Return("mock_tsdb").AnyTimes()
	query := &stmt.Query{Interval: timeutil.OneSecond}

	// query shards is empty
	exec := newStorageExecutor(engine, nil, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	// shards of engine is empty
	engine.EXPECT().NumOfShards().Return(0)
	exec = newStorageExecutor(engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	// num. of shard not match
	engine.EXPECT().NumOfShards().Return(2)
	exec = newStorageExecutor(engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	engine.EXPECT().NumOfShards().Return(3).AnyTimes()
	engine.EXPECT().GetShard(gomock.Any()).Return(nil).MaxTimes(3)
	exec = newStorageExecutor(engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	// normal case
	query, _ = sql.Parse("select f from cpu")
	engine1 := MockTSDBEngine(ctrl)
	exec = newStorageExecutor(engine1, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.Nil(t, exec.Error())
}

func TestStorageExecute_Simple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var scanners []series.DataFamilyScanner
	for i := 0; i < 3; i++ {
		seriesData := MockSumFieldSeries(ctrl, 10, 1, map[int]interface{}{
			5:  5.5,
			15: 5.5,
			17: 5.5,
			16: 5.5,
			56: 5.5,
		})
		itr := series.NewMockVersionIterator(ctrl)
		itr.EXPECT().Close()
		itr.EXPECT().HasNext().Return(true)
		itr.EXPECT().Next().Return(seriesData)

		// finish scanner
		itr.EXPECT().HasNext().Return(false)

		scanner := series.NewMockDataFamilyScanner(ctrl)
		scanner.EXPECT().Scan(gomock.Any()).Return(itr)

		scanners = append(scanners, scanner)
	}

	engine := MockTSDBEngine(ctrl, scanners...)

	// normal case
	query, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := newStorageExecutor(engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.Nil(t, exec.Error())

	execImpl := exec.(*storageExecutor)
	// mock scanner return nil
	mockScanner1 := series.NewMockDataFamilyScanner(ctrl)
	mockScanner1.EXPECT().Scan(gomock.Any()).Return(nil).Times(1)
	execImpl.familyLevelSearch(mockScanner1, nil)
	// mock scanner return iterator with nil ts
	mockScanner2 := series.NewMockDataFamilyScanner(ctrl)
	mockItr := series.NewMockVersionIterator(ctrl)
	mockItr.EXPECT().Close().Return(nil)
	mockItr.EXPECT().HasNext().Return(true)
	mockItr.EXPECT().Next().Return(nil)
	mockScanner2.EXPECT().Scan(gomock.Any()).Return(mockItr)
	execImpl.familyLevelSearch(mockScanner2, nil)
	// check shards error
	execImpl.shardIDs = nil
	assert.NotNil(t, execImpl.checkShards())
}
