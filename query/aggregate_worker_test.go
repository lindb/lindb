package query

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/aggregation"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series"
)

func TestAggregateWorker_emit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	resultCh := make(chan *series.TimeSeriesEvent)
	wait := sync.WaitGroup{}
	go func() {
		for range resultCh {
			wait.Done()
		}
	}()
	wait.Add(1)
	agg := createAggWorker(10000, &timeutil.TimeRange{
		Start: 10,
		End:   20,
	}, map[string]*aggregation.AggregatorSpec{"f1": {}}, resultCh)
	agg.emit(&series.FieldEvent{
		Interval:        10000,
		FamilyStartTime: 10,
	})
	agg.sendResult(nil)
	agg.close()
}
