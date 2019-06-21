package rpc

import (
	"testing"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/rpc/pkg/batch"
)

type rpcTestSuite struct {
}

var _ = check.Suite(&rpcTestSuite{})

func Test(t *testing.T) {
	check.TestingT(t)
}

func (ts *rpcTestSuite) TestRemoveEntries(c *check.C) {
	entries := []*batchRequestEntry{
		{canceled: 1},
		{canceled: 0},
		{canceled: 1},
		{canceled: 0},
		{canceled: 1},
	}

	firstEntry := &entries[0]

	reqs := make([]*batch.BatchRequest_Request, len(entries), len(entries))
	length := removeCanceledRequests(&entries, &reqs)
	c.Assert(length, check.Equals, 2)
	c.Assert(len(entries), check.Equals, 2)

	newFirstEntry := &entries[0]

	c.Assert(firstEntry, check.Equals, newFirstEntry)
}
