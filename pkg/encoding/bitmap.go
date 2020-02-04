package encoding

import "github.com/lindb/roaring"

// for testing
var (
	BitmapMarshal   = bitmapMarshal
	BitmapUnmarshal = bitmapUnmarshal
)

// bitmapMarshal marshals the bitmap data for testing
func bitmapMarshal(bitmap *roaring.Bitmap) ([]byte, error) {
	return bitmap.MarshalBinary()
}

// bitmapUnmarshal unmarshal the bitmap from data for testing
func bitmapUnmarshal(bitmap *roaring.Bitmap, data []byte) error {
	return bitmap.UnmarshalBinary(data)
}
