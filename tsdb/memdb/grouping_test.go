package memdb

import (
	"strconv"
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func TestGroupingContext_Build(t *testing.T) {
	hosts := newTagKVEntrySet("host")
	disks := newTagKVEntrySet("disk")
	partitions := newTagKVEntrySet("partition")
	id := uint32(0)
	count := 12500
	for i := 0; i < count; i++ {
		for j := 0; j < 4; j++ {
			for k := 0; k < 20; k++ {
				id++
				host := "host" + strconv.Itoa(i)
				disk := "/tmp" + strconv.Itoa(j)
				partition := "partition" + strconv.Itoa(k)
				h, ok := hosts.values[host]
				if !ok {
					hosts.values[host] = roaring.BitmapOf(id)
				} else {
					h.Add(id)
				}

				d, ok := disks.values[disk]
				if !ok {
					disks.values[disk] = roaring.BitmapOf(id)
				} else {
					d.Add(id)
				}

				p, ok := partitions.values[partition]
				if !ok {
					partitions.values[partition] = roaring.BitmapOf(id)
				} else {
					p.Add(id)
				}
			}
		}
	}

	mStore := newMetricStore()
	ms := mStore.(*metricStore)

	// test single group by tag keys
	ctx := &groupingContext{
		ms:             ms,
		tagKVEntrySets: []*tagKVEntrySet{disks},
	}
	total := roaring.New()
	total.AddRange(1, 1000001)
	keys := total.GetHighKeys()
	i := 0
	for idx, key := range keys {
		container := total.GetContainerAtIndex(idx)
		i += container.GetCardinality()
		k := key
		rs := ctx.BuildGroup(k, container)
		assert.Len(t, rs, 4)
	}
	// test single group by tag keys
	ctx = &groupingContext{
		ms:             ms,
		tagKVEntrySets: []*tagKVEntrySet{disks, partitions},
	}
	for idx, key := range keys {
		container := total.GetContainerAtIndex(idx)
		i += container.GetCardinality()
		k := key
		rs := ctx.BuildGroup(k, container)
		assert.Len(t, rs, 80)
	}
}
