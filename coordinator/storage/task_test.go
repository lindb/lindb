package storage

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/interval"
	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/pkg/util"
	"github.com/eleme/lindb/service"
)

var testPath = "test_data"
var validOption = option.ShardOption{Interval: time.Second * 10, IntervalType: interval.Day}

type testTaskExecutorSuite struct {
	mock.RepoTestSuite
}

func TestAdminStateMachine(t *testing.T) {
	check.Suite(&testTaskExecutorSuite{})
	check.TestingT(t)
}

func (ts *testTaskExecutorSuite) TestCreateShard(c *check.C) {
	defer func() {
		_ = util.RemoveDir(testPath)
	}()

	cfg := config.Engine{
		Path: testPath,
	}
	storageService := service.NewStorageService(cfg)

	repo, _ := state.NewRepo(state.Config{
		Namespace: "/admin/shard/test",
		Endpoints: ts.Cluster.Endpoints,
	})
	node := models.Node{IP: "127.0.0.5", Port: 2080}
	taskExecutor := NewTaskExecutor(context.TODO(), &node, repo, storageService)
	taskExecutor.Run()

	cluster, _ := newCluster(context.TODO(), models.StorageCluster{Config: state.Config{
		Namespace: "/admin/shard/test",
		Endpoints: ts.Cluster.Endpoints,
	}})
	nodes := make(map[int]models.Node)
	nodes[1] = node
	shardAssign := models.NewShardAssignment()
	shardAssign.AddReplica(1, 1)
	_ = cluster.SaveShardAssign("test",
		&models.ShardAssignment{
			Nodes:  nodes,
			Shards: shardAssign.Shards,
			Config: models.DatabaseCluster{
				ShardOption: validOption,
			},
		},
	)
	time.Sleep(100 * time.Millisecond)

	c.Assert(true, check.Equals, util.Exist(filepath.Join(testPath, "test", "shard")))
}
