package memdb

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/timeutil"
	pb "github.com/lindb/lindb/rpc/proto/field"
	"github.com/lindb/lindb/tsdb/metadb"
)

func BenchmarkMemoryDatabase_write(b *testing.B) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	kvStore, err := kv.NewStore("test", kv.StoreOption{
		Path:                 filepath.Join(testPath, "kv"),
		Levels:               2,
		CompactCheckInterval: 5 * 60,
	})
	assert.NoError(b, err)

	metadata, err := metadb.NewMetadata(context.TODO(), "test", filepath.Join(testPath, "meta"),
		kvStore.GetFamily("meta"))
	assert.NoError(b, err)

	metricID, err := metadata.MetadataDatabase().GenMetricID("ns", "test")
	assert.NoError(b, err)
	var cfg = MemoryDatabaseCfg{
		Interval: timeutil.Interval(10 * timeutil.OneSecond),
		Metadata: metadata,
		TempPath: filepath.Join(testPath, "data_temp"),
	}
	db, err := NewMemoryDatabase(cfg)
	if err != nil {
		b.Fatal(err)
	}
	now := timeutil.Now()
	for i := 0; i < 3200000; i++ {
		_ = db.Write("ns", "test", metricID, uint32(i), now, []*pb.Field{{
			Name:  "f1",
			Type:  pb.FieldType_Sum,
			Value: 10.0,
		}})
	}
	runtime.GC()
	fmt.Printf("cost:=%d\n", timeutil.Now()-now)
	now = timeutil.Now()
	for i := 0; i < 3200000; i++ {
		_ = db.Write("ns", "test", metricID, uint32(i), now, []*pb.Field{{
			Name:  "f1",
			Type:  pb.FieldType_Sum,
			Value: 10.0,
		}})
	}
	runtime.GC()
	fmt.Printf("cost:=%d\n", timeutil.Now()-now)
	_ = http.ListenAndServe("0.0.0.0:6060", nil)
}

func BenchmarkMemoryDatabase_write_sum(b *testing.B) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	kvStore, err := kv.NewStore("test", kv.StoreOption{
		Path:                 filepath.Join(testPath, "kv"),
		Levels:               2,
		CompactCheckInterval: 5 * 60,
	})
	assert.NoError(b, err)

	metadata, err := metadb.NewMetadata(context.TODO(), "test",
		filepath.Join(testPath, "meta"), kvStore.GetFamily("meta"))
	assert.NoError(b, err)

	metricID, err := metadata.MetadataDatabase().GenMetricID("ns", "test")
	assert.NoError(b, err)
	run := func(n int) {
		var cfg = MemoryDatabaseCfg{
			Interval: timeutil.Interval(10 * timeutil.OneSecond),
			Metadata: metadata,
			TempPath: filepath.Join(testPath, "data_temp", fmt.Sprintf("%d", n)),
		}
		db, err := NewMemoryDatabase(cfg)
		if err != nil {
			b.Fatal(err)
		}
		now := timeutil.Now()
		for i := 0; i < 400000; i++ {
			_ = db.Write("ns", "test", metricID, uint32(i), now, []*pb.Field{{
				Name:  "f1",
				Type:  pb.FieldType_Sum,
				Value: 10.0,
			}})
		}
		fmt.Printf("n:=%d, cost:=%d\n", n, timeutil.Now()-now)
	}
	now := timeutil.Now()
	var wait sync.WaitGroup
	n := 4
	wait.Add(n)
	go func() {
		run(0)
		wait.Done()
	}()
	go func() {
		run(1)
		wait.Done()
	}()
	go func() {
		run(2)
		wait.Done()
	}()
	go func() {
		run(3)
		wait.Done()
	}()
	wait.Wait()
	fmt.Println(timeutil.Now() - now)
	run(0)
}
