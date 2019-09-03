package parallel

import (
	"testing"

	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/series"
)

func TestResultMerger_Merge(t *testing.T) {
	ch := make(chan *series.TimeSeriesEvent)
	merger := newResultMerger(ch)
	go func() {
		for range ch {
		}
	}()

	merger.Merge(&pb.TaskResponse{TaskID: "taskID"})
}
