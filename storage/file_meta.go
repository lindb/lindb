package storage

type FileMeta struct {
	fileNumber int64
	minKey     uint32
	maxKey     uint32
	fileSize   int32
}

// New file meta
func NewFileMeta(fileNumber int64, minKey uint32, maxKey uint32, fileSize int32) *FileMeta {
	return &FileMeta{
		fileNumber: fileNumber,
		minKey:     minKey,
		maxKey:     maxKey,
		fileSize:   fileSize,
	}
}

func (f *FileMeta) GetFileNumber() int64 {
	return f.fileNumber
}

func (f *FileMeta) GetMinKey() uint32 {
	return f.minKey
}

func (f *FileMeta) GetMaxKey() uint32 {
	return f.maxKey
}

func (f *FileMeta) GetFileSize() int32 {
	return f.fileSize
}
