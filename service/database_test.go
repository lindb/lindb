package service

import (
	"testing"

	check "gopkg.in/check.v1"

	"github.com/lindb/lindb/mock"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

type testDatabaseSRVSuite struct {
	mock.RepoTestSuite
}

func TestDatabaseSRV(t *testing.T) {
	check.Suite(&testDatabaseSRVSuite{})
	check.TestingT(t)
}

func (ts *testDatabaseSRVSuite) TestDatabase(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})

	db := NewDatabaseService(repo)
	database := models.Database{
		Name: "test",
		Clusters: []models.DatabaseCluster{
			{
				Name:          "test",
				NumOfShard:    12,
				ReplicaFactor: 3,
			},
		},
	}
	err := db.Save(&database)
	if err != nil {
		c.Fatal(err)
	}
	database2, _ := db.Get("test")
	c.Assert(database, check.DeepEquals, *database2)
	database2, err = db.Get("test_not_exist")
	c.Assert(err, check.Equals, state.ErrNotExist)
	c.Assert(database2, check.IsNil)

	// test create database error
	err = db.Save(&models.Database{})
	c.Assert(err, check.NotNil)

	err = db.Save(&models.Database{
		Name: "test",
	})
	c.Assert(err, check.NotNil)

	err = db.Save(&models.Database{
		Name: "test",
		Clusters: []models.DatabaseCluster{
			{
				NumOfShard:    12,
				ReplicaFactor: 3,
			},
		},
	})
	c.Assert(err, check.NotNil)

	err = db.Save(&models.Database{
		Name: "test",
		Clusters: []models.DatabaseCluster{
			{
				Name:          "test",
				ReplicaFactor: 3,
			},
		},
	})
	c.Assert(err, check.NotNil)

	err = db.Save(&models.Database{
		Name: "test",
		Clusters: []models.DatabaseCluster{
			{
				Name:       "test",
				NumOfShard: 3,
			},
		},
	})
	c.Assert(err, check.NotNil)

	_ = repo.Close()

	err = db.Save(&database)
	c.Assert(err, check.NotNil)
	database2, err = db.Get("test")
	c.Assert(database2, check.IsNil)
	c.Assert(err, check.NotNil)
}
