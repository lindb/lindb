package util

import (
	"os"
	"bufio"
	"github.com/BurntSushi/toml"
)

// Create given dir if it not exist
func MkDirIfNotExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if e := os.MkdirAll(path, os.ModePerm); e != nil {
			return e
		}
	}
	return nil
}

func MkDir(path string) error {
	if e := os.MkdirAll(path, os.ModePerm); e != nil {
		return e
	}
	return nil
}

// Check file or dir if exist
func Exist(file string) bool {
	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func EncodeToml(fileName string, v interface{}) error {
	f, _ := os.Create(fileName)
	w := bufio.NewWriter(f)
	if err := toml.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}
