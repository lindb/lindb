package tblstore

import (
	"fmt"
	"strings"

	"github.com/hillbig/rsdic"
)

//go:generate mockgen -source ./rank_select.go -destination=./rank_select_mock.go -package tblstore

// Ref:
// [1] Engineering the LOUDS succinct tree representation:
//     http://citeseerx.ist.psu.edu/viewdoc/download?doi=10.1.1.106.4250&rep=rep1&type=pdf

// RSINTF abstracts the rank select implementation(RSDIC).
type RSINTF interface {
	///////////////////////////////////////////////
	// Raw functions of rsdic
	///////////////////////////////////////////////
	// Num returns the number of bits
	Num() uint64
	// OneNum returns the number of ones in bits
	OneNum() uint64
	// ZeroNum returns the number of zeros in bits
	ZeroNum() uint64
	// PushBack appends the bit to the end of B
	PushBack(bit bool)
	Bit(pos uint64) bool
	// MarshalBinary encodes the RSDic into a binary form and returns the result.
	MarshalBinary() (out []byte, err error)
	// UnmarshalBinary decodes the RSDic from a binary from generated MarshalBinary
	UnmarshalBinary(in []byte) (err error)
	///////////////////////////////////////////////
	// Rank and Select primitives
	///////////////////////////////////////////////
	// Rank returns the number of bit's in B[0...pos)
	Rank1(pos uint64) uint64
	Rank0(pos uint64) uint64
	// Select returns the position of (rank+1)-th occurrence of bit in B
	// Select returns num if rank+1 is larger than the possible range.
	// (i.e. Select(oneNum, true) = num, Select(zeroNum, false) = num)
	Select1(rank uint64) uint64
	Select0(rank uint64) uint64
	///////////////////////////////////////////////
	// Helper functions for building or traversing the tree
	///////////////////////////////////////////////
	// insert a pseudo root node for easier calculating
	PushPseudoRoot()
	// Stringer builds the LOUDS bit-string(LBS)
	fmt.Stringer
	// NodeNumber returns the node number from position
	NodeNumber(pos uint64) uint64
	// Position returns the position from a given nodeNumber
	Position(nodeNumber uint64) uint64
	// FirstChild returns the first-child from the given nodeNumber
	FirstChild(nodeNumber uint64) (childNodeNumber uint64, ok bool)
	// LastChild returns the last-child from the given nodeNumber
	LastChild(nodeNumber uint64) (childNodeNumber uint64, ok bool)
	// ChildrenNum returns the children count from the given nodeNumber
	ChildrenNum(nodeNumber uint64) uint64
	// Parent returns the parent from the given nodeNumber
	Parent(nodeNumber uint64) (parentNumber uint64)
}

func NewRankSelect() RSINTF {
	return &rankSelect{rs: rsdic.New()}
}

// rankSelect wraps the rsdic library
type rankSelect struct {
	rs *rsdic.RSDic
}

func (rs *rankSelect) Num() uint64                            { return rs.rs.Num() }
func (rs *rankSelect) OneNum() uint64                         { return rs.rs.OneNum() }
func (rs *rankSelect) ZeroNum() uint64                        { return rs.rs.ZeroNum() }
func (rs *rankSelect) Bit(pos uint64) bool                    { return rs.rs.Bit(pos) }
func (rs *rankSelect) MarshalBinary() (out []byte, err error) { return rs.rs.MarshalBinary() }
func (rs *rankSelect) UnmarshalBinary(in []byte) (err error)  { return rs.rs.UnmarshalBinary(in) }
func (rs *rankSelect) PushBack(bit bool)                      { rs.rs.PushBack(bit) }
func (rs *rankSelect) Rank1(pos uint64) uint64                { return rs.rs.Rank(pos+1, true) }
func (rs *rankSelect) Rank0(pos uint64) uint64                { return rs.rs.Rank(pos+1, false) }
func (rs *rankSelect) Select1(rank uint64) uint64             { return rs.rs.Select1(rank - 1) }
func (rs *rankSelect) Select0(rank uint64) uint64             { return rs.rs.Select0(rank - 1) }

func (rs *rankSelect) PushPseudoRoot() {
	if rs.rs.Num() > 0 {
		return
	}
	rs.rs.PushBack(true)
	rs.rs.PushBack(false)
}

func (rs *rankSelect) String() string {
	var s strings.Builder
	for i := uint64(0); i < rs.Num(); i++ {
		if rs.Bit(i) {
			s.WriteString("1")
		} else {
			s.WriteString("0")
		}
	}
	return s.String()
}

func (rs *rankSelect) NodeNumber(pos uint64) uint64      { return rs.Rank1(pos) }
func (rs *rankSelect) Position(nodeNumber uint64) uint64 { return rs.Select1(nodeNumber) }

func (rs *rankSelect) FirstChild(nodeNumber uint64) (childNodeNumber uint64, ok bool) {
	pos := rs.Position(nodeNumber)
	childPos := rs.Select0(rs.Rank1(pos)) + 1
	return rs.NodeNumber(childPos), rs.Bit(childPos)
}

func (rs *rankSelect) LastChild(nodeNumber uint64) (childNodeNumber uint64, ok bool) {
	pos := rs.Position(nodeNumber)
	childPos := rs.Select0(rs.Rank1(pos)+1) - 1
	return rs.NodeNumber(childPos), rs.Bit(childPos)
}

func (rs *rankSelect) ChildrenNum(nodeNumber uint64) uint64 {
	firstChildNumber, ok := rs.FirstChild(nodeNumber)
	if !ok {
		return 0
	}
	lastChildNumber, _ := rs.LastChild(nodeNumber)
	return lastChildNumber - firstChildNumber + 1
}

func (rs *rankSelect) Parent(nodeNumber uint64) (parentNumber uint64) {
	pos := rs.Position(nodeNumber)
	parentPos := rs.Select1(rs.Rank0(pos))
	return rs.NodeNumber(parentPos)
}
