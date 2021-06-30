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

package tagkeymeta

import (
	"errors"
	"testing"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/kv/table"
	"github.com/lindb/lindb/sql/stmt"

	"github.com/golang/mock/gomock"
	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func buildTrieBlock() (zoneBlock []byte, ipBlock []byte, hostBlock []byte) {
	// tag id mapping relation
	/////////////////////////
	ipMapping := map[uint32]string{
		1: "192.168.1.1",
		2: "192.168.1.2",
		3: "192.168.1.3",
		4: "192.168.2.4",
		5: "192.168.2.5",
		6: "192.168.2.6",
		7: "192.168.3.7",
		8: "192.168.3.8",
		9: "192.168.3.9"}
	zoneMapping := map[uint32]string{
		1: "nj",
		2: "sh",
		3: "bj"}
	hostMapping := map[uint32]string{
		1:    "eleme-dev-nj-1",
		2:    "eleme-dev-nj-2",
		3:    "eleme-dev-nj-3",
		4:    "eleme-dev-sh-4",
		5:    "eleme-dev-sh-5",
		6000: "eleme-dev-sh-6000",
		7:    "eleme-dev-bj-7",
		8:    "eleme-dev-bj-8",
		9:    "eleme-dev-bj-9"}
	flush := func(flusher Flusher, mapping map[uint32]string) {
		for id, value := range mapping {
			flusher.FlushTagValue([]byte(value), id)
		}
	}
	/////////////////////////
	// flush zone tag, tagKeyID: 20
	/////////////////////////
	nopKVFlusher1 := kv.NewNopFlusher()
	flusher1 := NewFlusher(nopKVFlusher1)
	flush(flusher1, zoneMapping)
	// pick the zoneBlock buffer
	_ = flusher1.FlushTagKeyID(20, 20)
	zoneBlock = append(zoneBlock, nopKVFlusher1.Bytes()...)

	/////////////////////////
	// flush ip tag, tagKeyID: 21
	/////////////////////////
	nopKVFlusher2 := kv.NewNopFlusher()
	flusher2 := NewFlusher(nopKVFlusher2)
	flush(flusher2, ipMapping)
	// pick the ipBlock buffer
	_ = flusher2.FlushTagKeyID(21, 21)
	ipBlock = append(ipBlock, nopKVFlusher2.Bytes()...)

	/////////////////////////
	// flush host tag, tagKeyID: 22
	/////////////////////////
	nopKVFlusher3 := kv.NewNopFlusher()
	flusher3 := NewFlusher(nopKVFlusher3)
	flush(flusher3, hostMapping)
	// pick the hostBlock buffer
	_ = flusher3.FlushTagKeyID(22, 22)
	hostBlock = append(hostBlock, nopKVFlusher3.Bytes()...)
	return zoneBlock, ipBlock, hostBlock
}

func mockTagReader(ctrl *gomock.Controller) Reader {
	zoneBlock, ipBlock, hostBlock := buildTrieBlock()
	// mock readers
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(10)).Return(nil, true).AnyTimes()
	mockReader.EXPECT().Get(uint32(19)).Return(nil, false).AnyTimes()
	mockReader.EXPECT().Get(uint32(20)).Return(zoneBlock, true).AnyTimes()
	mockReader.EXPECT().Get(uint32(21)).Return(ipBlock, true).AnyTimes()
	mockReader.EXPECT().Get(uint32(22)).Return(hostBlock, true).AnyTimes()
	// build tag reader
	return NewReader([]table.Reader{mockReader})
}

func mockBadTagReader(ctrl *gomock.Controller) Reader {
	zoneBlock, _, _ := buildTrieBlock()
	badZoneBlock := append(zoneBlock, byte(1), byte(1))
	mockReader := table.NewMockReader(ctrl)
	mockReader.EXPECT().Get(uint32(23)).Return(badZoneBlock, true).AnyTimes()
	return NewReader([]table.Reader{mockReader})
}

func TestReader_GetTagValueSeq(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mockTagReader(ctrl)
	// case 1: tag key id not exist
	id, err := reader.GetTagValueSeq(19)
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Equal(t, uint32(0), id)
	// case 2: get value
	id, err = reader.GetTagValueSeq(22)
	assert.NoError(t, err)
	assert.Equal(t, uint32(22), id)

	// case3: newTagKeyMeta error
	id, err = mockBadTagReader(ctrl).GetTagValueSeq(23)
	assert.Error(t, err)
	assert.Zero(t, id)
}

func TestReader_GetTagValueID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mockTagReader(ctrl)
	// case 1: tag key id not exist
	id, err := reader.GetTagValueID(19, "eleme-dev-sh-5")
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Equal(t, uint32(0), id)

	// case 2: get value
	id, err = reader.GetTagValueID(22, "eleme-dev-sh-5")
	assert.NoError(t, err)
	assert.Equal(t, uint32(5), id)

	// case 3: tag value not found
	id, err = reader.GetTagValueID(22, "eleme-dev-sh-5999")
	assert.True(t, errors.Is(err, constants.ErrNotFound))
	assert.Equal(t, uint32(0), id)

	// case 4: new tag key meta err
	id, err = mockBadTagReader(ctrl).GetTagValueID(23, "eleme-dev-sh-5")
	assert.Error(t, err)
	assert.Equal(t, uint32(0), id)
}

