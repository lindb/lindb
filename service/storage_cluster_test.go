package service

import (
	"testing"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"
)

type testStorageClusterSRVSuite struct {
	mock.RepoTestSuite
}

func TestStorageClusterSRV(t *testing.T) {
	check.Suite(&testStorageClusterSRVSuite{})
	check.TestingT(t)
}

func (ts *testStorageClusterSRVSuite) TestStorageCluster(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	cluster := models.StorageCluster{
		Name: "test1",
	}
	srv := NewStorageClusterService(repo)
	err := srv.Save(cluster)
	if err != nil {
		c.Fatal(err)
	}
	err = srv.Save(models.StorageCluster{})
	c.Assert(err, check.NotNil)

	cluster2, _ := srv.Get("test1")
	c.Assert(cluster, check.DeepEquals, cluster2)

	_ = srv.Save(models.StorageCluster{
		Name: "test2",
	})
	clusterList, _ := srv.List()
	c.Assert(len(clusterList), check.Equals, 2)

	_ = srv.Delete("test1")

	_, err2 := srv.Get("test1")
	c.Assert(err2, check.NotNil)
}
