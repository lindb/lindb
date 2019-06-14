package io

import "os"

type FileWriter struct {
	f *os.File
}

func NewWriter(fileName string) (*FileWriter, error) {
	f, err := os.Create(fileName)
	if nil != err {
		return nil, err
	}
	return &FileWriter{
		f: f,
	}, nil
}

func (w *FileWriter) Write(data []byte) (int, error) {
	return w.f.Write(data)
}

func (w *FileWriter) Close() error {
	return w.f.Close()
}
