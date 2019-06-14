package meta

type FamilyVersion struct {
	current *Version
}

func NewFamilyVersion() *FamilyVersion {
	fv := &FamilyVersion{}

	fv.current = newVersion(fv)
	return fv
}

func (fv *FamilyVersion) GetCurrent() *Version {
	return fv.current
}

func (fv *FamilyVersion) RemoveVersion(v *Version) {

}
