package meta

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/eleme/lindb/pkg/util"
)

type VersionSet struct {
	manifestFileNumber int64
	nextFileNumber     int64
	storePath          string
	familyVersion      map[string]*FamilyVersion
	//mutex              sync.Mutex
}

func NewVersionSet(storePath string) *VersionSet {
	vs := &VersionSet{
		manifestFileNumber: 1, // default value for initialize store
		nextFileNumber:     2, // default value
		storePath:          storePath,
		familyVersion:      make(map[string]*FamilyVersion),
	}
	return vs
}

func (vs *VersionSet) Recover() error {
	if err := vs.initializeIfNeeded(); err != nil {
		return err
	}
	//todo do recover log
	return nil
}

// Get next file number
func (vs *VersionSet) NextFileNumber() int64 {
	return atomic.AddInt64(&vs.nextFileNumber, 1) - 1
}

// Initialize if version file not exists
func (vs *VersionSet) initializeIfNeeded() error {
	if !util.Exist(filepath.Join(vs.storePath, Current())) {
		manifest := manifestFileName(vs.manifestFileNumber)
		tmp := filepath.Join(vs.storePath, fmt.Sprintf("%s.%s", current, tmpSuffix))
		if err := ioutil.WriteFile(tmp, []byte(manifest), 0666); err != nil {
			return err
		}
		if err := os.Rename(tmp, filepath.Join(vs.storePath, current)); err != nil {
			return err
		}
	}
	return nil
}
