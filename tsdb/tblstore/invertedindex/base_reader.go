package invertedindex

import (
	"fmt"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
)

const (
	indexFooterSize = 4 + // keys position
		4 + // offsets position
		4 // crc32 checksum
)

// baseReader represents the base index reader, include basic reader context
type baseReader struct {
	buf           []byte
	offsets       *encoding.FixedOffsetDecoder
	keys          *roaring.Bitmap
	crc32CheckSum uint32
}

// initReader initializes the basic index reader context
func (r *baseReader) initReader() error {
	if len(r.buf) <= indexFooterSize {
		return fmt.Errorf("block length no ok")
	}
	// read footer(4+4+4)
	footerPos := len(r.buf) - indexFooterSize
	keysStartPos := int(stream.ReadUint32(r.buf, footerPos))
	offsetsPos := int(stream.ReadUint32(r.buf, footerPos+4))
	r.crc32CheckSum = stream.ReadUint32(r.buf, footerPos+8)
	// validate offsets
	if keysStartPos > footerPos || offsetsPos > keysStartPos {
		return fmt.Errorf("bad offsets")
	}
	// read keys
	keys := roaring.New()
	if err := encoding.BitmapUnmarshal(keys, r.buf[keysStartPos:]); err != nil {
		return err
	}
	r.keys = keys
	// read high keys offsets
	r.offsets = encoding.NewFixedOffsetDecoder(r.buf[offsetsPos:])
	return nil
}
