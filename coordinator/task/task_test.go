package task

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	etcdcliv3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/pkg/capnslog"
	"gopkg.in/check.v1"

	"github.com/eleme/lindb/mock"
	"github.com/eleme/lindb/pkg/state"
)

func init() {
	capnslog.SetGlobalLogLevel(capnslog.CRITICAL)
}

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

var _ = check.Suite(&testTaskSuite{})

func TestElection(t *testing.T) {
	check.TestingT(t)
}

func (ts *testTaskSuite) Test_tasks(c *check.C) {
	config := etcdcliv3.Config{
		Endpoints: ts.Cluster.Endpoints,
	}
	repo, _ := state.NewRepo(state.Config{
		Endpoints: ts.Cluster.Endpoints,
	})
	cli, err := etcdcliv3.New(config)
	if err != nil {
		c.Fatal(err)
	}
	ctx := context.TODO()
	keypfx := "/let-me-through"

	controller := NewController(ctx, keypfx, repo)
	defer controller.Close()
	err = controller.Submit(kindDummy, "wtf-2019-07-05--1", []ControllerTaskParam{
		{NodeID: "her", Params: dummyParams{}},
		{NodeID: "him", Params: dummyParams{}},
	})
	if err != nil {
		c.Fatal(err)
	}

	processor := &dummyProcessor{}
	executor1 := NewExecutor(ctx, keypfx, "her", repo)
	executor1.Register(processor)
	go executor1.Run()
	time.Sleep(333 * time.Millisecond)

	c.Assert(1, check.Equals, processor.CallCount())

	executor1.Close()

	executor2 := NewExecutor(ctx, keypfx, "him", repo)
	executor2.Register(processor)
	go executor2.Run()
	time.Sleep(333 * time.Millisecond)
	c.Assert(2, check.Equals, processor.CallCount())
	executor2.Close()

	time.Sleep(666 * time.Millisecond)
	resp, err := cli.Get(ctx, keypfx, etcdcliv3.WithPrefix())
	if err != nil {
		c.Fatal(err)
	}
	c.Assert(1, check.Equals, len(resp.Kvs))

	var tasks groupedTasks
	(&tasks).UnsafeUnmarshal(resp.Kvs[0].Value)

	c.Assert(StateDoneOK, check.Equals, tasks.State)
	fail := true
	for _, task := range tasks.Tasks {
		c.Assert(StateDoneOK, check.Equals, task.State)
		fail = false
	}
	c.Assert(false, check.Equals, fail)
}
