// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package point

import (
	"bytes"
	"fmt"

	"github.com/lindb/lindb/series/tag"
)

const ( // the number of characters for the largest possible int64 (9223372036854775807)
	maxInt64Digits = 19

	// the number of characters for the smallest possible int64 (-9223372036854775808)
	minInt64Digits = 20

	// the number of characters for the largest possible uint64 (18446744073709551615)
	maxUint64Digits = 20

	// the number of characters required for the largest float64 before a range check
	// would occur during parsing
	maxFloat64Digits = 25

	// the number of characters required for smallest float64 before a range check occur
	// would occur during parsing
	minFloat64Digits = 27
)

// parseTags parses buf to the dst slice, then return it,
// the returned slice may be a new one
func parseTags(buf []byte, dst tag.Tags) tag.Tags {
	if len(buf) == 0 {
		return nil
	}
	n := bytes.Count(buf, []byte(","))
	if cap(dst) < n {
		dst = make(tag.Tags, n)
	} else {
		dst = dst[:n]
	}
	if dst == nil {
		dst = tag.Tags{}
	}
	// Series keys can contain escaped commas, therefore the number of commas
	// in a series key only gives an estimation of the upper bound on the number
	// of tags.
	var i int
	walkTags(buf, func(key, value []byte) bool {
		dst[i].Key, dst[i].Value = key, value
		i++
		return true
	})
	return dst[:i]
}

func walkTags(buf []byte, fn func(key, value []byte) bool) {
	if len(buf) == 0 {
		return
	}

	pos, name := scanTo(buf, 0, ',')

	// it's an empty key, so there are no tags
	if len(name) == 0 {
		return
	}

	hasEscape := bytes.IndexByte(buf, '\\') != -1
	i := pos + 1
	var key, value []byte
	for {
		if i >= len(buf) {
			break
		}
		i, key = scanTo(buf, i, '=')
		i, value = scanTagValue(buf, i+1)

		if len(value) == 0 {
			continue
		}

		if hasEscape {
			if !fn(tag.UnescapeTag(key), tag.UnescapeTag(value)) {
				return
			}
		} else {
			if !fn(key, value) {
				return
			}
		}

		i++
	}
}

// MakeKey creates a key for a set of tags.
func MakeKey(name []byte, tags tag.Tags) []byte {
	return AppendMakeKey(nil, name, tags)
}

// AppendMakeKey appends the key derived from name and tags to dst and returns the extended buffer.
func AppendMakeKey(dst []byte, name []byte, tags tag.Tags) []byte {
	// unescape the name and then re-escape it to avoid double escaping.
	// The key should always be stored in escaped form.
	dst = append(dst, EscapeMetricName(UnescapeMetricName(name))...)
	dst = tags.AppendHashKey(dst)
	return dst
}

// scanTo returns the end position in the buf and the next consecutive block
// of bytes
func scanTo(buf []byte, i int, stop byte) (int, []byte) {
	start := i
	for {
		// reached the end of buf?
		if i >= len(buf) {
			break
		}

		// Reached unescaped stop value?
		if buf[i] == stop && (i == 0 || buf[i-1] != '\\') {
			break
		}
		i++
	}

	return i, buf[start:i]
}

func scanTagValue(buf []byte, i int) (int, []byte) {
	start := i
	for {
		if i >= len(buf) {
			break
		}

		if buf[i] == ',' && buf[i-1] != '\\' {
			break
		}
		i++
	}
	if i > len(buf) {
		return i, nil
	}
	return i, buf[start:i]
}

// scanLine returns the end position in buf and the next line found within
// buf.
func scanLine(buf []byte, i int) (int, []byte) {
	start := i
	quoted := false
	fields := false

	// tracks how many '=' and commas we've seen
	// this duplicates some of the functionality in scanFields
	equals := 0
	commas := 0
	for {
		// reached the end of buf?
		if i >= len(buf) {
			break
		}

		// skip past escaped characters
		if buf[i] == '\\' && i+2 < len(buf) {
			i += 2
			continue
		}

		if buf[i] == ' ' {
			fields = true
		}

		// If we see a double quote, makes sure it is not escaped
		if fields {
			switch buf[i] {
			case '=':
				if !quoted {
					i++
					equals++
					continue
				}
			case ',':
				if !quoted {
					i++
					commas++
					continue
				}
			case '"':
				if equals > commas {
					i++
					quoted = !quoted
					continue
				}
			}
		}

		if buf[i] == '\n' && !quoted {
			break
		}

		i++
	}

	return i, buf[start:i]
}

