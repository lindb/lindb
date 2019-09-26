package stream_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/stream"
)

func Test_Stream_ReaderWriter(t *testing.T) {
	writer1 := stream.NewBufferWriter(nil)
	assert.NotNil(t, writer1)

	var buf bytes.Buffer
	writer2 := stream.NewBufferWriter(&buf)

	writer2.PutUint64(1)
	writer2.PutUint32(2)
	writer2.PutInt32(-3)
	writer2.PutInt64(-4)
	writer2.PutByte(5)
	writer2.PutBytes([]byte{6, 7, 8})
	writer2.PutUvarint64(1234567890)
	writer2.PutVarint64(-1234567890)
	writer2.PutUvarint32(12345)
	writer2.PutVarint32(-12345)
	writer2.PutUInt16(1)
	writer2.PutInt16(-2)
	assert.Nil(t, writer2.Error())

	data, err := writer2.Bytes()
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, 47, writer2.Len())
	writer2.Reset()
	writer2.SwitchBuffer(bytes.NewBuffer(nil))
	assert.Equal(t, 0, writer2.Len())

	reader := stream.NewReader(data)
	assert.Equal(t, uint64(1), reader.ReadUint64())
	assert.Equal(t, uint32(2), reader.ReadUint32())
	assert.Equal(t, int32(-3), reader.ReadInt32())
	assert.Equal(t, int64(-4), reader.ReadInt64())
	assert.Equal(t, byte(5), reader.ReadByte())
	assert.Equal(t, []byte{6, 7, 8}, reader.ReadBytes(3))
	assert.Equal(t, uint64(1234567890), reader.ReadUvarint64())
	assert.Equal(t, int64(-1234567890), reader.ReadVarint64())
	assert.Equal(t, uint32(12345), reader.ReadUvarint32())
	assert.Equal(t, int32(-12345), reader.ReadVarint32())
	assert.Equal(t, uint16(1), reader.ReadUint16())
	assert.Equal(t, int16(-2), reader.ReadInt16())

	assert.Nil(t, reader.Error())
	assert.True(t, reader.Empty())

	reader.ReadBytes(1)
	assert.NotNil(t, reader.Error())
	assert.True(t, reader.Empty())

	// shift test
	reader.Reset(data)
	assert.Nil(t, reader.Error())
	assert.False(t, reader.Empty())
	_ = reader.ReadSlice(8)
	_ = reader.ReadSlice(0)
	_ = reader.ReadSlice(0)
	assert.Equal(t, uint32(2), reader.ReadUint32())
	// 12 bytes
	assert.Equal(t, 12, reader.Position())
	_ = reader.ReadSlice(35)
	assert.Nil(t, reader.Error())
	_ = reader.ReadSlice(1)
	assert.NotNil(t, reader.Error())

	// read failure
	assert.Zero(t, reader.ReadUint64())
	assert.Zero(t, reader.ReadUint32())
	assert.Zero(t, reader.ReadUint16())
}

func Test_Stream_SliceWriter(t *testing.T) {
	w := stream.NewSliceWriter(nil)
	assert.Nil(t, w.Error())
	w.PutByte(byte(3))
	assert.NotNil(t, w.Error())

	_, err := w.Bytes()
	assert.NotNil(t, err)
}

func Test_Reader_ReadSlice(t *testing.T) {
	sl := make([]byte, 200)
	reader := stream.NewReader(nil)

	for x := 0; x < 2; x++ {
		reader.Reset(sl)

		assert.Zero(t, reader.Position())
		data := reader.ReadSlice(100)
		assert.Len(t, data, 100)
		assert.Equal(t, 100, reader.Position())

		data = reader.ReadSlice(100)
		assert.Len(t, data, 100)
		assert.Nil(t, reader.Error())
		assert.Equal(t, 200, reader.Position())

		for i := 0; i < 3; i++ {
			data = reader.ReadSlice(1)
			assert.Len(t, data, 0)
			assert.NotNil(t, reader.Error())
			assert.Equal(t, 200, reader.Position())
		}
	}

	reader.Reset(sl)
	_ = reader.ReadSlice(-1)
	assert.NotNil(t, reader.Error())

	reader.Reset(sl)
	_ = reader.ReadBytes(-1)
	assert.NotNil(t, reader.Error())
}

func Test_Reader_SeekStart(t *testing.T) {
	sl := make([]byte, 200)
	reader := stream.NewReader(sl)

	assert.Zero(t, reader.Position())
	_ = reader.ReadByte()
	reader.SeekStart()
	assert.Zero(t, reader.Position())

	reader.ReadSlice(300)
	assert.NotNil(t, reader.Error())
	assert.Equal(t, 200, reader.Position())

	reader.SeekStart()
	assert.Zero(t, reader.Position())
	assert.Nil(t, reader.Error())
}

func Test_Stream_UvariantSize(t *testing.T) {
	assert.Equal(t, 1, stream.UvariantSize(0))
	assert.Equal(t, 1, stream.UvariantSize(1))
	assert.Equal(t, 1, stream.UvariantSize(127))
	assert.Equal(t, 2, stream.UvariantSize(129))
	assert.Equal(t, 2, stream.UvariantSize(16383))
	assert.Equal(t, 3, stream.UvariantSize(16384))
	assert.Equal(t, 3, stream.UvariantSize(2097151))
	assert.Equal(t, 4, stream.UvariantSize(2097152))
}

func Test_Stream_VariantSize(t *testing.T) {
	assert.Equal(t, 1, stream.VariantSize(0))
	assert.Equal(t, 1, stream.VariantSize(-63))
	assert.Equal(t, 1, stream.VariantSize(63))
	assert.Equal(t, 1, stream.VariantSize(-64))
	assert.Equal(t, 2, stream.VariantSize(-65))
	assert.Equal(t, 2, stream.VariantSize(64))
	assert.Equal(t, 2, stream.VariantSize(-127))
	assert.Equal(t, 2, stream.VariantSize(127))
	assert.Equal(t, 2, stream.VariantSize(8191))
	assert.Equal(t, 2, stream.VariantSize(-8191))
	assert.Equal(t, 3, stream.VariantSize(8192))
	assert.Equal(t, 2, stream.VariantSize(-8192))
	assert.Equal(t, 3, stream.VariantSize(-8193))
}

func Benchmark_Reader_ReadBytes(b *testing.B) {
	sl := make([]byte, 1024*1024)
	reader := stream.NewReader(sl)
	for i := 0; i < b.N; i++ {
		reader.Reset(sl)
		_ = reader.ReadBytes(4096)
	}
}

func Benchmark_Reader_ReadSlice(b *testing.B) {
	sl := make([]byte, 1024*1024)
	reader := stream.NewReader(sl)
	for i := 0; i < b.N; i++ {
		reader.Reset(sl)
		_ = reader.ReadSlice(4096)
	}
}

func Benchmark_readUvariant64(b *testing.B) {
	var buf bytes.Buffer
	w := stream.NewBufferWriter(&buf)
	for i := 0; i < 1000; i++ {
		w.PutUvarint64(uint64(i))
	}

	r := stream.NewReader(buf.Bytes())
	for i := 0; i < b.N; i++ {
		r.SeekStart()
		for x := 0; x < 1000; x++ {
			_ = r.ReadUvarint64()
		}
	}
}

func Test_readUvarint(t *testing.T) {
	r := stream.NewReader(nil)
	_ = r.ReadUvarint64()
	assert.NotNil(t, r.Error())

	data := []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x1, 0, 0}
	r = stream.NewReader(data)
	_ = r.ReadUvarint64()
	assert.NotNil(t, r.Error())
}
