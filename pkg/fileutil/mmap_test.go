// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package fileutil

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"gopkg.in/check.v1"
)

type testSuite struct {
}

var _ = check.Suite(&testSuite{})

var filename = path.Join(os.TempDir(), "testdata")

func Test(t *testing.T) {
	check.TestingT(t)
}

func (ts *testSuite) TearDownTest(c *check.C) {
	if err := os.Remove(filename); err != nil {
		c.Error("tear down test remove tmp file error", err)
	}
}

func (ts *testSuite) TestRead(c *check.C) {
	file, err := os.Create(filename)
	if err != nil {
		c.Fatal(err)
	}

	content := "abc123"

	_, err = file.WriteString(content)
	if err != nil {
		c.Fatal(err)
	}

	bys, err := Map(filename)
	if err != nil {
		c.Fatal(c)
	}
	c.Assert(len(bys), check.Equals, len(content))
	c.Assert(bys, check.DeepEquals, []byte(content))

}

func (ts *testSuite) TestRWMap(c *check.C) {
	var content = []byte("12345")
	var size = 1024

	mapBytes, err := RWMap(filename, size)
	if err != nil {
		c.Error("RWMap", err)
	}
	if Unmap(nil) != nil {
		c.Error("unmap nil returns not nil")
	}

	buffer := bytes.NewBuffer(mapBytes[:0])

	if _, err := buffer.Write(content); err != nil {
		c.Error("buffer write", err)
	}

	if err := Sync(mapBytes); err != nil {
		c.Error(err)
	}

	if Unmap(mapBytes) != nil {
		c.Errorf("unmap mapBytes with error: %v", err)
	}

	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		c.Error("read file error", err)
	}

	c.Assert(len(fileContent), check.Equals, size)

	c.Assert(content, check.DeepEquals, fileContent[:len(content)])
}
