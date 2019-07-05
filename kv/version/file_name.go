package version

import (
	"fmt"
)

const sstSuffix = "sst"
const TmpSuffix = "tmp"

const Options = "OPTIONS"
const manifestPrefix = "MANIFEST-"

// FileType represent a file type.
type FileType int

// File types.
const (
	TypeManifest FileType = iota
	TypeJournal
	TypeTable
	TypeTemp
	TypeInfo
)

// FileDesc define file type and file number
type FileDesc struct {
	FileType   string
	FileNumber int64
}

// current return current file name for saving manifest file name
func current() string {
	return "CURRENT"
}

// Table file name
func Table(fileNumber int64) string {
	return fmt.Sprintf("%06d.%s", fileNumber, sstSuffix)
}

// manifestFileName return manifest file name
func manifestFileName(fileNumber int64) string {
	return fmt.Sprintf("%s%06d", manifestPrefix, fileNumber)
}
