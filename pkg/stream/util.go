package stream

import (
	"encoding/binary"
)

// ReadUint32 reads 4 bytes from buf as uint32
func ReadUint32(buf []byte, offset int) uint32 {
	return binary.LittleEndian.Uint32(buf[offset : offset+4])
}

// ReadUint16 reads 2 bytes from buf as uint16
func ReadUint16(buf []byte, offset int) uint16 {
	return binary.LittleEndian.Uint16(buf[offset : offset+2])
}

// ReadUvarint reads an encoded unsigned integer from bytes.Reader and returns it as a uint64.
func ReadUvarint(buf []byte, offset int) (value uint64, readLen int, err error) {
	var s uint
	for i := 0; ; i++ {
		b := buf[offset]
		offset++
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return value, i + 1, errOverflow
			}
			return value | uint64(b)<<s, i + 1, nil
		}
		value |= uint64(b&0x7f) << s
		s += 7
	}
}
