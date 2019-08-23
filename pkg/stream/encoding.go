package stream

// GetUVariantLength returns the length of variant-encoded-bytes.
func GetUVariantLength(value uint64) int {
	i := uint8(1)
	for ; ; i++ {
		if value < 2<<(i*7-1) {
			break
		}
	}
	return int(i)
}
