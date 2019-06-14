package meta

import "sync/atomic"

type Version struct {
	fv  *FamilyVersion
	ref *int32
}

func newVersion(fv *FamilyVersion) *Version {
	var ref int32
	return &Version{
		fv:  fv,
		ref: &ref,
	}
}

func (v *Version) Release() {
	val := atomic.AddInt32(v.ref, -1)
	if val == 0 {
		v.fv.RemoveVersion(v)
	}
}

func (v *Version) Retain() {
	atomic.AddInt32(v.ref, 1)
}
