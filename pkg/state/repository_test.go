package state

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/mock"
)

func TestNewRepo(t *testing.T) {
	cluster := mock.StartEtcdCluster(t)
	defer cluster.Terminate(t)
	cfg := Config{
		Endpoints: cluster.Endpoints,
	}

	factory := NewRepositoryFactory()
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