// skipWhitespace returns the end position within buf, starting at i after
// scanning over spaces in tags.
func skipWhitespace(buf []byte, i int) int {
	for i < len(buf) {
		if buf[i] != ' ' && buf[i] != '\t' && buf[i] != 0 {
			break
		}
		i++
	}
	return i
}

// scanKey scans buf starting at i for the metric-name and tag portion of the point.
// It returns the ending position and the byte slice of key within buf.  If there
// are tags, they will be sorted if they are not already.
func scanKey(buf []byte, i int) (int, []byte, error) {
	start := skipWhitespace(buf, i)

	i = start

	// Determines whether the tags are sort, assume they are
	sorted := true

	// indices holds the indexes within buf of the start of each tag.  For example,
	// a buf of 'cpu,host=a,region=b,zone=c' would have indices slice of [4,11,20]
	// which indicates that the first tag starts at buf[4], seconds at buf[11], and
	// last at buf[20]
	indices := make([]int, 100)

	// tracks how many commas we've seen so we know how many values are indices.
	// Since indices is an arbitrarily large slice,
	// we need to know how many values in the buffer are in use.
	commas := 0

	// First scan the Point's metric-name.
	state, i, err := scanMetricName(buf, i)
	if err != nil {
		return i, buf[start:i], err
	}

	// Optionally scan tags if needed.
	if state == tagKeyState {
		i, commas, indices, err = scanTags(buf, i, indices)
		if err != nil {
			return i, buf[start:i], err
		}
	}

	// Now we know where the key region is within buf, and the location of tags, we
	// need to determine if duplicate tags exist and if the tags are sorted. This iterates
	// over the list comparing each tag in the sequence with each other.
	for j := 0; j < commas-1; j++ {
		// get the left and right tags
		_, left := scanTo(buf[indices[j]:indices[j+1]-1], 0, '=')
		_, right := scanTo(buf[indices[j+1]:indices[j+2]-1], 0, '=')

		// If left is greater than right, the tags are not sorted. We do not have to
		// continue because the short path no longer works.
		// If the tags are equal, then there are duplicate tags, and we should abort.
		// If the tags are not sorted, this pass may not find duplicate tags and we
		// need to do a more exhaustive search later.
		if cmp := bytes.Compare(left, right); cmp > 0 {
			sorted = false
			break
		} else if cmp == 0 {
			return i, buf[start:i], ErrDuplicateTags
		}
	}

	// If the tags are not sorted, then sort them.  This sort is inline and
	// uses the tag indices we created earlier.  The actual buffer is not sorted, the
	// indices are using the buffer for value comparison.  After the indices are sorted,
	// the buffer is reconstructed from the sorted indices.
	if !sorted && commas > 0 {
		// Get the measurement name for later
		measurement := buf[start : indices[0]-1]

		// Sort the indices
		indices := indices[:commas]
		insertionSort(0, commas, buf, indices)

		// Create a new key using the measurement and sorted indices
		b := make([]byte, len(buf[start:i]))
		pos := copy(b, measurement)
		for _, i := range indices {
			b[pos] = ','
			pos++
			_, v := scanToSpaceOr(buf, i, ',')
			pos += copy(b[pos:], v)
		}

		// Check again for duplicate tags now that the tags are sorted.
		for j := 0; j < commas-1; j++ {
			// get the left and right tags
			_, left := scanTo(buf[indices[j]:], 0, '=')
			_, right := scanTo(buf[indices[j+1]:], 0, '=')

			// If the tags are equal, then there are duplicate tags, and we should abort.
			// If the tags are not sorted, this pass may not find duplicate tags and we
			// need to do a more exhaustive search later.
			if bytes.Equal(left, right) {
				return i, b, fmt.Errorf("duplicate tags")
			}
		}

		return i, b, nil
	}

	return i, buf[start:i], nil
}

// The following constants allow us to specify which state to move to
// next, when scanning sections of a Point.
const (
	tagKeyState = iota
	tagValueState
	fieldsState
)

