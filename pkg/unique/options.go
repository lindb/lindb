// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package unique

import (
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/bloom"
)

// DefaultOptions returns the default options of pebble kv store.
// ref: https://github.com/cockroachdb/cockroach/blob/master/pkg/storage/pebble.go#L445-L498
func DefaultOptions() *pebble.Options {
	// In RocksDB, the concurrency setting corresponds to both flushes and
	// compactions. In Pebble, there is always a slot for a flush, and
	// compactions are counted separately.
	// maxConcurrentCompactions := rocksdbConcurrency - 1
	maxConcurrentCompactions := 1

	opts := &pebble.Options{
		Comparer:                    pebble.DefaultComparer,
		L0CompactionThreshold:       2,
		L0StopWritesThreshold:       1000,
		LBaseMaxBytes:               64 << 20, // 64 MB
		Levels:                      make([]pebble.LevelOptions, 7),
		MaxConcurrentCompactions:    maxConcurrentCompactions,
		MemTableSize:                64 << 20, // 64 MB
		MemTableStopWritesThreshold: 4,
		Merger:                      pebble.DefaultMerger,
		DisableWAL:                  true,
	}
	// Automatically flush 10s after the first range tombstone is added to a
	// memtable. This ensures that we can reclaim space even when there's no
	// activity on the database generating flushes.
	opts.Experimental.DeleteRangeFlushDelay = 10 * time.Second
	// Enable deletion pacing. This helps prevent disk slowness events on some
	// SSDs, that kick off an expensive GC if a lot of files are deleted at
	// once.
	opts.Experimental.MinDeletionRate = 128 << 20 // 128 MB
	// Validate min/max keys in each SSTable when performing a compaction. This
	// serves as a simple protection against corruption or programmer-error in
	// Pebble.

	for i := 0; i < len(opts.Levels); i++ {
		l := &opts.Levels[i]
		l.BlockSize = 32 << 10       // 32 KB
		l.IndexBlockSize = 256 << 10 // 256 KB
		l.FilterPolicy = bloom.FilterPolicy(10)
		l.FilterType = pebble.TableFilter
		if i > 0 {
			l.TargetFileSize = opts.Levels[i-1].TargetFileSize * 2
		}
		l.EnsureDefaults()
	}

	// Do not create bloom filters for the last level (i.e. the largest level
	// which contains data in the LSM store). This configuration reduces the size
	// of the bloom filters by 10x. This is significant given that bloom filters
	// require 1.25 bytes (10 bits) per key which can translate into gigabytes of
	// memory given typical key and value sizes. The downside is that bloom
	// filters will only be usable on the higher levels, but that seems
	// acceptable. We typically see read amplification of 5-6x on clusters
	// (i.e. there are 5-6 levels of sstables) which means we'll achieve 80-90%
	// of the benefit of having bloom filters on every level for only 10% of the
	// memory cost.
	// opts.Levels[6].FilterPolicy = nil
	opts.Levels[len(opts.Levels)-1].FilterPolicy = nil

	// Set disk health check interval to min(5s, maxSyncDurationDefault). This
	// is mostly to ease testing; the default of 5s is too infrequent to test
	// conveniently. See the disk-stalled roachtest for an example of how this
	// is used.
	// diskHealthCheckInterval := 5 * time.Second
	// if diskHealthCheckInterval.Seconds() > maxSyncDurationDefault.Seconds() {
	//	diskHealthCheckInterval = maxSyncDurationDefault
	// }
	// Instantiate a file system with disk health checking enabled. This FS wraps
	// vfs.Default, and can be wrapped for encryption-at-rest.
	// opts.FS = vfs.WithDiskHealthChecks(vfs.Default, diskHealthCheckInterval,
	//	func(name string, duration time.Duration) {
	//		opts.EventListener.DiskSlow(pebble.DiskSlowInfo{
	//			Path:     name,
	//			Duration: duration,
	//		})
	//	})
	// If we encounter ENOSPC, exit with an informative exit code.
	// opts.FS = vfs.OnDiskFull(opts.FS, func() {
	//	exit.WithCode(exit.DiskFull())
	// })
	return opts
}
