package indexdb

import (
	"bytes"
	"compress/gzip"
	"encoding"
	"encoding/binary"
	"io"
	"io/ioutil"

	"github.com/lindb/lindb/pkg/stream"

	art "github.com/plar/go-adaptive-radix-tree"
)

//go:generate mockgen -source ./art_tree.go -destination=./art_tree_mock.go -package=indexdb

// artTreeINTF is a serializable/deserializable Adaptive-Radix-Tree
type artTreeINTF interface {
	encoding.BinaryUnmarshaler
	art.Tree
}

// artTree implements artTreeINTF
type artTree struct {
	art.Tree
}

// newArtTree returns a new ART-Tree
func newArtTree() artTreeINTF {
	return &artTree{art.New()}
}

// UnmarshalBinary set the tree from the binary.
func (tree *artTree) UnmarshalBinary(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	decoded, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return err
	}
	reader := stream.BinaryReader(decoded)
	for {
		// read length of metricName
		size := reader.ReadUvarint64()
		if reader.Error() != nil {
			break
		}
		metricName := reader.ReadBytes(int(size))
		metricID := reader.ReadUint32()
		if reader.Error() != nil {
			break
		}
		tree.Insert(art.Key(metricName), metricID)
	}
	if reader.Error() != nil && reader.Error() != io.EOF {
		return reader.Error()
	}
	return nil
}

// nameIDCompressor is used to compress newly-created metricName and metricID pairs.
type nameIDCompressor struct {
	buf         bytes.Buffer // storing metric-Name, metricID
	variableBuf [8]byte      // placeholder for uint32
	gzipWriter  *gzip.Writer
}

// newNameIDCompressor returns a new nameIDCompressor
func newNameIDCompressor() *nameIDCompressor {
	compressor := &nameIDCompressor{}
	compressor.gzipWriter, _ = gzip.NewWriterLevel(&compressor.buf, gzip.BestSpeed)
	return compressor
}

// AddNameID add a new metricName and metricID pair to buffer
func (c *nameIDCompressor) AddNameID(metricName string, metricID uint32) {
	// write metricName length
	size := binary.PutUvarint(c.variableBuf[:], uint64(len(metricName)))
	_, _ = c.gzipWriter.Write(c.variableBuf[:size])
	// write metricName
	_, _ = c.gzipWriter.Write([]byte(metricName))
	// write metricID
	binary.BigEndian.PutUint32(c.variableBuf[:], metricID)
	_, _ = c.gzipWriter.Write(c.variableBuf[:4])
}

// Close closes the underlying gzip writer and return compressed data
func (c *nameIDCompressor) Close() ([]byte, error) {
	if err := c.gzipWriter.Close(); err != nil {
		return nil, err
	}
	return c.buf.Bytes(), nil
}
