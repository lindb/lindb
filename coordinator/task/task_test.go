package task

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	etcdcliv3 "github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/integration"
	"github.com/coreos/pkg/capnslog"
	"github.com/stretchr/testify/require"

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

func Test_tasks(t *testing.T) {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)
	config := etcdcliv3.Config{
		Endpoints: []string{clus.Members[0].GRPCAddr()},
	}
	_ = state.New("etcd", config)
	cli, err := etcdcliv3.New(config)
	require.Nil(t, err)
	repo := state.GetRepo()
	ctx := context.TODO()
	keypfx := "/let-me-through"

	controller := NewController(ctx, keypfx, repo)
	defer controller.Close()
	err = controller.Submit(kindDummy, "wtf-2019-07-05--1", []ControllerTaskParam{
		{NodeID: "her", Params: dummyParams{}},
		{NodeID: "him", Params: dummyParams{}},
	})
	require.Nil(t, err)

	processor := &dummyProcessor{}
	executor1 := NewExecutor(ctx, keypfx, "her", repo)
	executor1.Register(processor)
	go executor1.Run()
	time.Sleep(333 * time.Millisecond)
	require.Equal(t, 1, processor.CallCount())
	executor1.Close()

	executor2 := NewExecutor(ctx, keypfx, "him", repo)
	executor2.Register(processor)
	go executor2.Run()
	time.Sleep(333 * time.Millisecond)
	require.Equal(t, 2, processor.CallCount())
	executor2.Close()

	time.Sleep(666 * time.Millisecond)
	resp, err := cli.Get(ctx, keypfx, etcdcliv3.WithPrefix())
	require.Nil(t, err)
	require.Equal(t, 1, len(resp.Kvs))
	var tasks groupedTasks
	(&tasks).UnsafeUnmarshal(resp.Kvs[0].Value)
	require.Equal(t, StateDoneOK, tasks.State)
	for _, task := range tasks.Tasks {
		require.Equal(t, StateDoneOK, task.State)
	}
}
