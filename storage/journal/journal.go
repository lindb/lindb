package journal

//import (
//	//"io/ioutil"
//)

//const (
//	//blockSize = 32 * 1024 // block size 32Kb
//)

// Writer writes journals(wal records) to an underlying io.Writer
type Writer struct {
	//bw *bufio.Writer
}

type Reader struct {
}

func NewWriter() *Writer {
	return &Writer{}
}

func NewReader(filename string) {
	//ioutil.ReadFile(filename)
	//os.Open()
}
func (w *Writer) Write(v []byte) error {

	return nil
}

func (r *Reader) Next() (bool, error) {

	return true, nil
}
