package meta

import (
	"fmt"
	"io/ioutil"
	"os"
)

const SST_SUFFIX = "sst"
const TMP_SUFFIX = "tmp"
const CURRENT = "CURRENT"
const OPTIONS = "OPTIONS"
const LOCK = "LOCK"
const MANIFEST_PREFIX = "MANIFEST-"

// FileType represent a file type.
type FileType int

// File types.
const (
	TypeManifest   FileType = 1 << iota
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
	return fmt.Sprintf("%06d.%s", fileNumber, SST_SUFFIX)
}

func Info() string {
	return "info"
}

func SetCurrentFile(manifestFileNumber int64) error {
	manifest := manifestFileName(manifestFileNumber)
	tmp := fmt.Sprintf("%s.%s", CURRENT, TMP_SUFFIX)
	if err := ioutil.WriteFile(tmp, []byte(manifest), 0666); err != nil {
		return err
	}
	if err := os.Rename(tmp, CURRENT); err != nil {
		return err
	}
	return nil
}

func manifestFileName(fileNumber int64) string {
	return fmt.Sprintf("%s%06d", MANIFEST_PREFIX, fileNumber)
}
