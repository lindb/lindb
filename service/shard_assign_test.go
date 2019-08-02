package service

import (
	"testing"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"

	"gopkg.in/check.v1"
)

type testShardAssignSRVSuite struct {
	mock.RepoTestSuite
}

func TestShardAssignSRV(t *testing.T) {
	check.Suite(&testShardAssignSRVSuite{})
	check.TestingT(t)
}

func (ts *testShardAssignSRVSuite) TestShardAssign(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	srv := NewShardAssignService(repo)

	shardAssign1 := models.NewShardAssignment()
	shardAssign1.AddReplica(1, 1)
	shardAssign1.AddReplica(1, 2)
	shardAssign1.AddReplica(1, 3)
	shardAssign1.AddReplica(2, 2)
	_ = srv.Save("db1", shardAssign1)

	shardAssign2 := models.NewShardAssignment()
	shardAssign2.AddReplica(1, 1)
	shardAssign2.AddReplica(2, 2)
	_ = srv.Save("db2", shardAssign2)

	shardAssign11, _ := srv.Get("db1")
	c.Assert(*shardAssign1, check.DeepEquals, *shardAssign11)

	shardAssign22, _ := srv.Get("db2")
	c.Assert(*shardAssign2, check.DeepEquals, *shardAssign22)

	_, err := srv.Get("not_exist")
	c.Assert(state.ErrNotExist, check.Equals, err)

	_ = repo.Close()

	// test error
	err = srv.Save("db2", shardAssign2)
	c.Assert(err, check.NotNil)
	_, err = srv.Get("db2")
	c.Assert(err, check.NotNil)
}
