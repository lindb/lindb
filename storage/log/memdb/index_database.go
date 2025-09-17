package memdb

type IndexDatabase interface {
	GetOrCreateFieldValueID(namespace, fieldName, feldValue string) uint32
}
