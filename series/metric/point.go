package metric

//
//import "github.com/lindb/lindb/series/field"
//
//type pointField struct {
//	name  []byte
//	typ   field.Type
//	value float64
//}
//
//type pointKeyValue struct {
//	key   []byte
//	value []byte
//}
//
//type Point struct {
//	name      []byte
//	namespace []byte
//	tags      []pointKeyValue
//	fields    []pointField
//	timestamp int64
//
//	tagIdx   int
//	fieldIdx int
//}
//
//func (p *Point) SetNamespace(namespace []byte) {
//
//}
//
//func (p *Point) Reset() {
//	p.name = p.name[:0]
//	p.namespace = p.namespace[:0]
//	p.timestamp = 0
//	p.tagIdx = 0
//	p.fieldIdx = 0
//}
