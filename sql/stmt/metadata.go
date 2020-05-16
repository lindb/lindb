package stmt

// MetadataType represents metadata suggest type
type MetadataType uint8

// Defines all types of metadata suggest
const (
	Database MetadataType = iota + 1
	Namespace
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
	Prefix     string
	Condition  Expr // tag filter condition expression

	Limit int // result set limit
}
