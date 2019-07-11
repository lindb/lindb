package page

import (
	"testing"

	"gopkg.in/check.v1"
)

type pageTestSuite struct {
}

var _ = check.Suite(&pageTestSuite{})

func Test(t *testing.T) {
	check.TestingT(t)
}

func (ts *pageTestSuite) TestMethods(c *check.C) {
	fileName := "fileName"
	bytes := []byte("12345")
	mp := NewMappedPage(fileName, bytes,
		func(mappedBytes []byte) error {
			return nil
		}, func(mappedBytes []byte) error {
			return nil
		})

	c.Check(mp.FilePath(), check.Equals, fileName)

	c.Check(mp.Size(), check.Equals, len(bytes))

	c.Check(mp.Buffer(0), check.DeepEquals, bytes[0:])
	c.Check(mp.Buffer(3), check.DeepEquals, bytes[3:])

	c.Check(mp.Data(0, 2), check.DeepEquals, bytes[0:2])
	c.Check(mp.Data(3, 2), check.DeepEquals, bytes[3:5])

	c.Check(mp.Closed(), check.Equals, false)

	c.Check(mp.Close(), check.IsNil)

	c.Check(mp.Closed(), check.Equals, true)

	c.Check(mp.Close(), check.IsNil)

}
