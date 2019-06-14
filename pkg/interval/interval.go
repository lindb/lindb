package interval

type Type interface {
	// get segment name
	GetSegment(timestamp int64) string

	// cal segment base time based on given timestamp
	CalSegmentTime(timestamp int64) int64

	// cal family base time based on given timestamp
	CalFamilyBaseTime(timestamp int64) int64

	CalSlot(timestamp int64) int32
}

type day struct {
}

func NewDayInterval() Type {
	return &day{}
}

func (d *day) CalSlot(timestamp int64) int32 {
	return 0
}

func (d *day) GetSegment(timestamp int64) string {
	return ""
}
func (d *day) CalSegmentTime(timestamp int64) int64 {
	return 0
}
func (d *day) CalFamilyBaseTime(timestamp int64) int64 {
	return 0
}
