package monitoring

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/pkg/state"
)

func Test_heartbeatReport_Report(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	reporter := NewHeartbeatReporter(context.TODO(), repo, "/tmp")
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	reporter.Report("test")

	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	reporter.Report("test")
}
