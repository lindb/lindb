package table

import (
	"fmt"

	"github.com/eleme/lindb/pkg/encoding"
	"github.com/eleme/lindb/pkg/io"
	"github.com/eleme/lindb/pkg/logger"

	"github.com/RoaringBitmap/roaring"
	"go.uber.org/zap"
)

type Builder interface {
	Add(key uint32, value []byte) bool
	Close() error
}

type StoreBuilder struct {
	keys    *roaring.Bitmap
	lastKey uint32
	logger  *zap.Logger
	fw      *io.FileWriter
	offset  *encoding.DeltaBitPackingEncoder

	pos int32

	first bool
}

func NewStoreBuilder(fileName string) (Builder, error) {
	keys := roaring.New()
	log := logger.GetLogger()
	writer, err := io.NewWriter(fileName)
	if err != nil {
		return nil, fmt.Errorf("create file write for store builder error:%s", err)
	}
	return &StoreBuilder{
		keys:   keys,
		logger: log,
		fw:     writer,
		first:  true,
		pos:    0,
		offset: encoding.NewDeltaBitPackingEncoder(),
	}, nil
}

func (b *StoreBuilder) Add(key uint32, value []byte) bool {
	if !b.first && key <= b.lastKey {
		b.logger.Warn("key is smaller then last key ignore current options.",
			zap.Any("last", b.lastKey), zap.Any("cur", key))
		return false
	}

	n, err := b.fw.Write(value)
	if err != nil {
		b.pos += int32(n)
		//TODO
		b.logger.Error("write file error")
		return false
	}
	// get write pos
	pos := b.pos
	// add pos into offset
	b.offset.Add(pos)

	b.pos = pos + int32(n)
	// add key into index block
	b.keys.Add(key)

	b.lastKey = key
	b.first = false

	return true
}

func (b *StoreBuilder) Close() error {
	offset, err := b.offset.Bytes()
	if err != nil {
		return fmt.Errorf("marshal store table offsets error:%s", err)
	}

	n, err := b.fw.Write(offset)
	if err != nil {
		return fmt.Errorf("write offsets error:%s", err)
	}

	b.keys.RunOptimize()
	keys, err := b.keys.MarshalBinary()
	if err != nil {
		return err
	}

	b.pos += int32(n)

	n, err = b.fw.Write(keys)
	if err != nil {
		return err
	}

	b.pos += int32(n)

	return nil
}
