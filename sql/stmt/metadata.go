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

// String returns string value of metadata type
func (m MetadataType) String() string {
	switch m {
	case Database:
		return "database"
	case Namespace:
		return "namespace"
	case Metric:
		return "measurement"
	case Field:
		return field
	case TagKey:
		return "tagKey"
	case TagValue:
		return "tagValue"
	default:
		return unknown
	}
}

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
