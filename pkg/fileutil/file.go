package fileutil

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// MkDirIfNotExist creates given dir if it not exist
func MkDirIfNotExist(path string) error {
	if !Exist(path) {
		if e := os.MkdirAll(path, os.ModePerm); e != nil {
			return e
		}
	}
	return nil
}

// RemoveDir deletes dir include children if exist
func RemoveDir(path string) error {
	if Exist(path) {
		if e := os.RemoveAll(path); e != nil {
			return e
		}
	}
	return nil
}

// MkDir create dir
func MkDir(path string) error {
	if e := os.MkdirAll(path, os.ModePerm); e != nil {
		return e
	}
	return nil
}

// ListDir reads the directory named by dirname and returns a list of filename.
func ListDir(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, file := range files {
		result = append(result, file.Name())
	}
	return result, nil
}

// Exist check file or dir if exist
func Exist(file string) bool {
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// GetExistPath get exist path based on given path
func GetExistPath(path string) string {
	if Exist(path) {
		return path
	}
	dir, _ := filepath.Split(path)
	length := len(dir)
	if length > 0 && os.IsPathSeparator(dir[length-1]) {
		dir = dir[:length-1]
	}
	return GetExistPath(dir)
}
