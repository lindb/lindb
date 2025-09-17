package wal

import (
	"errors"
	"fmt"
	"io"

	"go.uber.org/atomic"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/queue"
)

type WriteAheadLog interface {
	io.Closer

	Get(index int64) ([]byte, error)

	Write(msg []byte) error

	Replica(replicaIndex int64, msg []byte) (int64, error)
	// ReplicaAckIndex returns the index which replica appended index.
	ReplicaAckIndex() int64
	// ResetReplicaIndex resets replica index.
	ResetReplicaIndex(idx int64)

	Consume() (int64, []byte, error)
}

type writeAheadLog struct {
	data  queue.FanOutQueue
	local queue.ConsumerGroup

	closed atomic.Bool
}

func NewWriteAheadLog(path string) (WriteAheadLog, error) {
	// TODO: database level???
	pageSize := config.GlobalStorageConfig().WAL.PageSize
	data, err := queue.NewFanOutQueue(path, int64(pageSize))
	if err != nil {
		return nil, err
	}
	fmt.Printf("wal===%s=%s\n", path, err)
	local, err := data.GetOrCreateConsumerGroup("local")
	if err != nil {
		return nil, err
	}
	return &writeAheadLog{
		data:  data,
		local: local,
	}, nil
}

func (w *writeAheadLog) Replica(replicaIndex int64, msg []byte) (int64, error) {
	if w.closed.Load() {
		return -1, errors.New("write ahead log is closed")
	}
	appendIdx := w.data.Queue().AppendedSeq() + 1
	if replicaIndex != appendIdx {
		return appendIdx, nil
	}
	if err := w.data.Queue().Put(msg); err != nil {
		return -1, err
	}
	return appendIdx, nil
}

func (w *writeAheadLog) Consume() (int64, []byte, error) {
	seq := w.local.Consume()
	if seq >= 0 {
		data, err := w.data.Queue().Get(seq)
		if err != nil {
			return 0, nil, err
		}
		return seq, data, nil
	}
	return 0, nil, nil
}

func (w *writeAheadLog) ReplicaAckIndex() int64 {
	return w.data.Queue().AppendedSeq()
}

func (w *writeAheadLog) ResetReplicaIndex(idx int64) {
	w.data.SetAppendedSeq(idx - 1)
}

func (w *writeAheadLog) Get(index int64) ([]byte, error) {
	return w.data.Queue().Get(index)
}

func (w *writeAheadLog) Write(msg []byte) error {
	if len(msg) == 0 {
		return nil
	}
	if w.closed.Load() {
		return errors.New("write ahead log is closed")
	}
	// TODO: add metric
	return w.data.Queue().Put(msg)
}

func (w *writeAheadLog) Close() error {
	if w.closed.CompareAndSwap(false, true) {
		// close queue
		w.data.Close()
	}
	return nil
}
