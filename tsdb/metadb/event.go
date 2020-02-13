package metadb

import (
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/series/tag"
)

const full = 10000

// metricEvent represents metric name/id need store
type metricEvent struct {
	name string
	id   uint32
}

// metricMetadataEvent represents the metric metadata include fields/tags keys need store
type metricMetadataEvent struct {
	fieldIDSeq uint16
	fields     []field.Meta
	tagKeys    []tag.Meta
}

// newMetricMetadataEvent creates a metric metadata event
func newMetricMetadataEvent() *metricMetadataEvent {
	return &metricMetadataEvent{}
}

// namespaceEvent represents the namespace include metrics need store
type namespaceEvent struct {
	metrics []metricEvent // namespace => (metricName=>metric event)
}

// newNamespaceEvent creates a namespace event
func newNamespaceEvent() *namespaceEvent {
	return &namespaceEvent{}
}

// metadataUpdateEvent represents the metadata include namespace/metric  metadata need store
type metadataUpdateEvent struct {
	metricSeqID uint32
	tagKeySeqID uint32
	namespaces  map[string]*namespaceEvent
	metrics     map[uint32]*metricMetadataEvent

	pending int
}

// newMetadataUpdateEvent creates a metadata update event
func newMetadataUpdateEvent() *metadataUpdateEvent {
	return &metadataUpdateEvent{
		namespaces: make(map[string]*namespaceEvent),
		metrics:    make(map[uint32]*metricMetadataEvent),
	}
}

// addMetric adds metric into namespace
func (e *metadataUpdateEvent) addMetric(namespace, metricName string, metricID uint32) {
	// set metric seq id directly, because gen metric id in order
	e.metricSeqID = metricID

	ns, ok := e.namespaces[namespace]
	if !ok {
		ns = newNamespaceEvent()
		e.namespaces[namespace] = ns
	}
	ns.metrics = append(ns.metrics, metricEvent{
		name: metricName,
		id:   metricID,
	})
	e.pending++
}

// addField adds field into metric metadata event
func (e *metadataUpdateEvent) addField(metricID uint32, f field.Meta) {
	metricMeta, ok := e.metrics[metricID]
	if !ok {
		metricMeta = newMetricMetadataEvent()
		e.metrics[metricID] = metricMeta
	}
	// set field seq id directly, because gen field id in order
	metricMeta.fieldIDSeq = uint16(f.ID)

	metricMeta.fields = append(metricMeta.fields, f)
	e.pending++
}

// addTagKey adds tag key into metric metadata event
func (e *metadataUpdateEvent) addTagKey(metricID uint32, tagKey tag.Meta) {
	// set tag key seq id directly, because gen tag key id in order
	e.tagKeySeqID = tagKey.ID

	metricMeta, ok := e.metrics[metricID]
	if !ok {
		metricMeta = newMetricMetadataEvent()
		e.metrics[metricID] = metricMeta
	}
	metricMeta.tagKeys = append(metricMeta.tagKeys, tagKey)
	e.pending++
}

// isFull returns if update event is full
func (e *metadataUpdateEvent) isFull() bool {
	return e.pending >= full
}

// isEmpty returns if update event is empty
func (e *metadataUpdateEvent) isEmpty() bool {
	return e.pending == 0
}
