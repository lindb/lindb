package lockers

import (
	"io/ioutil"
	"os"
	"testing"

	. "gopkg.in/check.v1"
)

type TestSuite struct {
	path string
	lock *FileLock
}

var _ = Suite(&TestSuite{})

func Test(t *testing.T) { TestingT(t) }

func (t *TestSuite) SetUpTest(c *C) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "go-flock-")
	c.Assert(err, IsNil)
	c.Assert(tmpFile, Not(IsNil))

	t.path = tmpFile.Name()

	defer os.Remove(t.path)
	tmpFile.Close()

	t.lock = NewFileLock(t.path)
}

func (t *TestSuite) TestFileLock(c *C) {
	var lock = t.lock
	defer os.Remove(t.path)
	var err = lock.Lock()
	c.Assert(err, IsNil)

	err = lock.Lock()
	c.Assert(err, NotNil)

	c.Assert(lock.IsLocked(), Equals, true)

	err = lock.Unlock()
	c.Assert(err, IsNil)

}

func (t *TestSuite) TestFileLock_RLock(c *C) {
	var lock = t.lock
	defer os.Remove(t.path)
	var err = lock.RLock()
	c.Assert(err, IsNil)
	c.Assert(lock.IsRLocked(), Equals, true)
}

func (t *TestSuite) TestFileLock_RLock2(c *C) {
	var lock = t.lock
	defer os.Remove(t.path)
	_ = lock.RLock()
	var err = lock.RLock()
	c.Assert(err, NotNil)
}

func (t *TestSuite) TestFileLock_Unlock(c *C) {
	var lock = t.lock
	lock = NewFileLock(t.path)
	var err = lock.Unlock()
	c.Assert(err, NotNil)
}
