package query

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/sql"
	"github.com/lindb/lindb/sql/stmt"
	"github.com/lindb/lindb/tsdb"
)

func TestStorageExecute_validation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	engine := tsdb.NewMockEngine(ctrl)
	engine.EXPECT().Name().Return("mock_tsdb").AnyTimes()
	query := &stmt.Query{Interval: timeutil.OneSecond}

	// query shards is empty
	exec := NewStorageExecutor(engine, nil, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	// shards of engine is empty
	engine.EXPECT().NumOfShards().Return(0)
	exec = NewStorageExecutor(engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	// num. of shard not match
	engine.EXPECT().NumOfShards().Return(2)
	exec = NewStorageExecutor(engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	engine.EXPECT().NumOfShards().Return(3).AnyTimes()
	engine.EXPECT().GetShard(gomock.Any()).Return(nil).MaxTimes(3)
	exec = NewStorageExecutor(engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	// normal case
	query, _ = sql.Parse("select f from cpu")
	engine1 := MockTSDBEngine(ctrl)
	exec = NewStorageExecutor(engine1, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.Nil(t, exec.Error())
}

func TestStorageExecute_Simple(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var families []tsdb.DataFamily
	for i := 0; i < 3; i++ {
		seriesData := MockSumFieldSeries(ctrl, 10, 1, map[int]interface{}{
			5:  5.5,
			15: 5.5,
			17: 5.5,
			16: 5.5,
			56: 5.5,
		})
		scanner := tsdb.NewMockScanner(ctrl)
		scanner.EXPECT().Close()
		scanner.EXPECT().HasNext().Return(true)
		scanner.EXPECT().Next().Return(seriesData)

		// finish scanner
		scanner.EXPECT().HasNext().Return(false)

		family := tsdb.NewMockDataFamily(ctrl)
		family.EXPECT().Scan(gomock.Any()).Return(scanner)

		families = append(families, family)
	}

	engine := MockTSDBEngine(ctrl, families...)

	// normal case
	query, _ := sql.Parse("select f from cpu where time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	exec := NewStorageExecutor(engine, []int32{1, 2, 3}, query)
	_ = exec.Execute()
	assert.Nil(t, exec.Error())
}
