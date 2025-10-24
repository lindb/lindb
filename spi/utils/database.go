package utils

import (
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/storage/store"
)

func FindSegments(db store.Database, shardIDs []int,
	interval timeutil.Interval, timeRange timeutil.TimeRange,
	fn func(shard store.Shard, partition store.Partition, segments []store.Segment),
) {
	storageInterval := db.FindMatchSmallestInterval(interval)
	for _, id := range shardIDs {
		shard, ok := db.GetShard(models.ShardID(id))
		if ok {
			pList := shard.GetPartitions(storageInterval, timeRange)
			if len(pList) > 0 {
				for _, partition := range pList {
					segments := partition.GetSegments(timeRange)
					if len(segments) > 0 {
						fn(shard, partition, segments)
					}
				}
			}
		}
	}
}
