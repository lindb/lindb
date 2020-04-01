package stmt

// MetadataType represents metadata suggest type
type MetadataType uint8

// Defines all types of metadata suggest
const (
	Namespace MetadataType = iota + 1
	Metric
	TagKey
	TagValue
	Field
)

// Metadata represents search metadata statement
type Metadata struct {
	Namespace  string       // namespace
	MetricName string       // like table name
	Type       MetadataType // metadata suggest type
	TagKey     string
	TagValue   string
	Limit      int // result set limit
}
