package segment

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/queue/page"
	"github.com/lindb/lindb/pkg/stream"
)

var closeFunc = func([]byte) error {
	return nil
}

var syncFunc = func([]byte) error {
	return nil
}

func buildDataPage(t *testing.T, data ...[]byte) page.MappedPage {
	dataWriter := stream.NewBufferWriter(nil)

	for _, msg := range data {
		dataWriter.PutBytes(msg)
	}
	dataBytes, err := dataWriter.Bytes()
	if err != nil {
		t.Error(err)
	}

	dataPage := page.NewMappedPage("0.dat", dataBytes, closeFunc, syncFunc)

	return dataPage
}

func buildIndexPage(t *testing.T, data ...[]byte) page.MappedPage {
	indexWriter := stream.NewBufferWriter(nil)

	offset := int32(0)
	for _, msg := range data {
		indexWriter.PutInt32(offset)
		indexWriter.PutInt32(int32(len(msg)))
		offset += int32(len(msg))
	}

	indexBytes, err := indexWriter.Bytes()
	if err != nil {
		t.Error(err)
	}

	indexPage := page.NewMappedPage("0.idx", indexBytes, closeFunc, syncFunc)

	return indexPage
}

func TestRead(t *testing.T) {
	data := [][]byte{
		[]byte("welcome"),
		[]byte("to"),
		[]byte("LinDB"),
	}

	seg, err := NewSegment(buildIndexPage(t, data...), buildDataPage(t, data...), int64(0), int64(len(data)))

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, seg.Begin(), int64(0))
	assert.Equal(t, seg.End(), int64(len(data)))

	assert.Equal(t, seg.Contains(0), true)
	assert.Equal(t, seg.Contains(1), true)
	assert.Equal(t, seg.Contains(2), true)
	assert.Equal(t, seg.Contains(3), false)

	for i, msg := range data {
		rmsg, err := seg.Read(int64(i))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, rmsg, msg)
	}

	_, err = seg.Read(3)
	if err == nil {
		t.Error(errors.New("should be nil"))
	}

}

func TestWrite(t *testing.T) {
	dataBytes := [3]byte{}
	indexBytes := [24]byte{}

	indexPage := page.NewMappedPage("0.idx", indexBytes[:], closeFunc, syncFunc)
	dataPage := page.NewMappedPage("0.dat", dataBytes[:], closeFunc, syncFunc)

	seg, err := NewSegment(indexPage, dataPage, int64(0), int64(0))

	if err != nil {
		t.Fatal(err)
	}

	var (
		msg, rmsg []byte
		seq       int64
	)

	// first msg
	msg = []byte("1")

	seq, err = seg.Append(msg)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seq, int64(0))

	rmsg, err = seg.Read(seq)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rmsg, msg)

	// second msg
	msg = []byte("23")

	seq, err = seg.Append(msg)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seq, int64(1))

	rmsg, err = seg.Read(seq)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, rmsg, msg)

	if _, err = seg.Append([]byte("4")); err != ErrExceedPageSize {
		t.Fatal(err)
	}

}

func TestSegment_Append(t *testing.T) {
	msg0 := []byte("123")
	dataWriter := stream.NewSliceWriter(make([]byte, 10))
	dataWriter.PutBytes(msg0)

	indexWriter := stream.NewSliceWriter(make([]byte, 16))
	indexWriter.PutInt32(0)
	indexWriter.PutInt32(3)

	dataBytes, err := dataWriter.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	indexBytes, err := indexWriter.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	seg, err := NewSegment(page.NewMappedPage("0.idx", indexBytes, closeFunc, syncFunc),
		page.NewMappedPage("0.dat", dataBytes, closeFunc, syncFunc), 0, 1)

	if err != nil {
		t.Fatal(err)
	}

	msg1 := []byte("4567890")
	seq, err := seg.Append(msg1)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seq, int64(1))

	msgr0, err := seg.Read(0)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, msgr0, msg0)

	msgr1, err := seg.Read(1)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, msgr1, msg1)

}
