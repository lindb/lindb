package broker

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/rpc"
	brokerpb "github.com/eleme/lindb/rpc/pkg/broker"
)

const (
	requestTimeout = 2 * time.Second
	address        = ":9001"
)

type testSuite struct {
	server rpc.Server
}

func Test(t *testing.T) {
	check.TestingT(t)
}

var _ = check.Suite(&testSuite{
	server: NewBrokerServer(address),
})

func (ts *testSuite) SetUpSuite(c *check.C) {
	ts.server.Init()
	ts.server.Listen()
	ts.server.Register()
	go func() {
		ts.server.Serve()
	}()
}

func (ts *testSuite) TearDownSuite(c *check.C) {
}

func (ts *testSuite) TestRequest(c *check.C) {
	cli := NewBrokerClient(address, requestTimeout)

	wg := sync.WaitGroup{}

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			points := buildPoints()
			err := cli.WritePoints(points)
			c.Assert(err, check.IsNil)
		}()
	}

	wg.Wait()
	cli.Close()
}

func (ts *testSuite) TestStreamRequest(c *check.C) {
	cli := NewBrokerClient(address, requestTimeout)

	wg := sync.WaitGroup{}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			points := buildPoints()

			err := cli.WritePointsStream(points)
			c.Assert(err, check.IsNil)

			err = cli.WritePointsStream(points)
			c.Assert(err, check.IsNil)

			err = cli.WritePointsStream(points)
			c.Assert(err, check.IsNil)

		}()
	}

	wg.Wait()
	cli.Close()
}

func buildPoints() (points []*brokerpb.Point) {
	for i := range []int{1, 2, 3} {
		name := fmt.Sprintf("name-%d", i)
		tags := map[string]string{"t1": "v1"}
		f1 := buildField("f1", brokerpb.AggType_GAUGE, 1)
		f2 := buildField("f2", brokerpb.AggType_SUM, 2)
		fields := []*brokerpb.Field{f1, f2}
		points = append(points, buildPoint(name, int64(i), tags, fields))
	}
	return points
}

func buildPoint(name string, ts int64, tags map[string]string, fields []*brokerpb.Field) *brokerpb.Point {
	return &brokerpb.Point{
		Name:      name,
		Timestamp: ts,
		Tags:      tags,
		Fields:    fields,
	}
}

func buildField(name string, aggType brokerpb.AggType, value float64) *brokerpb.Field {
	return &brokerpb.Field{
		Name:    name,
		AggType: aggType,
		Value:   value,
	}
}
