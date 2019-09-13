package series

import (
	"time"

	"github.com/lindb/lindb/pkg/timeutil"
)

// Version represents a metric version,
// it is the default created-time in milliseconds
type Version int64

// NewVersion returns a new Version
func NewVersion() Version {
	return Version(timeutil.Now())
}

// Int64 returns the version as int64
func (v Version) Int64() int64 {
	return int64(v)
}

// Time converts the version into Time
func (v Version) Time() time.Time {
	return time.Unix(0, v.Int64()*1000*1000)
}

// Elapsed returns the elapsed time since version start.
func (v Version) Elapsed() time.Duration {
	return time.Since(v.Time())
}

// IsExpired detects if this version has been expired
func (v Version) IsExpired(ttl time.Duration) bool {
	return v.Time().Add(ttl).Before(time.Now())
}

// Before returns if this version is before other.
func (v Version) Before(other Version) bool {
	return v < other
}

// Before returns if this version is after other.
func (v Version) After(other Version) bool {
	return v > other
}

// Before checks if this version is equals to other.
func (v Version) Equal(other Version) bool {
	return v == other
}

// String implements Stringer
func (v Version) String() string {
	return v.Time().Format("2006-01-02 15:04:05")
}
