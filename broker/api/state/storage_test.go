package state

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
)

func TestStorageAPI(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMachine := broker.NewMockStorageStateMachine(ctrl)
	api := NewStorageAPI(stateMachine)

	storageState := models.NewStorageState()
	storageState.Name = "LinDB_Storage"
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}})

	stateMachine.EXPECT().List().Return([]*models.StorageState{storageState})
	// get success
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/storage/state",
		HandlerFunc:    api.ListStorageCluster,
		ExpectHTTPCode: 200,
		ExpectResponse: []*models.StorageState{storageState},
	})
}
