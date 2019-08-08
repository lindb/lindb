package storage

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

func TestTaskExecutor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageService := service.NewMockStorageService(ctrl)
	repo := state.NewMockRepository(ctrl)
	exec := NewTaskExecutor(context.TODO(), &models.Node{IP: "1.1.1.1", Port: 5000}, repo, storageService)
	assert.NotNil(t, exec)

	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any()).Return(nil)
	exec.Run()
	time.Sleep(100 * time.Millisecond)
	err := exec.Close()
	if err != nil {
		t.Fatal(err)
	}
}