// scanMetricName examines the metric-name part of a Point, returning
// the next state to move to, and the current location in the buffer.
func scanMetricName(buf []byte, i int) (int, int, error) {
	// Check first byte of metric, anything except a comma is fine.
	// It can't be a space, since whitespace is stripped prior to this
	// function call.
	if i >= len(buf) || buf[i] == ',' {
		return -1, i, fmt.Errorf("missing metric-name")
	}

	for {
		i++
		if i >= len(buf) {
			// cpu
			return -1, i, fmt.Errorf("missing fields")
		}

		if buf[i-1] == '\\' {
			// Skip character (it's escaped).
			continue
		}

		// Unescaped comma; move onto scanning the tags.
		if buf[i] == ',' {
			return tagKeyState, i + 1, nil
		}

		// Unescaped space; move onto scanning the fields.
		if buf[i] == ' ' {
			// cpu value=1.0
			return fieldsState, i, nil
		}
	}
}

// scanTags examines all the tags in a Point, keeping track of and
// returning the updated indices slice, number of commas and location
// in buf where to start examining the Point fields.
func scanTags(buf []byte, i int, indices []int) (int, int, []int, error) {
	var (
		err    error
		commas int
		state  = tagKeyState
	)

	for {
		switch state {
		case tagKeyState:
			// Grow our indices slice if we have too many tags.
			if commas >= len(indices) {
				newIndices := make([]int, cap(indices)*2)
				copy(newIndices, indices)
				indices = newIndices
			}
			indices[commas] = i
			commas++

			i, err = scanTagsKey(buf, i)
			state = tagValueState // tag value always follows a tag key
		case tagValueState:
			state, i, err = scanTagsValue(buf, i)
		case fieldsState:
			indices[commas] = i + 1
			return i, commas, indices, nil
		}

		if err != nil {
			return i, commas, indices, err
		}
	}
}

// scanTagsKey scans each character in a tag key.
func scanTagsKey(buf []byte, i int) (int, error) {
	// First character of the key.
	if i >= len(buf) || buf[i] == ' ' || buf[i] == ',' || buf[i] == '=' {
		// cpu,{'', ' ', ',', '='}
		return i, fmt.Errorf("missing tag key")
	}

	// Examine each character in the tag key until we hit an unescaped
	// equals (the tag value), or we hit an error (i.e., unescaped
	// space or comma).
	for {
		i++

		// Either we reached the end of the buffer or we hit an
		// unescaped comma or space.
		if i >= len(buf) ||
			((buf[i] == ' ' || buf[i] == ',') && buf[i-1] != '\\') {
			// cpu,tag{'', ' ', ','}
			return i, fmt.Errorf("missing tag value")
		}

		if buf[i] == '=' && buf[i-1] != '\\' {
			// cpu,tag=
			return i + 1, nil
		}
	}
}

// scanTagsValue scans each character in a tag value.
func scanTagsValue(buf []byte, i int) (int, int, error) {
	// Tag value cannot be empty.
	if i >= len(buf) || buf[i] == ',' || buf[i] == ' ' {
		// cpu,tag={',', ' '}
		return -1, i, ErrMissingTagValue
	}

	// Examine each character in the tag value until we hit an unescaped
	// comma (move onto next tag key), an unescaped space (move onto
	// fields), or we error out.
	for {
		i++
		if i >= len(buf) {
			// cpu,tag=value
			return -1, i, ErrMissingFields
		}

		// An unescaped equals sign is an invalid tag value.
		if buf[i] == '=' && buf[i-1] != '\\' {
			// cpu,tag={'=', 'fo=o'}
			return -1, i, fmt.Errorf("invalid tag format")
		}

		if buf[i] == ',' && buf[i-1] != '\\' {
			// cpu,tag=foo,
			return tagKeyState, i + 1, nil
		}

		// cpu,tag=foo value=1.0
		// cpu, tag=foo\= value=1.0
		if buf[i] == ' ' && buf[i-1] != '\\' {
			return fieldsState, i, nil
		}
	}
}

func insertionSort(l, r int, buf []byte, indices []int) {
	for i := l + 1; i < r; i++ {
		for j := i; j > l && less(buf, indices, j, j-1); j-- {
			indices[j], indices[j-1] = indices[j-1], indices[j]
		}
	}
}

func less(buf []byte, indices []int, i, j int) bool {
	// This grabs the tag names for i & j, it ignores the values
	_, a := scanTo(buf, indices[i], '=')
	_, b := scanTo(buf, indices[j], '=')
	return bytes.Compare(a, b) < 0
}

