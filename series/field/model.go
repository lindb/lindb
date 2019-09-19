package field

// Meta is the meta-data for field, which contains field-name, fieldID and field-type
type Meta struct {
	ID   uint16
	Type Type
	Name string
}
