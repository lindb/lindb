package service

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestStorageSateService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	storageState := models.NewStorageState()
	storageState.Name = "LinDB_Storage"
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}})

	srv := NewStorageStateService(repo)

	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err := srv.Save("Test_LinDB", storageState)
	if err != nil {
		t.Fatal(err)
	}

	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = srv.Save("Test_LinDB", storageState)
	assert.NotNil(t, err)

	data, _ := json.Marshal(&storageState)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(data, nil)
	storageState1, _ := srv.Get("Test_LinDB")
	assert.Equal(t, storageState, storageState1)

	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte{1, 1, 3}, nil)
	storageState1, err = srv.Get("Test_LinDB")
	assert.NotNil(t, err)
	assert.Nil(t, storageState1)

	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, state.ErrNotExist)
	_, err = srv.Get("Test_LinDB_Not_Exist")
	assert.Equal(t, state.ErrNotExist, err)
}
