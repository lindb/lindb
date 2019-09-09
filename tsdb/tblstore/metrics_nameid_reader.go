package tblstore

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/pkg/stream"

	art "github.com/plar/go-adaptive-radix-tree"
)

//go:generate mockgen -source ./metrics_nameid_reader.go -destination=./metrics_nameid_reader_mock.go -package tblstore

const (
	metricNameIDSequenceSize = 4 + // metricID sequence
		4 // tagKeyID sequence
)

// MetricsNameIDReader reads metricNameID info from the kv table
type MetricsNameIDReader interface {
	// ReadMetricNS read metricNameID data by the namespace-id
	ReadMetricNS(nsID uint32) (data [][]byte, metricIDSeq, tagKeyIDSeq uint32, ok bool)
	// UnmarshalBinaryToART de-compresses the compressed block, then insert the metricName-id pair to the tree
	UnmarshalBinaryToART(tree art.Tree, data []byte) error
}

// metricsNameIDReader implements MetricsNameIDReader
type metricsNameIDReader struct {
	gzipReader *gzip.Reader
	sr         *stream.Reader
	readers    []table.Reader
}

// NewMetricsNameIDReader returns a new MetricsNameIDReader
func NewMetricsNameIDReader(readers []table.Reader) MetricsNameIDReader {
	return &metricsNameIDReader{
		sr:      stream.NewReader(nil),
		readers: readers}
}

// UnmarshalBinaryToART de-compresses the compressed block, then insert the metricName-id pair to the tree
func (r *metricsNameIDReader) UnmarshalBinaryToART(
	tree art.Tree,
	data []byte,
) error {
	decompressed, err := r.DeCompress(data)
	if err != nil {
		return err
	}
	r.sr.Reset(decompressed)

	for !r.sr.Empty() {
		// read length of metricName
		size := r.sr.ReadUvarint64()
		metricName := r.sr.ReadBytes(int(size))
		metricID := r.sr.ReadUint32()
		if r.sr.Error() != nil {
			return r.sr.Error()
		}
		tree.Insert(art.Key(metricName), metricID)
	}
	return nil
}

// DeCompress decompresses the compressed block
func (r *metricsNameIDReader) DeCompress(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, nil
	}
	if r.gzipReader == nil {
		gzipReader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		r.gzipReader = gzipReader
	}
	if err := r.gzipReader.Reset(bytes.NewReader(data)); err != nil {
		return nil, err
	}
	defer r.gzipReader.Close()
	decompressed, err := ioutil.ReadAll(r.gzipReader)
	return decompressed, err
}

// ReadMetricNS read metricNameID data by the namespace-id
func (r *metricsNameIDReader) ReadMetricNS(
	nsID uint32,
) (
	dataList [][]byte,
	maxMetricIDSeq,
	maxTagKeyIDSeq uint32,
	ok bool,
) {
	for _, reader := range r.readers {
		block := reader.Get(nsID)
		data, metricIDSeq, tagKeyIDSeq, thisOK := r.ReadBlock(block)
		if !thisOK {
			continue
		}
		dataList = append(dataList, data)
		maxMetricIDSeq = metricIDSeq
		maxTagKeyIDSeq = tagKeyIDSeq
		ok = true
	}
	return
}

// ReadBlock splits the block into sequence-part and compressed-part
func (r *metricsNameIDReader) ReadBlock(
	block []byte,
) (
	compressed []byte,
	metricIDSeq,
	tagKeyIDSeq uint32,
	ok bool,
) {
	if len(block) < metricNameIDSequenceSize {
		return nil, 0, 0, false
	}
	idSequencePos := uint32(len(block) - metricNameIDSequenceSize)
	compressed = block[:idSequencePos]
	r.sr.Reset(block)
	r.sr.ShiftAt(idSequencePos)

	metricIDSeq = r.sr.ReadUint32()
	tagKeyIDSeq = r.sr.ReadUint32()
	return compressed, metricIDSeq, tagKeyIDSeq, true
}
