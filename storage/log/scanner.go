package log

import (
	"encoding/binary"

	"github.com/apache/arrow-go/v18/arrow"
	logspkg "github.com/lindb/arrow/pkg/logs"
	"github.com/lindb/roaring"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/storage/store"
)

type Scanner interface {
	HasNext() bool
	Next(fn func(reader *logspkg.Reader, rowNum int)) error
	Close()
}

type scanner struct {
	segment *Segment

	logIDs *roaring.Bitmap
	it     roaring.IntPeekable

	wal map[models.NodeID]store.WriteAheadLog // leader => write ahead log

	current struct {
		leader   models.NodeID
		sequence uint32

		data       []byte
		logs       arrow.RecordBatch
		attributes arrow.RecordBatch
	}

	reader    *logspkg.BinaryReader
	logReader *logspkg.Reader
}

func NewScanner(segment *Segment, logIDs *roaring.Bitmap) Scanner {
	return &scanner{
		segment: segment,
		logIDs:  logIDs,
		it:      logIDs.Iterator(),
		wal:     segment.GetWALs(),
	}
}

func (s *scanner) HasNext() bool {
	return s.it.HasNext()
}

func (s *scanner) Next(fn func(reader *logspkg.Reader, rowNum int)) error {
	logID := s.it.Next()
	index, err := s.segment.GetIndex(logID)
	if err != nil {
		return err
	}
	leader := models.NodeID(index[0])
	sequence := binary.LittleEndian.Uint32(index[1:])
	rowNum := binary.LittleEndian.Uint32(index[5:])
	if leader != s.current.leader || sequence != s.current.sequence || len(s.current.data) == 0 {
		// update current log info, if leader or sequence changed, read log data from wal
		s.current.leader = leader
		s.current.sequence = sequence

		// read log data from wal
		data, err := s.wal[leader].Get(int64(sequence))
		if err != nil {
			return err
		}
		s.current.data = data

		if err := s.initializeReader(); err != nil {
			return err
		}
	}

	fn(s.logReader, int(rowNum))

	return nil
}

func (s *scanner) initializeReader() (err error) {
	var logs, attributes arrow.RecordBatch
	if s.logReader == nil {
		reader, err := logspkg.NewBinaryReader()
		if err != nil {
			return err
		}
		s.reader = reader
		logs, attributes, err = reader.ReadFrom(s.current.data)
		if err != nil {
			return err
		}
		s.logReader = logspkg.NewReader(logs, attributes)
	} else {
		// release previous log and attribute record batch before read new data
		s.current.logs.Release()
		s.current.attributes.Release()

		logs, attributes, err = s.reader.ReadFrom(s.current.data)
		if err != nil {
			return err
		}
		s.logReader.Reset(logs, attributes)
	}

	// update current log and attribute record batch
	s.current.logs = logs
	s.current.attributes = attributes
	return nil
}

func (s *scanner) Close() {
	if s.reader != nil {
		s.current.logs.Release()
		s.current.attributes.Release()

		s.reader.Release()
	}
}
