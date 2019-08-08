package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/service"
)

func TestCreateShardProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageService := service.NewMockStorageService(ctrl)
	processor := newCreateShardProcessor(storageService)
	assert.Equal(t, 1, processor.Concurrency())
	assert.Equal(t, time.Duration(0), processor.RetryBackOff())
	assert.Equal(t, 0, processor.RetryCount())
	assert.Equal(t, constants.CreateShard, processor.Kind())

	err := processor.Process(context.TODO(), task.Task{Params: []byte{1, 1, 1}})
	assert.NotNil(t, err)
	param := models.CreateShardTask{}
	storageService.EXPECT().CreateShards(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = processor.Process(context.TODO(), task.Task{Params: encoding.JSONMarshal(&param)})
	assert.NotNil(t, err)

	storageService.EXPECT().CreateShards(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = processor.Process(context.TODO(), task.Task{Params: encoding.JSONMarshal(&param)})
	assert.Nil(t, err)
}
