package segment

import (
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/queue/page"
)

const (
	indexFileSuffix = ".idx"
	dataFileSuffix  = ".dat"
)

// ErrSegmentNotFound represents error when no suitable segment is found.
var ErrSegmentNotFound = errors.New("segment not found")

// SeqRange represents a sequence of begin sequence in ascending order.
type SeqRange []int64

// implemented interface for sort
func (sr SeqRange) Len() int           { return len(sr) }
func (sr SeqRange) Swap(i, j int)      { sr[i], sr[j] = sr[j], sr[i] }
func (sr SeqRange) Less(i, j int) bool { return sr[i] < sr[j] }

// Last returns the last(the largest seq) of a SeqRange.
func (sr SeqRange) Last() int64 {
	return sr[len(sr)-1]
}

// Find returns a tuple (index, seq, ok).
// index: the index of seq in SeqRange.
// seq: the sequence num.
// ok: true means found, false means not found.
// index and seq satisfy SeqRange[index] <= seq and (SeqRange[index + 1] > seq or index == SeqRange.Len() - 1).
// SeqRange should be sorted before this call.
func (sr SeqRange) Find(seq int64) (int, int64, bool) {
	if sr.Len() == 0 || seq < sr[0] {
		return 0, 0, false
	}

	// binary search
	l, r := 0, sr.Len()-1

	for l <= r {
		m := (l + r) >> 1
		if sr[m] <= seq {
			if m == r || sr[m+1] > seq {
				return m, sr[m], true
			}
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return 0, 0, false
}

// Factory manages loading, retrieving and creating segment.
type Factory interface {
	// GetSegment returns a segment contains seq.
	GetSegment(seq int64) (Segment, error)
	// NewSegment creates a segment with given beginSeq.
	NewSegment(begSeq int64) (Segment, error)
	// RemoveSegments removes segments with tailSeq <= ackSeq
	RemoveSegments(ackSeq int64) error
	// SegmentsSize returns segments size hold in factory, mainly for test
	SegmentsSize() int
	// Close closes the segments.
	Close()
}

// factory implements Factory, loading segments from file when start up.
type factory struct {
	// dirPath for segment files
	dirPath string
	// the max size limit in bytes for data file
	dataFileSizeLimit int
	// segments in ascending order
	segments []Segment
	// segments beg sequence slice
	seqRange SeqRange
	// lock for segments
	lock4segments sync.RWMutex
	logger        *logger.Logger
}

// NewFactory builds a segment factory by loading file from dirPath.
// HeadSeq and  TailSeq are used to filter segments in use.
func NewFactory(dirPath string, dataFileSizeLimit int, headSeq, tailSeq int64) (Factory, error) {
	if err := fileutil.MkDir(dirPath); err != nil {
		return nil, err
	}

	fct := &factory{
		dirPath:           dirPath,
		dataFileSizeLimit: dataFileSizeLimit,
		segments:          make([]Segment, 0),
		seqRange:          make(SeqRange, 0),
		logger:            logger.GetLogger("pkg/queue", "SegmentFactory"),
	}

	// load from file
	if err := fct.load(headSeq, tailSeq); err != nil {
		return nil, err
	}

	return fct, nil
}

// load loads segments from dirPath, filtering segment file with headSeq and tailSeq.
func (fct *factory) load(headSeq, tailSeq int64) error {
	fileNames, err := fileutil.ListDir(fct.dirPath)
	if err != nil {
		return err
	}

	// empty
	if len(fileNames) == 0 {
		return nil
	}

	// dirPath should only contains .idx and .dat files
	seqRange := make(SeqRange, 0, len(fileNames)/2)

	// file set, check both .idx and .dat files exist
	filePathSet := make(map[string]struct{})

	for _, fn := range fileNames {
		filePath := path.Join(fct.dirPath, fn)
		filePathSet[filePath] = struct{}{}
		if strings.HasSuffix(fn, indexFileSuffix) {
			seqNumStr := fn[0:strings.Index(fn, indexFileSuffix)]
			seq, err := strconv.ParseInt(seqNumStr, 10, 64)
			if err != nil {
				return err
			}
			seqRange = append(seqRange, seq)
		}
	}

	// ensure segments corresponding to headSeq and tailSeq exist
	sort.Sort(seqRange)
	hi, _, ok := seqRange.Find(headSeq)
	if !ok {
		return fmt.Errorf("dirPath:%s, segment file for headSeq:%d not found", fct.dirPath, headSeq)
	}

	ti, _, ok := seqRange.Find(tailSeq)
	if !ok {
		return fmt.Errorf("dirPath:%s, segment file for tailSeq:%d not found", fct.dirPath, tailSeq)
	}

	for i := ti; i <= hi; i++ {
		begin := seqRange[i]

		end := headSeq

		if i != hi {
			end = seqRange[i+1]
		}

		if err := fct.loadOrCreateSegment(begin, end, filePathSet); err != nil {
			return err
		}

	}

	return nil
}

// loadOrCreateSegment loads a segment from {begin}.idx, {begin}.dat if both of them exist,
// if not, creates corresponding files and load.
func (fct *factory) loadOrCreateSegment(begin, end int64, filePathSet map[string]struct{}) error {
	indexFilePath, dataFilePath := fct.buildIndexAndDataFilePath(begin)
	// .idx files has been checked before
	_, ok := filePathSet[dataFilePath]
	if !ok {
		return fmt.Errorf("segemnt file %s not found", dataFilePath)
	}

	dataMappedBytes, err := fileutil.RWMap(dataFilePath, fct.dataFileSizeLimit)
	if err != nil {
		return err
	}

	// the worse case, all messages in a dataFile is one byte, one message takes 8 bytes for index
	// don't worry about the disk usage since init size doesn't occupy real disk storage,
	// and the index file will be truncated to proper size when close.
	indexMappedBytes, err := fileutil.RWMap(indexFilePath, indexItemSize*fct.dataFileSizeLimit)
	if err != nil {
		return err
	}

	dataMappedPage := page.NewMappedPage(dataFilePath, dataMappedBytes, page.MMapCloseFunc, page.MMapSyncFunc)
	indexMappedPage := page.NewMappedPage(indexFilePath, indexMappedBytes, page.MMapCloseFunc, page.MMapSyncFunc)

	seg, err := NewSegment(indexMappedPage, dataMappedPage, begin, end)
	if err != nil {
		return err
	}
	fct.segments = append(fct.segments, seg)
	fct.seqRange = append(fct.seqRange, begin)

	return nil
}

// buildFilePath concatenates the dirPath and fileName as a filePath
func (fct *factory) buildFilePath(fileName string) string {
	return path.Join(fct.dirPath, fileName)
}

// buildIndexAndDataFilePath returns the indexFilePath and dataFilePath for segment with beginSeq.
func (fct *factory) buildIndexAndDataFilePath(beginSeq int64) (indexFilePath, dataFilePath string) {
	seqNumStr := strconv.FormatInt(beginSeq, 10)

	dataFileName := seqNumStr + dataFileSuffix
	indexFileName := seqNumStr + indexFileSuffix
	dataFilePath = fct.buildFilePath(dataFileName)
	indexFilePath = fct.buildFilePath(indexFileName)
	return
}

// GetSegment returns a segment contains seq.
func (fct *factory) GetSegment(seq int64) (Segment, error) {
	fct.lock4segments.RLock()
	defer fct.lock4segments.RUnlock()
	ix, _, ok := fct.seqRange.Find(seq)
	if ok {
		return fct.segments[ix], nil
	}
	return nil, ErrSegmentNotFound
}

// NewSegment creates a segment with given beginSeq.
func (fct *factory) NewSegment(beginSeq int64) (Segment, error) {
	fct.lock4segments.Lock()
	defer fct.lock4segments.Unlock()
	if fct.seqRange.Len() > 0 {
		lastSeq := fct.seqRange.Last()
		if beginSeq <= lastSeq {
			return nil, errors.New("new segment beginSeq should be larger than exists segments")
		}
	}

	fakeSet := map[string]struct{}{
		fct.buildFilePath(strconv.FormatInt(beginSeq, 10) + dataFileSuffix): {},
	}

	err := fct.loadOrCreateSegment(beginSeq, beginSeq, fakeSet)
	if err != nil {
		return nil, err
	}

	return fct.segments[len(fct.segments)-1], nil
}

// RemoveSegments removes segments with tailSeq <= ackSeq, removing files error will only be logged, not returned.
func (fct *factory) RemoveSegments(ackSeq int64) error {
	fct.lock4segments.Lock()
	if fct.seqRange.Len() == 0 {
		fct.lock4segments.Unlock()
		return nil
	}

	lastSeg := fct.segments[len(fct.segments)-1]
	if lastSeg.End() <= ackSeq {
		fct.lock4segments.Unlock()
		return nil
	}

	index := 0
	for _, seg := range fct.segments {
		if seg.End() > ackSeq {
			break
		}
		index++
	}

	ackedSegments := fct.segments[0:index]

	fct.segments = fct.segments[index:]
	fct.seqRange = fct.seqRange[index:]

	fct.lock4segments.Unlock()

	for _, seg := range ackedSegments {

		fct.logger.Info("remove segment",
			logger.String("dirPath", fct.dirPath),
			logger.Int64("begin", seg.Begin()),
			logger.Int64("end", seg.End()))

		seg.Close()

		indexFilePath, dataFilePath := fct.buildIndexAndDataFilePath(seg.Begin())
		if fileutil.Exist(indexFilePath) {
			if err := os.Remove(indexFilePath); err != nil {
				fct.logger.Error("error rm indexFile:"+indexFilePath, logger.String("dirPath", fct.dirPath), logger.Error(err))
			}
		}
		if fileutil.Exist(dataFilePath) {
			if err := os.Remove(dataFilePath); err != nil {
				fct.logger.Error("error rm indexFile:"+indexFilePath, logger.String("dirPath", fct.dirPath), logger.Error(err))
			}
		}
	}
	return nil
}

// SegmentsSize returns segments size hold in factory, mainly for test
func (fct *factory) SegmentsSize() int {
	fct.lock4segments.RLock()
	defer fct.lock4segments.RUnlock()
	return fct.seqRange.Len()
}

// Close closes the segments.
func (fct *factory) Close() {
	fct.lock4segments.RLock()
	defer fct.lock4segments.RUnlock()
	for _, seg := range fct.segments {
		seg.Close()
	}
}
