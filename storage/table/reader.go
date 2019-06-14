package table

type Reader interface {
	Get(key uint32) ([]byte, error)
}

type StoreReader struct {
}
