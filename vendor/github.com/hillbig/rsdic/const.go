package rsdic

const (
	kSmallBlockSize          = 64
	kLargeBlockSize          = 1024
	kSelectBlockSize         = 4096
	kUseRawLen               = 48
	kSmallBlockPerLargeBlock = kLargeBlockSize / kSmallBlockSize
)
