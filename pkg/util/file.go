package util

import (
	"bufio"
	"os"

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

// Exist check file or dir if exist
func Exist(file string) bool {
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// EncodeToml encodes data into file using toml format
func EncodeToml(fileName string, v interface{}) error {
	f, _ := os.Create(fileName)
	w := bufio.NewWriter(f)
	if err := toml.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

// DecodeToml encodes data from file using toml format
func DecodeToml(fileName string, v interface{}) error {
	if _, err := toml.DecodeFile(fileName, v); err != nil {
		return err
	}
	return nil
}
