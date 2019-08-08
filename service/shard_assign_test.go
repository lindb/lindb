package service

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestShardAssignService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	srv := NewShardAssignService(repo)

	shardAssign1 := models.NewShardAssignment()
	shardAssign1.AddReplica(1, 1)
	shardAssign1.AddReplica(1, 2)
	shardAssign1.AddReplica(1, 3)
	shardAssign1.AddReplica(2, 2)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	_ = srv.Save("db1", shardAssign1)

	shardAssign2 := models.NewShardAssignment()
	shardAssign2.AddReplica(1, 1)
	shardAssign2.AddReplica(2, 2)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	_ = srv.Save("db2", shardAssign2)

	data1, _ := json.Marshal(shardAssign1)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(data1, nil)
	shardAssign11, _ := srv.Get("db1")
	assert.Equal(t, *shardAssign1, *shardAssign11)

	data2, _ := json.Marshal(shardAssign2)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(data2, nil)
	shardAssign22, _ := srv.Get("db2")
	assert.Equal(t, *shardAssign2, *shardAssign22)

	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, state.ErrNotExist)
	_, err := srv.Get("not_exist")
	assert.Equal(t, state.ErrNotExist, err)

	// unmarshal error
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte{1, 3, 34}, nil)
	_, err = srv.Get("not_exist")
	assert.NotNil(t, err)
}
