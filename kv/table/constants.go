package table

import (
	"errors"

	"github.com/lindb/lindb/pkg/logger"
)

var (
	ErrEmptyKeys = errors.New("empty keys under store builder")
)

const (
	// magic-number in the footer of sst file
	magicNumberOffsetFile uint64 = 0x69632d656d656c65
	// current file layout version
	version0 = 0

	sstFileFooterSize = 1 + // entry length wrote by bufioutil
		4 + // posOfOffset(4)
		4 + // posOfKeys(4)
		1 + // version(1)
		8 // magicNumber(8)
	// footer-size, offset(1), keys(1)
	sstFileMinLength = sstFileFooterSize + 2
)

var tableLogger = logger.GetLogger("kv", "table")
