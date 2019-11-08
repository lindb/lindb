package fileutil

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
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

// EncodeToml encodes data into file using toml format,
// encode data to tmp file, if success then rename tmp => target file
func EncodeToml(fileName string, v interface{}) error {
	tmp := fmt.Sprintf("%s.tmp", fileName)
	f, _ := os.Create(tmp)
	w := bufio.NewWriter(f)
	// write info using toml format
	if err := toml.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	if err := os.Rename(tmp, fileName); err != nil {
		return fmt.Errorf("rename tmp file[%s] name error:%s", tmp, err)
	}
	return nil
}

// DecodeToml decodes data from file using toml format
func DecodeToml(fileName string, v interface{}) error {
	if _, err := toml.DecodeFile(fileName, v); err != nil {
		return err
	}
	return nil
}
