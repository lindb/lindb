package edit

import "github.com/eleme/lindb/storage"

type AddFile struct {
	level int32
	file  *storage.FileMeta
}

func NewAddFile(level int32, file *storage.FileMeta) Log {
	return &AddFile{
		level: level,
		file:  file,
	}
}
func (a *AddFile) Name() string {
	return "a"
}
func (a *AddFile) Encode() {

}

func (a *AddFile) Decode() {

}

type DeleteFile struct {
	level      int32
	fileNumber int64
}

func NewDeleteFile(level int32, fileNumber int64) *DeleteFile {
	return &DeleteFile{
		level:      level,
		fileNumber: fileNumber,
	}
}
func (a *DeleteFile) Name() string {
	return "b"
}
func (a *DeleteFile) Encode() {

}

func (a *DeleteFile) Decode() {

}
