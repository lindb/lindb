package point

import "bytes"

type escapeSet struct {
	k   [1]byte
	esc [2]byte
}

var (
	metricNameEscapeCodes = [...]escapeSet{
		{k: [1]byte{','}, esc: [2]byte{'\\', ','}},
		{k: [1]byte{' '}, esc: [2]byte{'\\', ' '}},
	}
)

func UnescapeMetricName(in []byte) []byte {
	if bytes.IndexByte(in, '\\') == -1 {
		return in
	}

	for i := range metricNameEscapeCodes {
		c := &metricNameEscapeCodes[i]
		if bytes.IndexByte(in, c.k[0]) != -1 {
			in = bytes.Replace(in, c.esc[:], c.k[:], -1)
		}
	}
	return in
}

func EscapeMetricName(in []byte) []byte {
	for _, c := range metricNameEscapeCodes {
		if bytes.IndexByte(in, c.k[0]) != -1 {
			in = bytes.Replace(in, c.k[:], c.esc[:], -1)
		}
	}
	return in
}
