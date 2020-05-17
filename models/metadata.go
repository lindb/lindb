package models

// Metadata represents metadata query result model
type Metadata struct {
	Type   string      `json:"type"`
	Values interface{} `json:"values"`
}

// Field represents field metadata
type Field struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
