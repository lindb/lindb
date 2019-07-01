package interval

import (
	"fmt"
	"time"

	"github.com/eleme/lindb/pkg/util"
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
func GetCalculator(intervalType Type) Calculator {
	return intervalTypes[intervalType]
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
	// CalSegmentTime calculates family base time based on given timestamp
	CalFamilyBaseTime(timestamp int64) int64
	// CalSlot calculates field store slot index based on given timestamp
	CalSlot(timestamp int64) int32
}

// day implements Calculator interface for day interval type
type day struct {
}

func (d *day) CalSlot(timestamp int64) int32 {
	//TODO
	return 0
}

// GetSegment returns segment name by given timestamp for day interval type
func (d *day) GetSegment(timestamp int64) string {
	return util.FormatTimestamp(timestamp, "20060102")
}

// ParseSegmentTime parses segment base time based on given segment name for day interval type
func (d *day) ParseSegmentTime(segmentName string) (int64, error) {
	return util.ParseTimestamp(segmentName, "20060102")
}

func (d *day) CalSegmentTime(timestamp int64) int64 {
	t := time.Unix(timestamp/1000, 0)
	t2 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return t2.UnixNano() / 1000000
}
func (d *day) CalFamilyBaseTime(timestamp int64) int64 {
	//TODO
	return 0
}

// month implements Calculator interface for month interval type
type month struct {
}

func (m *month) CalSlot(timestamp int64) int32 {
	return 0
}

// GetSegment returns segment name by given timestamp for month interval type
func (m *month) GetSegment(timestamp int64) string {
	return util.FormatTimestamp(timestamp, "200601")
}

// ParseSegmentTime parses segment base time based on given segment name for month interval type
func (m *month) ParseSegmentTime(segmentName string) (int64, error) {
	return util.ParseTimestamp(segmentName, "200601")
}
func (m *month) CalSegmentTime(timestamp int64) int64 {
	return 0
}
func (m *month) CalFamilyBaseTime(timestamp int64) int64 {
	//TODO
	return 0
}

// year implements Calculator interface for year interval type
type year struct {
}

func (y *year) CalSlot(timestamp int64) int32 {
	//TODO
	return 0
}

// GetSegment returns segment name by given timestamp for day interval type
func (y *year) GetSegment(timestamp int64) string {
	return util.FormatTimestamp(timestamp, "2006")
}

// ParseSegmentTime parses segment base time based on given segment name for year interval type
func (y *year) ParseSegmentTime(segmentName string) (int64, error) {
	return util.ParseTimestamp(segmentName, "2006")
}
func (y *year) CalSegmentTime(timestamp int64) int64 {
	return 0
}
func (y *year) CalFamilyBaseTime(timestamp int64) int64 {
	//TODO
	return 0
}
