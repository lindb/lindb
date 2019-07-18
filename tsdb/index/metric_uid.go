package index

import (
	"go.uber.org/zap"

	"github.com/eleme/lindb/kv"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/tree"
	"github.com/eleme/lindb/pkg/util"
)

//MetricUid represents metric name unique id under the database
type MetricUID struct {
	partition  uint32
	metrics    *tree.BTree //Key is the ascii of the first letter of metric
	sequenceID uint32      //unique sequence id
	family     kv.Family
	dbField    zap.Field
}

//NewMetricUID creation requires kvFamily
func NewMetricUID(metricFamily kv.Family) *MetricUID {
	//Get the last written sequenceID, otherwise start from 0
	seq := uint32(0)
	metricFamily.Lookup(MetricSequenceIDKey, func(bytes []byte) bool {
		seq = util.BytesToUint32(bytes)
		return true
	})
	return &MetricUID{
		metrics:    tree.NewBTree(),
		sequenceID: seq,
		family:     metricFamily,
		dbField:    zap.String("db", "db"),
	}
}

//GetOrCreateMetricID returns find the metric ID associated with a given name or create it.
func (m *MetricUID) GetOrCreateMetricID(metricName string, create bool) (uint32, bool) {
	if len(metricName) > 0 {
		nameBytes := []byte(metricName)
		partition := getPartition(nameBytes)

		id := m.getMetricIDFromDisk(partition, nameBytes)
		if id == NotFoundMetricID {
			// if not exists
			if partition != m.partition {
				err := m.Flush()
				if nil != err {
					logger.GetLogger("tsdb/index").Error("flush metric tree error!", m.dbField)
				}
				m.partition = partition
				m.metrics.Clear()
			}
			if create {
				m.sequenceID++
				m.metrics.Put(nameBytes, int(m.sequenceID))
				return m.sequenceID, true
			}
			return NotFoundMetricID, false
		}
		return id, true
	}
	return NotFoundMetricID, false
}

//SuggestMetrics returns suggestions of metric names given a search prefix.
func (m *MetricUID) SuggestMetrics(prefix string, limit int16) map[string]struct{} {
	if len(prefix) > 0 {
		nameBytes := []byte(prefix)
		partition := getPartition(nameBytes)

		m.family.Lookup(partition, func(bytes []byte) bool {
			treeReader := tree.NewReader(bytes)
			it := treeReader.Seek(nameBytes)
			if nil != it {
				//todo
				logger.GetLogger("tsdb/index").Error("", m.dbField)
			}
			return true
		})
	} else {
		//todo
		logger.GetLogger("tsdb/index").Warn("", m.dbField)
		//partitions := m.getSortPartition()
	}
	return nil
}

//Flush represents forces a flush of in-memory data, and clear it
func (m *MetricUID) Flush() error {
	if m.metrics.Len() == 0 {
		return nil
	}
	flusher := m.family.NewFlusher()
	bs := util.Uint32ToBytes(m.sequenceID)
	err := flusher.Add(MetricSequenceIDKey, bs)
	if nil != err {
		logger.GetLogger("tsdb/index").Error("write metric sequenceID error!", m.dbField, logger.Error(err))
		return err
	}

	writer := tree.NewWriter(m.metrics)
	byteArray, err := writer.Encode()
	if nil != err {
		logger.GetLogger("tsdb/index").Error(" metricTree encode error!", m.dbField, logger.Error(err))
	}
	err = flusher.Add(m.partition, byteArray)
	if nil != err {
		logger.GetLogger("tsdb/index").Error("write metric tree error!",
			m.dbField, zap.String("partition", string(m.partition)))
		return err
	}
	m.metrics.Clear()

	return flusher.Commit()
}

// getMetricIdFromDisk return unique int32 id, return -1 if not found
func (m *MetricUID) getMetricIDFromDisk(partition uint32, metric []byte) uint32 {
	var metricID = NotFoundMetricID
	m.family.Lookup(partition, func(bytes []byte) bool {
		treeReader := tree.NewReader(bytes)
		v, ok := treeReader.Get(metric)
		if ok {
			metricID = uint32(v)
			return true
		}
		return false
	})
	return metricID
}

//getPartition returns determine partition according to the first byte
func getPartition(name []byte) uint32 {
	return uint32(name[0])
}
