package stmt

// MetadataType represents metadata suggest type
type MetadataType uint8

// Defines all types of metadata suggest
const (
	Metric MetadataType = iota + 1
	TagKey
	TagValue
)

// Metadata represents search metadata statement
type Metadata struct {
	MetricName string       // like table name
	Type       MetadataType // metadata suggest type
	TagKey     string
	TagValue   string
	Limit      int // result set limit
}
