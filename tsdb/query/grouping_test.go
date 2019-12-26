package query

import (
	"strconv"
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func TestGroupingContext_Build(t *testing.T) {
	hosts := NewTagValuesEntrySet()
	disks := NewTagValuesEntrySet()
	partitions := NewTagValuesEntrySet()
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

	// test single group by tag keys
	ctx := NewGroupContext(1)
	ctx.SetTagValuesEntrySet(0, disks)
	assert.Equal(t, 1, ctx.Len())
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
	ctx = NewGroupContext(2)
	ctx.SetTagValuesEntrySet(0, disks)
	ctx.SetTagValuesEntrySet(1, partitions)
	assert.Equal(t, 2, ctx.Len())
	for idx, key := range keys {
		container := total.GetContainerAtIndex(idx)
		i += container.GetCardinality()
		k := key
		rs := ctx.BuildGroup(k, container)
		assert.Len(t, rs, 80)
	}
}

func TestTagValuesEntrySet(t *testing.T) {
	entry := NewTagValuesEntrySet()
	entry.AddTagValue("a", roaring.BitmapOf(1))
	entry.AddTagValue("a", roaring.BitmapOf(10))
	entry.AddTagValue("b", roaring.BitmapOf(10))
	assert.Len(t, entry.values, 2)
	assert.Equal(t, roaring.BitmapOf(1, 10), entry.Values()["a"])

	entry.SetTagValues(map[string]*roaring.Bitmap{"c": roaring.BitmapOf(200)})
	assert.Len(t, entry.values, 1)
	assert.Equal(t, roaring.BitmapOf(200), entry.Values()["c"])
}
