package monitoring

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"

	"github.com/golang/mock/gomock"
)

func Test_NewSystemCollector(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	ctx, cancel := context.WithCancel(context.TODO())

	collector := NewSystemCollector(
		ctx,
		10*time.Millisecond,
		"/tmp",
		repo,
		"",
		models.ActiveNode{},
	)

	go func() {
		time.Sleep(time.Second)
		cancel()
	}()

	// collect system monitoring stat
	collector.Run()
}

func Test_SystemCollector_Collect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	repo := state.NewMockRepository(ctrl)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	collector := NewSystemCollector(
		ctx,
		10*time.Millisecond,
		"/tmp",
		repo,
		"",
		models.ActiveNode{},
	)

	collector.MemoryStatGetter = func() (*models.MemoryStat, error) {
		return nil, fmt.Errorf("error")
	}
	collector.collect()
	collector.MemoryStatGetter = GetMemoryStat

	collector.CPUStatGetter = func() (*models.CPUStat, error) {
		return nil, fmt.Errorf("error")
	}
	collector.collect()
	collector.CPUStatGetter = GetCPUStat

	collector.DiskStatGetter = func(path string) (*models.DiskStat, error) {
		return nil, fmt.Errorf("error")
	}
	collector.collect()
	collector.DiskStatGetter = GetDiskStat
}