func TestReader_GetTagValueIDsForTagKeyID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mockTagReader(ctrl)
	// read unexisted tagKeyID key
	idSet, err := reader.GetTagValueIDsForTagKeyID(19)
	assert.Error(t, err)
	assert.Nil(t, idSet)
	idSet, err = reader.GetTagValueIDsForTagKeyID(10)
	assert.Error(t, err)
	assert.Nil(t, idSet)

	// read zone block
	idSet, err = reader.GetTagValueIDsForTagKeyID(20)
	assert.NoError(t, err)
	assert.NotNil(t, idSet)
	assert.EqualValues(t, roaring.BitmapOf(1, 2, 3).ToArray(), idSet.ToArray())
}

func TestReader_FindValueIDsByExprForTagKeyID_bad_case(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mockTagReader(ctrl)

	// tagKeyID not exist
	idSet, err := reader.FindValueIDsByExprForTagKeyID(19, nil)
	assert.Error(t, err)
	assert.Nil(t, idSet)

	// find zone with bad expression
	idSet, err = reader.FindValueIDsByExprForTagKeyID(20, nil)
	assert.Error(t, err)
	assert.Nil(t, idSet)

	// value not exist
	idSet, err = reader.FindValueIDsByExprForTagKeyID(20, &stmt.EqualsExpr{Key: "zone", Value: "not-exist"})
	assert.Error(t, err)
	assert.Nil(t, idSet)
}

func TestReader_FindSeriesIDsByExprForTagID_EqualExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reader := mockTagReader(ctrl)

	idSet, err := reader.FindValueIDsByExprForTagKeyID(22, &stmt.EqualsExpr{Key: "host", Value: "eleme-dev-sh-4"})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(4), idSet)
	// find not existed host
	_, err = reader.FindValueIDsByExprForTagKeyID(22, &stmt.EqualsExpr{Key: "host", Value: "eleme-dev-sh-41"})
	assert.Error(t, err)
}

func TestReader_FindValueIDsByExprForTagKeyID_InExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reader := mockTagReader(ctrl)

	// find existed host
	idSet, err := reader.FindValueIDsByExprForTagKeyID(22, &stmt.InExpr{
		Key: "host", Values: []string{"eleme-dev-sh-4", "eleme-dev-sh-5", "eleme-dev-sh-55"}},
	)
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(4, 5), idSet)

	// find not existed host
	_, err = reader.FindValueIDsByExprForTagKeyID(22, &stmt.InExpr{
		Key: "host", Values: []string{"eleme-dev-sh-55"}},
	)
	assert.Error(t, err)
}

func TestReader_FindSeriesIDsByExprForTagID_LikeExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reader := mockTagReader(ctrl)

	// find existed host
	idSet, err := reader.FindValueIDsByExprForTagKeyID(22, &stmt.LikeExpr{Key: "host", Value: "eleme-dev-sh-*"})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(4, 5, 6000), idSet)
	// find not existed host
	_, err = reader.FindValueIDsByExprForTagKeyID(22, &stmt.InExpr{Key: "host", Values: []string{"eleme-dev-sh---"}})
	assert.Error(t, err)
}

func TestReader_FindSeriesIDsByExprForTagID_RegexExpr(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	reader := mockTagReader(ctrl)

	idSet, err := reader.FindValueIDsByExprForTagKeyID(22, &stmt.RegexExpr{Key: "host", Regexp: "eleme-dev-sh-"})
	assert.NoError(t, err)
	assert.Equal(t, roaring.BitmapOf(4, 5, 6000), idSet)

	// find not existed host
	_, err = reader.FindValueIDsByExprForTagKeyID(22, &stmt.RegexExpr{Key: "host", Regexp: "eleme-prod-sh-"})
	assert.Error(t, err)
}

func TestReader_SuggestTagValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mockTagReader(ctrl)

	// tagKeyID not exist
	assert.Nil(t, reader.SuggestTagValues(19, "", 10000000))
	// search ip
	assert.Len(t, reader.SuggestTagValues(21, "192", 1000), 9)
	assert.Len(t, reader.SuggestTagValues(21, "192", 3), 3)

	// mock corruption
	badReader := mockBadTagReader(ctrl)
	assert.Nil(t, badReader.SuggestTagValues(23, "", 10000000))
}

func Test_Reader_WalkTagValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := mockTagReader(ctrl)

	// tagKeyID not exist
	assert.NotPanics(t, func() {
		_ = reader.WalkTagValues(
			19,
			"",
			func(tagValue []byte, tagValueID uint32) bool {
				panic("tagKeyID doesn't exist!")
			})
	})
	assert.NotPanics(t, func() {
		_ = reader.WalkTagValues(
			10,
			"",
			func(tagValue []byte, tagValueID uint32) bool {
				panic("tagKeyID doesn't exist!")
			})
	})

	// search ip
	var ipCount1 int
	assert.Nil(t, reader.WalkTagValues(
		21,
		"192",
		func(tagValue []byte, tagValueID uint32) bool {
			ipCount1++
			return true
		}))
	assert.Equal(t, 9, ipCount1)

	// break case
	var ipCount2 int
	assert.Nil(t, reader.WalkTagValues(
		21,
		"192",
		func(tagValue []byte, tagValueID uint32) bool {
			ipCount2++
			return ipCount2 != 3
		}))
	assert.Equal(t, 3, ipCount2)
}

func TestReader_CollectTagValues(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockReader := mockTagReader(ctrl)

	// case 1: ok
	err := mockReader.CollectTagValues(21, roaring.BitmapOf(1, 2, 3), map[uint32]string{})
	assert.Nil(t, err)
	// case 2: tag value ids is empty
	err = mockReader.CollectTagValues(21, roaring.New(), map[uint32]string{})
	assert.NoError(t, err)
	// case 3: tag key not found
	err = mockReader.CollectTagValues(19, roaring.BitmapOf(1, 2, 3), map[uint32]string{})
	assert.NoError(t, err)
}
