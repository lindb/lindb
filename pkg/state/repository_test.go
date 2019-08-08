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
