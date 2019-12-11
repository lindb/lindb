package replication

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestReplicatorService_Report(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	srv := NewReplicatorStateReport(models.Node{IP: "1.1.1.1", Port: 9000}, repo)

	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err := srv.Report(&models.BrokerReplicaState{})
	assert.NotNil(t, err)

	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err = srv.Report(&models.BrokerReplicaState{})
	if err != nil {
		t.Fatal(err)
	}
}
