package parallel

import (
	"context"
	"testing"

	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
)

func TestResultMerger_Merge(t *testing.T) {
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(context.TODO(), ch)
	go func() {
		for range ch {
		}
	}()

	merger.merge(&pb.TaskResponse{TaskID: "taskID"})

	merger.close()
}
