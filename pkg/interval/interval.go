package interval

import (
	"fmt"
	"time"

	"github.com/eleme/lindb/pkg/timeutil"
)

// Type defines interval type
type Type int

const dayStr = "day"
const monthStr = "month"
const yearStr = "year"

// Interval types.
const (
	Day Type = iota + 1
	Month
	Year
	Unknown
)

// String returns string value of interval type
func (t Type) String() string {
	switch t {
	case Month:
		return monthStr
	case Year:
		return yearStr
	default:
		return dayStr
	}
}

// ParseType returns interval type based on string value, return error if type not in type list
func ParseType(s string) (Type, error) {
	switch s {
	case dayStr:
		return Day, nil
	case monthStr:
		return Month, nil
	case yearStr:
		return Year, nil
	default:
		return Unknown, fmt.Errorf("unknown interval type[%s]", s)
	}
}

// intervalTypes defines calculator for interval type
var intervalTypes = make(map[Type]Calculator)

// register adds calculator for interval type
func register(intervalType Type, calc Calculator) {
	if _, ok := intervalTypes[intervalType]; ok {
		panic(fmt.Sprintf("calculator of interval type already registered: %d", intervalType))
	}
	intervalTypes[intervalType] = calc
}

// GetCalculator returns calculator for given interval type
func GetCalculator(intervalType Type) (Calculator, error) {
	calc, ok := intervalTypes[intervalType]
	if !ok {
		return nil, fmt.Errorf("cannot found interval calculator by type:%d", intervalType)
	}
	return calc, nil
}

// init register interval types when system init
func init() {
	register(Day, &day{})
	register(Month, &month{})
	register(Year, &year{})
}

// Calculator represents calculate timestamp for each interval type
type Calculator interface {
	// GetSegment returns segment name by given timestamp
	GetSegment(timestamp int64) string
	// ParseSegmentTime parses segment base time based on given segment name
	ParseSegmentTime(segmentName string) (int64, error)
	// CalSegmentTime calculates segment base time based on given segment name
	CalSegmentTime(timestamp int64) int64
	// CalFamily calculates family base time based on given timestamp
	CalFamily(timestamp int64, segmentTime int64) int
	// CalFamilyStartTime calculates famliy start time based on segment time and family
	CalFamilyStartTime(segmentTime int64, family int) int64
	// CalSlot calculates field store slot index based on given timestamp and base time
	CalSlot(timestamp, baseTime, interval int64) int
}

// day implements Calculator interface for day interval type
type day struct {
}

// CalSlot calculates field store slot index based on given timestamp and base time for day interval type
func (d *day) CalSlot(timestamp, baseTime, interval int64) int {
	return int(((timestamp - baseTime) % timeutil.OneHour) / interval)
}

// GetSegment returns segment name by given timestamp for day interval type
func (d *day) GetSegment(timestamp int64) string {
	return timeutil.FormatTimestamp(timestamp, "20060102")
}

// ParseSegmentTime parses segment base time based on given segment name for day interval type
func (d *day) ParseSegmentTime(segmentName string) (int64, error) {
	return timeutil.ParseTimestamp(segmentName, "20060102")
}

// CalSegmentTime calculates segment base time based on given segment name for day interval type
func (d *day) CalSegmentTime(timestamp int64) int64 {
	t := time.Unix(timestamp/1000, 0)
	t2 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}

// CalFamily calculates family base time based on given timestamp for day interval type
func (d *day) CalFamily(timestamp int64, segmentTime int64) int {
	return int((timestamp - segmentTime) / timeutil.OneHour)
}

// CalFamilyStartTime calculates famliy start time based on segment time and family for day interval type
func (d *day) CalFamilyStartTime(segmentTime int64, family int) int64 {
	return segmentTime + int64(family)*timeutil.OneHour
}

// month implements Calculator interface for month interval type
type month struct {
}

// CalSlot calculates field store slot index based on given timestamp and base time for month interval type
func (m *month) CalSlot(timestamp, baseTime, interval int64) int {
	return int(((timestamp - baseTime) % timeutil.OneDay) / interval)
}

// GetSegment returns segment name by given timestamp for month interval type
func (m *month) GetSegment(timestamp int64) string {
	return timeutil.FormatTimestamp(timestamp, "200601")
}

// ParseSegmentTime parses segment base time based on given segment name for month interval type
func (m *month) ParseSegmentTime(segmentName string) (int64, error) {
	return timeutil.ParseTimestamp(segmentName, "200601")
}

// CalSegmentTime calculates segment base time based on given segment name for month interval type
func (m *month) CalSegmentTime(timestamp int64) int64 {
	t := time.Unix(timestamp/1000, 0)
	t2 := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}

// CalFamily calculates family base time based on given timestamp for month interval type
func (m *month) CalFamily(timestamp int64, segmentTime int64) int {
	t := time.Unix(timestamp/1000, 0)
	return t.Day()
}

// CalFamilyStartTime calculates famliy start time based on segment time and family for month interval type
func (m *month) CalFamilyStartTime(segmentTime int64, family int) int64 {
	t := time.Unix(segmentTime/1000, 0)
	t2 := time.Date(t.Year(), t.Month(), family, 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}

// year implements Calculator interface for year interval type
type year struct {
}

// CalSlot calculates field store slot index based on given timestamp and base time for year interval type
func (y *year) CalSlot(timestamp, baseTime, interval int64) int {
	return int((timestamp - baseTime) / interval)
}

// GetSegment returns segment name by given timestamp for day interval type
func (y *year) GetSegment(timestamp int64) string {
	return timeutil.FormatTimestamp(timestamp, "2006")
}

// ParseSegmentTime parses segment base time based on given segment name for year interval type
func (y *year) ParseSegmentTime(segmentName string) (int64, error) {
	return timeutil.ParseTimestamp(segmentName, "2006")
}

// CalSegmentTime calculates segment base time based on given segment name for year interval type
func (y *year) CalSegmentTime(timestamp int64) int64 {
	t := time.Unix(timestamp/1000, 0)
	t2 := time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}

// CalFamily calculates family base time based on given timestamp for year interval type
func (y *year) CalFamily(timestamp int64, segmentTime int64) int {
	t := time.Unix(timestamp/1000, 0)
	return int(t.Month())
}

// CalFamilyStartTime calculates famliy start time based on segment time and family for year interval type
func (y *year) CalFamilyStartTime(segmentTime int64, family int) int64 {
	t := time.Unix(segmentTime/1000, 0)
	t2 := time.Date(t.Year(), time.Month(family), 1, 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}
