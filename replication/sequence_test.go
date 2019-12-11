package replication

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
)

func TestSequence(t *testing.T) {
	seq, err := NewSequence("/test")
	assert.Error(t, err)
	assert.Nil(t, seq)

	tmp := path.Join(os.TempDir(), "sequence_test")
	if err := fileutil.RemoveDir(tmp); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := fileutil.RemoveDir(tmp); err != nil {
			t.Error(err)
		}
	}()

	seq, err = NewSequence(tmp)
	assert.NoError(t, err)

	assert.Equal(t, seq.GetHeadSeq(), int64(0))
	assert.Equal(t, seq.GetAckSeq(), int64(0))
	assert.False(t, seq.Synced())

	seq.SetHeadSeq(int64(10))
	seq.SetAckSeq(int64(5))

	assert.Equal(t, seq.GetHeadSeq(), int64(10))
	assert.Equal(t, seq.GetAckSeq(), int64(5))

	if err := seq.Sync(); err != nil {
		t.Fatal(err)
	}

	assert.True(t, seq.Synced())
	seq.ResetSynced()
	assert.False(t, seq.Synced())

	// new sequence
	newSeq, err := NewSequence(tmp)
	assert.NoError(t, err)

	assert.Equal(t, newSeq.GetAckSeq(), int64(5))
	assert.Equal(t, newSeq.GetHeadSeq(), int64(5))
	assert.False(t, newSeq.Synced())
}
