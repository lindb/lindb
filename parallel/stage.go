package parallel

// Stage is the definition of query stage
type Stage int

const (
	Filtering Stage = iota + 1
	Grouping
	Scanner
	DownSampling
)

func (qs Stage) String() string {
	switch qs {
	case Filtering:
		return "filtering"
	case Grouping:
		return "grouping"
	case Scanner:
		return "scanner"
	case DownSampling:
		return "downSampling"
	default:
		return "unknown"
	}
}