// scanTo returns the end position in buf and the next consecutive block
// of bytes, starting from i and ending with stop byte.  If there are leading
// spaces, they are skipped.
func scanToSpaceOr(buf []byte, i int, stop byte) (int, []byte) {
	start := i
	if buf[i] == stop || buf[i] == ' ' {
		return i, buf[start:i]
	}

	for {
		i++
		if buf[i-1] == '\\' {
			continue
		}

		// reached the end of buf?
		if i >= len(buf) {
			return i, buf[start:i]
		}

		// reached end of block?
		if buf[i] == stop || buf[i] == ' ' {
			return i, buf[start:i]
		}
	}
}

// scanFields scans buf, starting at i for the fields section of a point.  It returns
// the ending position and the byte slice of the fields within buf.
func scanFields(buf []byte, i int) (int, []byte, error) {
	start := skipWhitespace(buf, i)
	i = start
	quoted := false

	// tracks how many '=' we've seen
	equals := 0

	// tracks how many commas we've seen
	commas := 0

	for {
		// reached the end of buf?
		if i >= len(buf) {
			break
		}

		// escaped characters?
		if buf[i] == '\\' && i+1 < len(buf) {
			i += 2
			continue
		}

		// If the value is quoted, scan until we get to the end quote
		// Only quote values in the field value since quotes are not significant
		// in the field key
		if buf[i] == '"' && equals > commas {
			quoted = !quoted
			i++
			continue
		}

		// If we see an =, ensure that there is at least on char before and after it
		if buf[i] == '=' && !quoted {
			equals++

			// check for "... =123" but allow "a\ =123"
			if buf[i-1] == ' ' && buf[i-2] != '\\' {
				return i, buf[start:i], ErrMissingFieldName
			}

			// check for "...a=123,=456" but allow "a=123,a\,=456"
			if buf[i-1] == ',' && buf[i-2] != '\\' {
				return i, buf[start:i], ErrMissingFieldName
			}

			// check for "... value="
			if i+1 >= len(buf) {
				return i, buf[start:i], ErrMissingFieldValue
			}

			// check for "... value=,value2=..."
			if buf[i+1] == ',' || buf[i+1] == ' ' {
				return i, buf[start:i], ErrMissingFieldValue
			}

			if isNumeric(buf[i+1]) || buf[i+1] == '-' || buf[i+1] == 'N' || buf[i+1] == 'n' {
				var err error
				i, err = scanNumber(buf, i+1)
				if err != nil {
					return i, buf[start:i], err
				}
				continue
			}
			// If next byte is not a double-quote, the value must be a boolean
			if buf[i+1] != '"' {
				return i, buf[start:i], fmt.Errorf("invalid field value")
			}
		}

		if buf[i] == ',' && !quoted {
			commas++
		}

		// reached end of block?
		if buf[i] == ' ' && !quoted {
			break
		}
		i++
	}

	if quoted {
		return i, buf[start:i], fmt.Errorf("unbalanced quotes")
	}

	// check that all field sections had key and values (e.g. prevent "a=1,b"
	if equals == 0 || commas != equals-1 {
		return i, buf[start:i], fmt.Errorf("invalid field format")
	}

	return i, buf[start:i], nil
}

// scanTime scans buf, starting at i for the time section of a point. It
// returns the ending position and the byte slice of the timestamp within buf
// and and error if the timestamp is not in the correct numeric format.
func scanTime(buf []byte, i int) (int, []byte, error) {
	start := skipWhitespace(buf, i)
	i = start

	for {
		// reached the end of buf?
		if i >= len(buf) {
			break
		}

		// Reached end of block or trailing whitespace?
		if buf[i] == '\n' || buf[i] == ' ' {
			break
		}

		// Handle negative timestamps
		if i == start && buf[i] == '-' {
			i++
			continue
		}

		// Timestamps should be integers, make sure they are so we don't need
		// to actually  parse the timestamp until needed.
		if buf[i] < '0' || buf[i] > '9' {
			return i, buf[start:i], fmt.Errorf("bad timestamp")
		}
		i++
	}
	return i, buf[start:i], nil
}

func isNumeric(b byte) bool {
	return (b >= '0' && b <= '9') || b == '.'
}

