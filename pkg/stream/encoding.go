package stream

// UvariantSize returns the bytes-size of a uint64 uvariant encoded number.
func UvariantSize(value uint64) int {
	i := 0
	for value >= 0x80 {
		value >>= 7
		i++
	}
	return i + 1
}

// VariantSize returns the bytes-size of a int64 variant encoded number.
func VariantSize(value int64) int {
	ux := uint64(value) << 1
	if value < 0 {
		ux = ^ux
	}
	return UvariantSize(ux)
}
