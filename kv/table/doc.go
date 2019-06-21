package table

// Layout of sst-file:
//
// +-----------+
// | length    | 1-4 Byte
// +-----------+
// | value     | N Byte
// +-----------+
// | length    | 1-4 Byte
// +-----------+
// | value     | N Byte
// +-----------+
// | length    | 1-4 Byte
// +-----------+
// | value     | N Byte
// +-----------+
// .............
// +-----------+ <- posOfOffset
// | length    | 1-4 Byte
// +-----------+
// |  Offset   | N Byte
// +-----------+ <- posOfKeys
// | length    | 1-4 Byte
// +-----------+
// |  Keys     | N Byte
// +-----------+ <- Footer
// | length    | 1 Byte
// +-----------+
// |posOfOffset| 4 Byte
// +-----------+
// | posOfKeys | 4 Byte
// +-----------+
// |MagicNumber| 8 Byte
// +-----------+
