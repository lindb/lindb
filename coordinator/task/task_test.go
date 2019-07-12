package task

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"gopkg.in/check.v1"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/state"
)

const kindDummy Kind = "you-guess"

type dummyParams struct{}

func (p dummyParams) Bytes() []byte { return []byte("{}") }

type dummyProcessor struct{ callcnt int32 }

func (p *dummyProcessor) Kind() Kind                  { return kindDummy }
func (p *dummyProcessor) RetryCount() int             { return 0 }
func (p *dummyProcessor) RetryBackOff() time.Duration { return 0 }
func (p *dummyProcessor) Concurrency() int            { return 1 }
func (p *dummyProcessor) Process(ctx context.Context, task Task) error {
	atomic.AddInt32(&p.callcnt, 1)
	return nil
}
func (p *dummyProcessor) CallCount() int { return int(atomic.LoadInt32(&p.callcnt)) }

type testTaskSuite struct {
	mock.RepoTestSuite
}

func TestElection(t *testing.T) {
	check.Suite(&testTaskSuite{})
	check.TestingT(t)
}

func (ts *testTaskSuite) Test_tasks(c *check.C) {
	repo, _ := state.NewRepo(state.Config{
		Namespace: "/coordinator/test/task",
		Endpoints: ts.Cluster.Endpoints,
	})
	ctx := context.TODO()

	controller := NewController(ctx, repo)
	node1 := &models.Node{IP: "1.1.1.1", Port: 8000}
	node2 := &models.Node{IP: "1.1.1.2", Port: 8000}
	defer func() {
		_ = controller.Close()
	}()
	err := controller.Submit(kindDummy, "wtf-2019-07-05--1", []ControllerTaskParam{
		{NodeID: node1.String(), Params: dummyParams{}},
		{NodeID: node2.String(), Params: dummyParams{}},
	})
	if err != nil {
		c.Fatal(err)
	}

	processor := &dummyProcessor{}
	executor1 := NewExecutor(ctx, node1, repo)
	executor1.Register(processor)
	go executor1.Run()
	time.Sleep(333 * time.Millisecond)

	c.Assert(1, check.Equals, processor.CallCount())

	_ = executor1.Close()

	executor2 := NewExecutor(ctx, node2, repo)
	executor2.Register(processor)
	go executor2.Run()
	time.Sleep(333 * time.Millisecond)
	c.Assert(2, check.Equals, processor.CallCount())
	_ = executor2.Close()

	//time.Sleep(666 * time.Millisecond)
	//resp, err := cli.Get(ctx, keypfx, etcdcliv3.WithPrefix())
	//if err != nil {
	//	c.Fatal(err)
	//}
	//c.Assert(1, check.Equals, len(resp.Kvs))

	//var tasks groupedTasks
	//(&tasks).UnsafeUnmarshal(resp.Kvs[0].Value)

	//c.Assert(StateDoneOK, check.Equals, tasks.State)
	//fail := true
	//for _, task := range tasks.Tasks {
	//	c.Assert(StateDoneOK, check.Equals, task.State)
	//	fail = false
	//}
	//c.Assert(false, check.Equals, fail)
}
