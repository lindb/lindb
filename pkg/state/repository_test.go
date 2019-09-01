package state

import (
	"testing"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/mock"

	"github.com/stretchr/testify/assert"
)

func TestNewRepo(t *testing.T) {
	cluster := mock.StartEtcdCluster(t)
	defer cluster.Terminate(t)
	cfg := config.RepoState{
		Endpoints: cluster.Endpoints,
	}

	factory := NewRepositoryFactory("nobody")
	repo, err := factory.CreateRepo(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, repo)
}

func TestEventType_String(t *testing.T) {
	assert.Equal(t, "delete", EventTypeDelete.String())
	assert.Equal(t, "modify", EventTypeModify.String())
	assert.Equal(t, "all", EventTypeAll.String())
	assert.Equal(t, "unknown", EventType(111).String())
}
