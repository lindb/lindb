package field

// Meta is the meta-data for field, which contains field-name, fieldID and field-type
type Meta struct {
	ID   uint16 // query not use ID, don't get id in query phase
	Type Type
	Name string
}
