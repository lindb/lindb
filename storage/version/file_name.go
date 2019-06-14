package meta

import (
	"fmt"
	"io/ioutil"
	"os"
)

const sstSuffix = "sst"
const tmpSuffix = "tmp"
const current = "CURRENT"

//const options = "OPTIONS"
const Lock = "LOCK"
const manifestPrefix = "MANIFEST-"

// FileType represent a file type.
type FileType int

// File types.
const (
	TypeManifest FileType = 1 << iota
	TypeJournal
	TypeTable
	TypeTemp
	TypeStoreInfo
	TypeFamilyInfo

	TypeAll = TypeManifest | TypeJournal | TypeTable | TypeTemp | TypeStoreInfo | TypeFamilyInfo
)

type FileDesc struct {
	FileType   string
	FileNumber int64
}

// Get current file name for saving manifest file name
func Current() string {
	return "CURRENT"
}

// Table file name
func Table(fileNumber int64) string {
	return fmt.Sprintf("%06d.%s", fileNumber, sstSuffix)
}

func Info() string {
	return "info"
}

func SetCurrentFile(manifestFileNumber int64) error {
	manifest := manifestFileName(manifestFileNumber)
	tmp := fmt.Sprintf("%s.%s", current, tmpSuffix)
	if err := ioutil.WriteFile(tmp, []byte(manifest), 0666); err != nil {
		return err
	}
	if err := os.Rename(tmp, current); err != nil {
		return err
	}
	return nil
}

func manifestFileName(fileNumber int64) string {
	return fmt.Sprintf("%s%06d", manifestPrefix, fileNumber)
}
