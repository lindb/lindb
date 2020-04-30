package wal

import (
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/series/field"
)

//go:generate mockgen -source=./metric_meta_wal.go -destination=./metric_meta_wal_mock.go -package=wal

var metaWAlLogger = logger.GetLogger("wal", "meta")

const (
	metricMetaPageSize = 64 * 1024 * 1024 // 64M
	// type(1 byte) + ns length (1 byte) + metric length (1 byte) + metric id (4 bytes)
	metricBaseLength = 1 + 1 + 1 + 4
	// type(1 byte) + field id (1 byte) + field type (1 byte) + field length (1 byte) + metric id (4 bytes)
	fieldBaseLength = 1 + 1 + 1 + 1 + 4
	// type(1 byte) + tag key length (1 byte) + metric id (4 bytes) + tag key id (4 bytes)
	tagKeyBaseLength = 1 + 1 + 4 + 4
)

// metaType represents meta type
type metaType uint8

// Defines all meta types
const (
	metricType metaType = iota + 1
	fieldType
	tagKeyType
)

// MetricMetaWAL represents write ahead log which stores metric metadata for meta database
type MetricMetaWAL interface {
	// AppendMetric appends namespace/metricName/metricID into wal log
	AppendMetric(namespace, metricName string, metricID uint32) error
	// AppendField appends metricID/fieldID/fieldName/fieldType into wal log
	AppendField(metricID uint32, fID field.ID, fieldName string, fType field.Type) error
	// AppendTagKey appends metricID/tagKeyID/tagKey into wal log
	AppendTagKey(metricID uint32, tagKeyID uint32, tagKey string) error
	// NeedRecovery checks if wal log need to recover
	NeedRecovery() bool
	// Recovery recoveries wal log, then writes data via recovery function
	Recovery(metricRecovery MetricRecoveryFunc,
		fieldRecovery FieldRecoveryFunc,
		tagKeyRecovery TagKeyRecoveryFunc,
		commit CommitFunc)
	// Sync flushes data into disk
	Sync() error
	// Close closes the wal log
	Close() error
}

// metricMetaWAL implements MetricMetaWAL interface
type metricMetaWAL struct {
	base *baseWAL
}

// NewMetricMetaWAL creates a new metric meta write ahead log
func NewMetricMetaWAL(path string) (MetricMetaWAL, error) {
	base, err := newBaseWAL(path, metricMetaPageSize)
	if err != nil {
		return nil, err
	}
	return &metricMetaWAL{base: base}, nil
}

// AppendMetric appends namespace/metricName/metricID into wal log
func (m *metricMetaWAL) AppendMetric(namespace, metricName string, metricID uint32) error {
	if err := m.base.checkPage(len(namespace) + len(metricName) + metricBaseLength); err != nil {
		return err
	}
	m.base.putUint8(uint8(metricType))
	m.base.putString(namespace)
	m.base.putString(metricName)
	m.base.putUint32(metricID)
	return nil
}

// AppendField appends metricID/fieldID/fieldName/fieldType into wal log
func (m *metricMetaWAL) AppendField(metricID uint32, fID field.ID, fieldName string, fType field.Type) error {
	if err := m.base.checkPage(len(fieldName) + fieldBaseLength); err != nil {
		return err
	}
	m.base.putUint8(uint8(fieldType))
	m.base.putUint32(metricID)
	m.base.putUint8(uint8(fID))
	m.base.putString(fieldName)
	m.base.putUint8(uint8(fType))
	return nil
}

// AppendTagKey appends metricID/tagKeyID/tagKey into wal log
func (m *metricMetaWAL) AppendTagKey(metricID uint32, tagKeyID uint32, tagKey string) error {
	if err := m.base.checkPage(len(tagKey) + tagKeyBaseLength); err != nil {
		return err
	}
	m.base.putUint8(uint8(tagKeyType))
	m.base.putUint32(metricID)
	m.base.putUint32(tagKeyID)
	m.base.putString(tagKey)
	return nil
}

// NeedRecovery checks if wal log need to recover
func (m *metricMetaWAL) NeedRecovery() bool {
	return m.base.needRecovery()
}

// Recovery recoveries wal log, then writes data via recovery function
func (m *metricMetaWAL) Recovery(metricRecovery MetricRecoveryFunc,
	fieldRecovery FieldRecoveryFunc,
	tagKeyRecovery TagKeyRecoveryFunc,
	commit CommitFunc) {
	current := m.base.pageIndex.Load()
	committed := m.base.commitPageIndex.Load()
	for i := committed; i < current; i++ {
		walPage, ok := m.base.walFactory.GetPage(i)
		if !ok {
			continue
		}
		offset := 0
		completed := false
		for !completed {
			mType := metaType(walPage.ReadUint8(offset))
			offset++
			switch mType {
			case metricType: // recovery metric
				ns, n := readString(walPage, offset)
				offset += n
				metricName, n := readString(walPage, offset)
				offset += n
				metricID := walPage.ReadUint32(offset)
				offset += 4
				if err := metricRecovery(ns, metricName, metricID); err != nil {
					//TODO add metric?????
					metaWAlLogger.Error("invoke metric recovery func error",
						logger.String("wal", m.base.path), logger.Error(err))
					return
				}
			case fieldType: // recovery field
				metricID := walPage.ReadUint32(offset)
				offset += 4
				fID := walPage.ReadUint8(offset)
				offset++
				fieldName, n := readString(walPage, offset)
				offset += n
				fType := walPage.ReadUint8(offset)
				offset++
				if err := fieldRecovery(metricID, field.ID(fID), fieldName, field.Type(fType)); err != nil {
					//TODO add metric?????
					metaWAlLogger.Error("invoke field recovery func error",
						logger.String("wal", m.base.path), logger.Error(err))
					return
				}
			case tagKeyType: // recovery tag key
				metricID := walPage.ReadUint32(offset)
				offset += 4
				tagKeyID := walPage.ReadUint32(offset)
				offset += 4
				tagKey, n := readString(walPage, offset)
				offset += n
				if err := tagKeyRecovery(metricID, tagKeyID, tagKey); err != nil {
					//TODO add metric?????
					metaWAlLogger.Error("invoke tag key recovery func error",
						logger.String("wal", m.base.path), logger.Error(err))
					return
				}
			default:
				completed = true // no data
			}
		}

		if err := commit(); err != nil {
			//TODO add metric?????
			metaWAlLogger.Error("invoke commit func error",
				logger.String("wal", m.base.path), logger.Error(err))
			return
		}

		if err := m.base.walFactory.ReleasePage(i); err != nil {
			//TODO add metric?????
			metaWAlLogger.Error("release meta wal page error",
				logger.String("wal", m.base.path), logger.Error(err))
		}

		m.base.commitPageIndex.Inc()
	}
}

// Sync flushes metric meta into disk
func (m *metricMetaWAL) Sync() error {
	return m.base.sync()
}

// Close closes the wal log
func (m *metricMetaWAL) Close() error {
	return m.base.close()
}
