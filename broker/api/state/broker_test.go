package state

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/eleme/lindb/coordinator/broker"
	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/timeutil"
)

func TestBrokerAPI_ListBrokerNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	node := models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 2080}, OnlineTime: timeutil.Now()}
	nodes := []models.ActiveNode{node}

	stateMachine := broker.NewMockNodeStateMachine(ctrl)
	stateMachine.EXPECT().GetActiveNodes().Return(nodes)
	api := NewBrokerAPI(stateMachine)

	// get success
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state",
		HandlerFunc:    api.ListBrokerNodes,
		ExpectHTTPCode: 200,
		ExpectResponse: nodes,
	})
}
