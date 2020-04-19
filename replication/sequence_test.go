package replication

import (
	"fmt"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/queue/page"
)

var testPath = "test"

func TestSequence_new_err(t *testing.T) {
	ctrl := gomock.NewController(t)
	tmp := path.Join(testPath, "sequence_test")
	defer func() {
		newPageFactoryFunc = page.NewFactory
		_ = fileutil.RemoveDir(testPath)
		ctrl.Finish()
	}()
	// case 1: new page factory err
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return nil, fmt.Errorf("err")
	}
	seq, err := NewSequence(tmp)
	assert.Error(t, err)
	assert.Nil(t, seq)
	// case 2: AcquirePage err
	fct := page.NewMockFactory(ctrl)
	newPageFactoryFunc = func(path string, pageSize int) (page.Factory, error) {
		return fct, nil
	}
	fct.EXPECT().Close().Return(fmt.Errorf("err"))
	fct.EXPECT().AcquirePage(gomock.Any()).Return(nil, fmt.Errorf("err"))
	seq, err = NewSequence(tmp)
	assert.Error(t, err)
	assert.Nil(t, seq)
}

func TestSequence(t *testing.T) {
	tmp := path.Join(testPath, "sequence_test")

	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()

	seq, err := NewSequence(tmp)
	assert.NoError(t, err)
	assert.NotNil(t, seq)

	assert.Equal(t, seq.GetHeadSeq(), int64(0))
	assert.Equal(t, seq.GetAckSeq(), int64(0))
	assert.False(t, seq.Synced())

	seq.SetHeadSeq(int64(10))
	seq.SetAckSeq(int64(5))

	assert.Equal(t, seq.GetHeadSeq(), int64(10))
	assert.Equal(t, seq.GetAckSeq(), int64(5))

	err = seq.Sync()
	assert.NoError(t, err)

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
