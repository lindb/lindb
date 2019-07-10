package state

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eleme/lindb/mock"
)

func TestNewRepo(t *testing.T) {
	cluster := mock.StartEtcdCluster(t)
	defer cluster.Terminate(t)
	cfg := Config{
		Endpoints: cluster.Endpoints,
	}

	repo, err := NewRepo(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, repo)
}
