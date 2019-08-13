package indexdb

import (
	"bytes"
	"compress/flate"
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
	encoding.BinaryMarshaler
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

// MarshalBinary encodes the artTree into a binary form and returns the result.
func (tree *artTree) MarshalBinary() (data []byte, err error) {
	var (
		buffer        bytes.Buffer // storing metric-Name, metricID
		variableBuf   [8]byte      // placeholder for uint32
		gzipWriter, _ = gzip.NewWriterLevel(&buffer, flate.BestSpeed)
	)
	for it := tree.Iterator(); it.HasNext(); {
		item, _ := it.Next()
		metricName := item.Key()
		metricID, ok := item.Value().(uint32)
		if !ok {
			indexDBLogger.Error("ART-Tree node type error")
			continue
		}
		// write metricName length
		size := binary.PutUvarint(variableBuf[:], uint64(len(metricName)))
		_, _ = gzipWriter.Write(variableBuf[:size])
		// write metricName
		_, _ = gzipWriter.Write([]byte(metricName))
		// write metricID
		binary.BigEndian.PutUint32(variableBuf[:], metricID)
		_, _ = gzipWriter.Write(variableBuf[:4])
	}
	if err = gzipWriter.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
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
		if reader.Error() != nil {
			break
		}
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
