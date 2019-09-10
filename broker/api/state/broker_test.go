package state

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
)

func TestBrokerAPI_ListBrokerNodes(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	node := models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 2080}, OnlineTime: timeutil.Now()}
	nodes := []models.ActiveNode{node}

	stateMachine := broker.NewMockNodeStateMachine(ctrl)
	stateMachine.EXPECT().GetActiveNodesByType(models.NodeTypeRPC).Return(nodes)
	stateMachine.EXPECT().GetActiveNodesByType(models.NodeTypeTCP).Return(nodes)
	api := NewBrokerAPI(stateMachine)

	// get success
	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state",
		HandlerFunc:    api.ListBrokerNodes,
		ExpectHTTPCode: 200,
		ExpectResponse: nodes,
	})

	mock.DoRequest(t, &mock.HTTPHandler{
		Method:         http.MethodGet,
		URL:            "/broker/state?type=tcp",
		HandlerFunc:    api.ListBrokerNodes,
		ExpectHTTPCode: 200,
		ExpectResponse: nodes,
	})
}