// scanNumber returns the end position within buf, start at i after
// scanning over buf for an integer, or float.  It returns an
// error if a invalid number is scanned.
func scanNumber(buf []byte, i int) (int, error) {
	start := i
	var isInt, isUnsigned bool

	// Is negative number?
	if i < len(buf) && buf[i] == '-' {
		i++
		// There must be more characters now, as just '-' is illegal.
		if i == len(buf) {
			return i, ErrInvalidNumber
		}
	}

	// how many decimal points we've see
	decimal := false

	// indicates the number is float in scientific notation
	scientific := false

	for {
		if i >= len(buf) {
			break
		}

		if buf[i] == ',' || buf[i] == ' ' {
			break
		}

		if buf[i] == 'i' && i > start && !(isInt || isUnsigned) {
			isInt = true
			i++
			continue
		} else if buf[i] == 'u' && i > start && !(isInt || isUnsigned) {
			isUnsigned = true
			i++
			continue
		}

		if buf[i] == '.' {
			// Can't have more than 1 decimal (e.g. 1.1.1 should fail)
			if decimal {
				return i, ErrInvalidNumber
			}
			decimal = true
		}

		// `e` is valid for floats but not as the first char
		if i > start && (buf[i] == 'e' || buf[i] == 'E') {
			scientific = true
			i++
			continue
		}

		// + and - are only valid at this point if they follow an e (scientific notation)
		if (buf[i] == '+' || buf[i] == '-') && (buf[i-1] == 'e' || buf[i-1] == 'E') {
			i++
			continue
		}

		// NaN is an unsupported value
		if i+2 < len(buf) && (buf[i] == 'N' || buf[i] == 'n') {
			return i, ErrInvalidNumber
		}

		if !isNumeric(buf[i]) {
			return i, ErrInvalidNumber
		}
		i++
	}

	if (isInt || isUnsigned) && (decimal || scientific) {
		return i, ErrInvalidNumber
	}

	numericDigits := i - start
	if isInt {
		numericDigits--
	}
	if decimal {
		numericDigits--
	}
	if buf[start] == '-' {
		numericDigits--
	}

	if numericDigits == 0 {
		return i, ErrInvalidNumber
	}

	// It's more common that numbers will be within min/max range for their type but we need to prevent
	// out or range numbers from being parsed successfully.  This uses some simple heuristics to decide
	// if we should parse the number to the actual type.  It does not do it all the time because it incurs
	// extra allocations and we end up converting the type again when writing points to disk.
	switch {
	case isInt:
		// Make sure the last char is an 'i' for integers (e.g. 9i10 is not valid)
		if buf[i-1] != 'i' {
			return i, ErrInvalidNumber
		}
		// Parse the int to check bounds the number of digits could be larger than the max range
		// We subtract 1 from the index to remove the `i` from our tests
		if len(buf[start:i-1]) >= maxInt64Digits || len(buf[start:i-1]) >= minInt64Digits {
			if _, err := parseInt64Bytes(buf[start : i-1]); err != nil {
				return i, fmt.Errorf("unable to parse integer %s: %s", buf[start:i-1], err)
			}
		}
	case isUnsigned:
		// Make sure the last char is a 'u' for unsigned
		if buf[i-1] != 'u' {
			return i, ErrInvalidNumber
		}
		// Make sure the first char is not a '-' for unsigned
		if buf[start] == '-' {
			return i, ErrInvalidNumber
		}
		// Parse the uint to check bounds the number of digits could be larger than the max range
		// We subtract 1 from the index to remove the `u` from our tests
		if len(buf[start:i-1]) >= maxUint64Digits {
			if _, err := parseUint64Bytes(buf[start : i-1]); err != nil {
				return i, fmt.Errorf("unable to parse unsigned %s: %s", buf[start:i-1], err)
			}
		}
	case scientific || len(buf[start:i]) >= maxFloat64Digits || len(buf[start:i]) >= minFloat64Digits:
		// Parse the float to check bounds if it's scientific or the number of digits could be larger than the max range
		if _, err := parseFloat64Bytes(buf[start:i]); err != nil {
			return i, fmt.Errorf("invalid float")
		}
	}
	return i, nil
}

// walkFields walks each field key and value via fn.  If fn returns false, the iteration
// is stopped.  The values are the raw byte slices and not the converted types.
func walkFields(buf []byte, fn func(key, value, data []byte) bool) error {
	var i int
	var key, val []byte
	for len(buf) > 0 {
		data := buf

		i, key = scanTo(buf, 0, '=')
		if i > len(buf)-2 {
			return fmt.Errorf("invalid value: field-key=%s", key)
		}
		buf = buf[i+1:]
		i, val = scanFieldValue(buf, 0)
		buf = buf[i:]
		if !fn(key, val, data[:len(data)-len(buf)]) {
			break
		}

		// slice off comma
		if len(buf) > 0 {
			buf = buf[1:]
		}
	}
	return nil
}
