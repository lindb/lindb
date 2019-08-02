package replication

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/fileutil"
)

func TestSequence(t *testing.T) {
	tmp := path.Join(os.TempDir(), "sequence_test")
	if err := fileutil.RemoveDir(tmp); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := fileutil.RemoveDir(tmp); err != nil {
			t.Error(err)
		}
	}()

	seq, err := NewSequence(tmp)
	if err != nil {
		t.Fatal(err)
	}

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
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, newSeq.GetAckSeq(), int64(5))
	assert.Equal(t, newSeq.GetHeadSeq(), int64(5))
	assert.False(t, newSeq.Synced())

}

func TestNewSequenceManager(t *testing.T) {
	tmp := path.Join(os.TempDir(), "sequence_manager_test")
	if err := fileutil.RemoveDir(tmp); err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := fileutil.RemoveDir(tmp); err != nil {
			t.Error(err)
		}
	}()

	sm, err := NewSequenceManager(tmp)
	if err != nil {
		t.Fatal(err)
	}

	node := models.Node{
		IP:   "1.1.1.1",
		Port: 12345,
	}

	_, ok := sm.GetSequence("db", 0, node)
	if ok {
		t.Fatal("should not exists")
	}

	seq1, err := sm.CreateSequence("db", 0, node)
	if err != nil {
		t.Fatal(err)
	}

	seq11, err := sm.CreateSequence("db", 0, node)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, seq1 == seq11)

	seq2, err := sm.CreateSequence("db", 1, node)
	if err != nil {
		t.Fatal(err)
	}
	assert.False(t, seq1 == seq2)
}
