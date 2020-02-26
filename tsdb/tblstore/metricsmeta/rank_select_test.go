package metricsmeta

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RankSelect(t *testing.T) {
	rs := NewRankSelect()
	assert.NotNil(t, rs)

	assert.Zero(t, rs.MaxNodeNumber())
	assert.Zero(t, rs.Num())
	assert.Zero(t, rs.OneNum())
	assert.Zero(t, rs.ZeroNum())
	rs.PushPseudoRoot()
	rs.PushPseudoRoot()
	rs.PushPseudoRoot()
	assert.Equal(t, uint64(1), rs.MaxNodeNumber())
	assert.Equal(t, uint64(2), rs.Num())
	assert.Equal(t, uint64(1), rs.OneNum())
	assert.Equal(t, uint64(1), rs.ZeroNum())

}

/*
                +--------+
                |        |
                |  0-10  |
                +---+----+
                    |
                    |
                +--------+
                |        |
                |  1-110 |
                +---+----+
                   / \
                  /   \
                 /     \
       +--------+       +--------+
       |        |       |        |
       | 2-110  |       |  3-0   |
       +---+--+-+       +---+----+
           |  |
          /   |
         /    +----------+
        /                 \
   +--------+          +--------+
   |        |          |        |
   |   4-0  |          | 5-110  |
   +---+----+          +-+--+---+
                         |   \
                         |    ---------+
                     +--------+        +--------+
                     |        |        |        |
                     |  6-0   |        |  7-0   |
                     +--------+        +--------+

*/

func Test_RankSelect_Helpers(t *testing.T) {
	rs := NewRankSelect()
	rs.PushPseudoRoot()
	for _, char := range "1101100011000" {
		if string(char) == "1" {
			rs.PushBack(true)
		} else {
			rs.PushBack(false)
		}
	}
	assert.Equal(t, "101101100011000", rs.String())

	type result struct {
		hasParent   bool
		parent      uint64
		hasChild    bool
		firstChild  uint64
		lastChild   uint64
		childrenNum uint64
	}
	expects := []result{
		{hasParent: false, hasChild: false},
		{hasParent: false, hasChild: true, firstChild: 2, lastChild: 3, childrenNum: 2},
		{hasParent: true, parent: 1, hasChild: true, firstChild: 4, lastChild: 5, childrenNum: 2},
		{hasParent: true, parent: 1, hasChild: false, childrenNum: 0},
		{hasParent: true, parent: 2, hasChild: false, childrenNum: 0},
		{hasParent: true, parent: 2, hasChild: true, firstChild: 6, lastChild: 7, childrenNum: 2},
		{hasParent: true, parent: 5, hasChild: false, childrenNum: 0},
		{hasParent: true, parent: 5, hasChild: false, childrenNum: 0},
	}

	for nodeNumber, expect := range expects {
		nodeNumber := uint64(nodeNumber)
		firstChildNodeNumber, ok1 := rs.FirstChild(nodeNumber)
		lastChildNodeNumber, ok2 := rs.LastChild(nodeNumber)
		childrenNum := rs.ChildrenNum(nodeNumber)
		if expect.hasChild {
			assert.Equalf(t, firstChildNodeNumber, expect.firstChild,
				"node: %d firstChildNodeNumber not match", nodeNumber)
			assert.Equalf(t, lastChildNodeNumber, expect.lastChild,
				"node: %d lastChildNodeNumber not match", nodeNumber)
			assert.True(t, ok1)
			assert.True(t, ok2)
			assert.Equalf(t, childrenNum, expect.childrenNum,
				"node: %d childrenNum not match", nodeNumber)
		} else {
			assert.False(t, ok1)
			assert.False(t, ok2)
			assert.Zerof(t, childrenNum,
				"node: %d childrenNum not zero", nodeNumber)
		}
		if expect.hasParent {
			assert.Equalf(t, expect.parent, rs.Parent(nodeNumber),
				"node: %d parent number not match", nodeNumber)
		}
	}
}

func Test_RankSelect_Marshal(t *testing.T) {
	rs1 := NewRankSelect()
	out, err := rs1.MarshalBinary()
	assert.Nil(t, err)
	rs2 := NewRankSelect()
	assert.Nil(t, rs2.UnmarshalBinary(out))

	rs1.PushBack(true)
	rs1.PushBack(true)
	rs1.PushBack(false)
	rs1.PushBack(true)
	rs1.PushBack(false)

	out, _ = rs1.MarshalBinary()
	assert.Len(t, out, 17)
	rs2.UnmarshalBinary(out)
	out2, _ := rs2.MarshalBinary()
	assert.Equal(t, out2, out)
}
