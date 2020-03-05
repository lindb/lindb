package version

import "github.com/lindb/lindb/kv/table"

// level stores sst files of level
type level struct {
	files map[table.FileNumber]*FileMeta
}

// newLevel new level instance
func newLevel() *level {
	return &level{
		files: make(map[table.FileNumber]*FileMeta),
	}
}

// addFile adds new file into file list
func (l *level) addFile(file *FileMeta) {
	l.files[file.GetFileNumber()] = file
}

// addFiles adds new files into file list
func (l *level) addFiles(files ...*FileMeta) {
	for _, file := range files {
		l.addFile(file)
	}
}

// deleteFile removes file from file list
func (l *level) deleteFile(fileNumber table.FileNumber) {
	delete(l.files, fileNumber)
}

// getFiles returns all files in current level
func (l *level) getFiles() []*FileMeta {
	var values []*FileMeta
	for _, v := range l.files {
		values = append(values, v)
	}
	return values
}

// numberOfFiles returns the number of files in current level
func (l *level) numberOfFiles() int {
	return len(l.files)
}
