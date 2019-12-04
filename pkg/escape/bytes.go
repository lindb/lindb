// Package escape contains utilities for escaping an InfluxDB-like
// LinDB line protocol.
//
// Reference: github.com/influxdata/influxdb/pkg/escape
package escape

import (
	"bytes"
	"strings"
)

// escapeItem is a escape mapping from byte to slice.
type escapeItem struct {
	From byte
	To   []byte
}

var escapeMappings = []escapeItem{
	{',', []byte(`\,`)},
	{'"', []byte(`\"`)},
	{' ', []byte(`\ `)},
	{'=', []byte(`\=`)},
}

// Bytes escape characters on the input slice
func Bytes(in []byte) (to []byte) {
	for _, mapping := range escapeMappings {
		in = bytes.Replace(in, []byte{mapping.From}, mapping.To, -1)
	}
	return in
}

const escapeChars = `," =`

// IsEscaped returns whether b has any escaped characters,
// i.e. whether b seems to have been processed by Bytes.
func IsEscaped(b []byte) bool {
	for len(b) > 0 {
		i := bytes.IndexByte(b, '\\')
		if i < 0 {
			return false
		}

		if i+1 < len(b) && strings.IndexByte(escapeChars, b[i+1]) >= 0 {
			return true
		}
		b = b[i+1:]
	}
	return false
}

// AppendUnescaped appends the unescaped version of src to dst
// and returns the resulting slice.
func AppendUnescaped(dst, src []byte) []byte {
	var pos int
	for len(src) > 0 {
		next := bytes.IndexByte(src[pos:], '\\')
		if next < 0 || pos+next+1 >= len(src) {
			return append(dst, src...)
		}

		if pos+next+1 < len(src) && strings.IndexByte(escapeChars, src[pos+next+1]) >= 0 {
			if pos+next > 0 {
				dst = append(dst, src[:pos+next]...)
			}
			src = src[pos+next+1:]
			pos = 0
		} else {
			pos += next + 1
		}
	}

	return dst
}

// Unescape returns a new slice containing the unescaped version of in.
func Unescape(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}

	if bytes.IndexByte(in, '\\') == -1 {
		return in
	}

	i := 0
	inLen := len(in)

	// The output size will be no more than inLen. Preallocating the
	// capacity of the output is faster and uses less memory than
	// letting append() do its own (over)allocation.
	out := make([]byte, 0, inLen)

	for {
		if i >= inLen {
			break
		}
		if in[i] == '\\' && i+1 < inLen {
			switch in[i+1] {
			case ',':
				out = append(out, ',')
				i += 2
				continue
			case '"':
				out = append(out, '"')
				i += 2
				continue
			case ' ':
				out = append(out, ' ')
				i += 2
				continue
			case '=':
				out = append(out, '=')
				i += 2
				continue
			}
		}
		out = append(out, in[i])
		i++
	}
	return out
}
