package database

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/constants"
	"github.com/eleme/lindb/coordinator/storage"
	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/interval"
	"github.com/eleme/lindb/pkg/option"
	"github.com/eleme/lindb/pkg/pathutil"
	"github.com/eleme/lindb/pkg/state"
	"github.com/eleme/lindb/pkg/util"
	"github.com/eleme/lindb/service"
)

var testPath = "test_data"
var validOption = option.ShardOption{Interval: time.Second * 10, IntervalType: interval.Day}
var test *testing.T

type testAdminStateMachineSuite struct {
	mock.RepoTestSuite
}

func TestAdminStateMachine(t *testing.T) {
	check.Suite(&testAdminStateMachineSuite{})
	test = t
	check.TestingT(t)
}

func (ts *testAdminStateMachineSuite) TestDatabaseShardAssign(c *check.C) {
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
	databaseSRV := service.NewDatabaseService(repo)

	clusterStateMachine, _ := storage.NewClusterStateMachine(context.TODO(), repo)
	stateMachine, _ := NewAdminStateMachine(context.TODO(), repo, clusterStateMachine)
	defer func() {
		_ = stateMachine.Close()
	}()
	ts.prepareStorageCluster(repo)
	time.Sleep(100 * time.Millisecond)

	cluster := clusterStateMachine.GetCluster("storage1")

	taskExecutor := storage.NewTaskExecutor(context.TODO(),
		&models.Node{IP: "127.0.0.5", Port: 2080},
		cluster.GetRepo(),
		storageService)
	taskExecutor.Run()

	dbCfg := models.Database{
		Name: "test",
		Clusters: []models.DatabaseCluster{
			{
				Name:          "storage1",
				NumOfShard:    10,
				ReplicaFactor: 3,
				ShardOption:   validOption,
			},
		},
	}
	_ = databaseSRV.Save(&dbCfg)

	time.Sleep(200 * time.Millisecond)

	shardAssign, err := cluster.GetShardAssign("test")
	if err != nil {
		c.Fatal(err)
	}
	checkShardAssignResult(shardAssign, test)
	c.Assert(models.DatabaseCluster{
		Name:          "storage1",
		ShardOption:   validOption,
		NumOfShard:    10,
		ReplicaFactor: 3,
	},
		check.Equals,
		shardAssign.Config)

	c.Assert(true, check.Equals, util.Exist(filepath.Join(testPath, "test", "shard")))
}

func (ts *testAdminStateMachineSuite) TestWrongCfg(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Namespace: "/admin/shard/test/wrong",
		Endpoints: ts.Cluster.Endpoints,
	})
	databaseSRV := service.NewDatabaseService(repo)

	clusterStateMachine, _ := storage.NewClusterStateMachine(context.TODO(), repo)
	stateMachine, _ := NewAdminStateMachine(context.TODO(), repo, clusterStateMachine)
	defer func() {
		_ = stateMachine.Close()
	}()

	storage1 := state.Config{
		Namespace: "/admin/shard/test/wrong",
		Endpoints: ts.Cluster.Endpoints,
	}
	data1, _ := json.Marshal(models.StorageCluster{
		Name:   "storage1",
		Config: storage1,
	})
	_ = repo.Put(context.TODO(), constants.StorageClusterConfigPath+"/storage1", data1)
	time.Sleep(200 * time.Millisecond)

	cluster := clusterStateMachine.GetCluster("storage1")

	// replica factor > num. of storage nodes
	dbCfg := models.Database{
		Name: "test2",
		Clusters: []models.DatabaseCluster{
			{
				Name:          "storage1",
				NumOfShard:    10,
				ReplicaFactor: 30,
			},
		},
	}
	_ = databaseSRV.Save(&dbCfg)
	time.Sleep(100 * time.Millisecond)
	_, err := cluster.GetShardAssign("test2")
	c.Assert(state.ErrNotExist, check.Equals, err)

	// cfg unmarshal error
	_ = repo.Put(context.TODO(), pathutil.GetDatabaseConfigPath("test3"), []byte("ddd"))
	time.Sleep(100 * time.Millisecond)
	_, err = cluster.GetShardAssign("test3")
	c.Assert(state.ErrNotExist, check.Equals, err)

	// no database name
	dbCfg = models.Database{
		Clusters: []models.DatabaseCluster{
			{
				Name:          "storage1",
				NumOfShard:    10,
				ReplicaFactor: 3,
			},
		},
	}
	data, _ := json.Marshal(dbCfg)
	_ = repo.Put(context.TODO(), pathutil.GetDatabaseConfigPath("test4"), data)
	_ = databaseSRV.Save(&dbCfg)
	time.Sleep(100 * time.Millisecond)
	_, err = cluster.GetShardAssign("test4")
	c.Assert(state.ErrNotExist, check.Equals, err)

	// storage cluster not exist
	dbCfg = models.Database{
		Name: "test5",
		Clusters: []models.DatabaseCluster{
			{
				Name:          "storage122",
				NumOfShard:    10,
				ReplicaFactor: 3,
			},
		},
	}
	_ = databaseSRV.Save(&dbCfg)
	time.Sleep(100 * time.Millisecond)
	_, err = cluster.GetShardAssign("test5")
	c.Assert(state.ErrNotExist, check.Equals, err)

	// storage has no node
	storage2 := state.Config{
		Namespace: "/admin/shard/test2",
		Endpoints: ts.Cluster.Endpoints,
	}
	data2, _ := json.Marshal(models.StorageCluster{
		Name:   "storage2",
		Config: storage2,
	})
	_ = repo.Put(context.TODO(), constants.StorageClusterConfigPath+"/storage2", data2)
	time.Sleep(100 * time.Millisecond)

	// storage cluster not exist
	dbCfg = models.Database{
		Name: "test6",
		Clusters: []models.DatabaseCluster{
			{
				Name:          "storage2",
				NumOfShard:    10,
				ReplicaFactor: 3,
			},
		},
	}
	_ = databaseSRV.Save(&dbCfg)
	time.Sleep(100 * time.Millisecond)
	cluster2 := clusterStateMachine.GetCluster("storage2")
	_, err = cluster2.GetShardAssign("test6")
	c.Assert(state.ErrNotExist, check.Equals, err)
}

func (ts *testAdminStateMachineSuite) prepareStorageCluster(repo state.Repository) {

	storage1 := state.Config{
		Namespace: "/admin/shard/test",
		Endpoints: ts.Cluster.Endpoints,
	}
	data1, _ := json.Marshal(models.StorageCluster{
		Name:   "storage1",
		Config: storage1,
	})
	_ = repo.Put(context.TODO(), constants.StorageClusterConfigPath+"/storage1", data1)
	node1, _ := json.Marshal(models.Node{IP: "127.0.0.1", Port: 2080})
	_ = repo.Put(context.TODO(), constants.ActiveNodesPath+"/node1", node1)
	node2, _ := json.Marshal(models.Node{IP: "127.0.0.2", Port: 2080})
	_ = repo.Put(context.TODO(), constants.ActiveNodesPath+"/node2", node2)
	node3, _ := json.Marshal(models.Node{IP: "127.0.0.3", Port: 2080})
	_ = repo.Put(context.TODO(), constants.ActiveNodesPath+"/node3", node3)
	node4, _ := json.Marshal(models.Node{IP: "127.0.0.4", Port: 2080})
	_ = repo.Put(context.TODO(), constants.ActiveNodesPath+"/node4", node4)
	node5, _ := json.Marshal(models.Node{IP: "127.0.0.5", Port: 2080})
	_ = repo.Put(context.TODO(), constants.ActiveNodesPath+"/node5", node5)

}
