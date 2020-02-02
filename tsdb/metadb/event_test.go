package metadb

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

func TestNamespaceUpdateEvent(t *testing.T) {
	e := newMetadataUpdateEvent()
	assert.True(t, e.isEmpty())
	e.addMetric("ns-1", "name1", 1)
	e.addMetric("ns-1", "name2", 2)
	e.addMetric("ns-2", "name3", 3)
	e.addMetric("ns-2", "name2", 4)
	assert.False(t, e.isEmpty())
	assert.Equal(t, uint32(4), e.metricSeqID)
	assert.Len(t, e.namespaces, 2)
	assert.Len(t, e.metrics, 0)
	metrics := e.namespaces["ns-1"]
	assert.Len(t, metrics.metrics, 2)
	assert.Equal(t, uint32(1), metrics.metrics[0].id)
	assert.Equal(t, "name1", metrics.metrics[0].name)
	assert.Equal(t, uint32(2), metrics.metrics[1].id)
	assert.Equal(t, "name2", metrics.metrics[1].name)
	metrics = e.namespaces["ns-2"]
	assert.Len(t, metrics.metrics, 2)
	assert.Equal(t, uint32(3), metrics.metrics[0].id)
	assert.Equal(t, "name3", metrics.metrics[0].name)
	assert.Equal(t, uint32(4), metrics.metrics[1].id)
	assert.Equal(t, "name2", metrics.metrics[1].name)

	assert.False(t, e.isFull())
	for i := 0; i < full; i++ {
		e.addMetric("ns-2", "name2", 4)
	}
	assert.True(t, e.isFull())
}

func TestMetricUpdateEvent_TagKeys(t *testing.T) {
	e := newMetadataUpdateEvent()
	e.addTagKey(1, tag.Meta{Key: "tagKey-1", ID: 1})
	e.addTagKey(1, tag.Meta{Key: "tagKey-2", ID: 2})
	e.addTagKey(2, tag.Meta{Key: "tagKey-3", ID: 3})
	e.addTagKey(2, tag.Meta{Key: "tagKey-2", ID: 4})
	assert.Equal(t, uint32(4), e.tagKeySeqID)
	assert.Len(t, e.namespaces, 0)
	assert.Len(t, e.metrics, 2)

	metric := e.metrics[1]
	assert.Len(t, metric.tagKeys, 2)
	assert.Len(t, metric.fields, 0)
	assert.Equal(t, uint16(0), metric.fieldIDSeq)
	assert.Equal(t, tag.Meta{Key: "tagKey-1", ID: 1}, metric.tagKeys[0])
	assert.Equal(t, tag.Meta{Key: "tagKey-2", ID: 2}, metric.tagKeys[1])

	metric = e.metrics[2]
	assert.Len(t, metric.tagKeys, 2)
	assert.Len(t, metric.fields, 0)
	assert.Equal(t, uint16(0), metric.fieldIDSeq)
	assert.Equal(t, tag.Meta{Key: "tagKey-3", ID: 3}, metric.tagKeys[0])
	assert.Equal(t, tag.Meta{Key: "tagKey-2", ID: 4}, metric.tagKeys[1])
}

func TestMetricUpdateEvent_Fields(t *testing.T) {
	e := newMetadataUpdateEvent()
	e.addField(1, field.Meta{ID: 1})
	e.addField(1, field.Meta{ID: 2})
	e.addField(2, field.Meta{ID: 1})
	e.addField(2, field.Meta{ID: 3})
	assert.Equal(t, uint32(0), e.tagKeySeqID)
	assert.Len(t, e.namespaces, 0)
	assert.Len(t, e.metrics, 2)

	metric := e.metrics[1]
	assert.Len(t, metric.tagKeys, 0)
	assert.Len(t, metric.fields, 2)
	assert.Equal(t, uint16(2), metric.fieldIDSeq)
	assert.Equal(t, field.Meta{ID: 1}, metric.fields[0])
	assert.Equal(t, field.Meta{ID: 2}, metric.fields[1])

	metric = e.metrics[2]
	assert.Len(t, metric.tagKeys, 0)
	assert.Len(t, metric.fields, 2)
	assert.Equal(t, uint16(3), metric.fieldIDSeq)
	assert.Equal(t, field.Meta{ID: 1}, metric.fields[0])
	assert.Equal(t, field.Meta{ID: 3}, metric.fields[1])
}
