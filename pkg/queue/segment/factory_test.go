package segment

import (
	"os"
	"path"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/stream"
)

func TestSeqRange(t *testing.T) {
	var (
		sr     SeqRange
		index  int
		seqNum int64
		ok     bool
	)

	sr = SeqRange{}
	_, _, ok = sr.Find(0)
	assert.False(t, ok)

	sr = SeqRange{10}
	index, seqNum, ok = sr.Find(10)
	assert.Equal(t, index, 0)
	assert.Equal(t, seqNum, sr[index])
	assert.True(t, ok)

	index, seqNum, ok = sr.Find(20)
	assert.Equal(t, index, 0)
	assert.Equal(t, seqNum, sr[index])
	assert.True(t, ok)

	_, _, ok = sr.Find(0)
	assert.False(t, ok)

	sr = SeqRange{10, 5, 20}
	sort.Sort(sr)

	index, seqNum, ok = sr.Find(10)
	assert.Equal(t, index, 1)
	assert.Equal(t, seqNum, sr[index])
	assert.True(t, ok)

	index, seqNum, ok = sr.Find(11)
	assert.Equal(t, index, 1)
	assert.Equal(t, seqNum, sr[index])
	assert.True(t, ok)

	index, seqNum, ok = sr.Find(21)
	assert.Equal(t, index, 2)
	assert.Equal(t, seqNum, sr[index])
	assert.True(t, ok)

	_, _, ok = sr.Find(4)
	assert.False(t, ok)

}

func TestEmptyFactory(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), "segment_empty_factory")

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Error(err)
		}

	}()

	fct, err := NewFactory(tmpDir, 1024, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	seg, err := fct.NewSegment(0)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seg.Begin(), int64(0))
	assert.Equal(t, seg.End(), int64(0))

	msg := []byte("123")

	seq, err := seg.Append(msg)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seq, int64(0))
	assert.Equal(t, seg.Begin(), int64(0))
	assert.Equal(t, seg.End(), int64(1))

	msgr, err := seg.Read(seq)
	assert.Nil(t, err)
	assert.Equal(t, msgr, msg)
}

func TestFactory(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), "segment_factory")

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Error(err)
		}

	}()

	writeFile(t, tmpDir, 0, []byte("123"))
	writeFile(t, tmpDir, 1, []byte("456"), []byte("789"))

	fat, err := NewFactory(tmpDir, 10, 3, 0)
	if err != nil {
		t.Fatal(err)
	}

	seg0, err := fat.GetSegment(0)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seg0.Begin(), int64(0))
	assert.Equal(t, seg0.End(), int64(1))

	bys, err := seg0.Read(0)
	assert.Nil(t, err)
	assert.Equal(t, []byte("123"), bys)

	seg11, err := fat.GetSegment(1)
	if err != nil {
		t.Fatal(err)
	}

	seg12, err := fat.GetSegment(2)
	if err != nil {
		t.Fatal(err)
	}

	seg13, err := fat.GetSegment(3)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, seg11, seg12)
	assert.Equal(t, seg11, seg13)

	assert.Equal(t, seg11.Begin(), int64(1))
	assert.Equal(t, seg11.End(), int64(3))

	bys, err = seg11.Read(1)
	assert.Nil(t, err)
	assert.Equal(t, []byte("456"), bys)

	bys, err = seg11.Read(2)
	assert.Nil(t, err)
	assert.Equal(t, []byte("789"), bys)

}

func writeFile(t *testing.T, tmpDir string, seqNum int, msgs ...[]byte) {
	dataFile, err := os.Create(path.Join(tmpDir, strconv.Itoa(seqNum)+dataFileSuffix))
	if err != nil {
		t.Fatal(err)
	}

	indexFile, err := os.Create(path.Join(tmpDir, strconv.Itoa(seqNum)+indexFileSuffix))
	if err != nil {
		t.Fatal(err)
	}

	offset := int32(0)
	for _, msg := range msgs {
		dataLen := int32(len(msg))
		if _, err := dataFile.Write(msg); err != nil {
			t.Fatal(err)
		}

		writer := stream.BinaryWriter()
		writer.PutInt32(offset)
		writer.PutInt32(dataLen)

		offset += dataLen

		bys, err := writer.Bytes()
		if err != nil {
			t.Fatal(err)
		}

		if _, err := indexFile.Write(bys); err != nil {
			t.Fatal(err)
		}
	}

	if err := dataFile.Close(); err != nil {
		t.Fatal(err)
	}

	if err := indexFile.Close(); err != nil {
		t.Fatal(err)
	}

}

func TestFactory_RemoveSegments(t *testing.T) {
	tmpDir := path.Join(os.TempDir(), "segment_factory")

	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Error(err)
		}

	}()

	//[0, 1)
	writeFile(t, tmpDir, 0, []byte("123"))
	//[1, 3)
	writeFile(t, tmpDir, 1, []byte("456"), []byte("789"))
	//[3, 5)
	writeFile(t, tmpDir, 3, []byte("456"), []byte("789"))

	fct, err := NewFactory(tmpDir, 10, 5, 0)
	if err != nil {
		t.Fatal(err)
	}

	// remove the first one
	if err := fct.RemoveSegments(1); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fct.SegmentsSize(), 2)

	// remove nothing
	if err := fct.RemoveSegments(2); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fct.SegmentsSize(), 2)

	// remove the second one
	if err := fct.RemoveSegments(4); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, fct.SegmentsSize(), 1)

	// list file to check
	files, err := fileutil.ListDir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, len(files), 2)

	// new segment fails
	if _, err := fct.NewSegment(0); err == nil {
		t.Fatal(err)
	}

}
