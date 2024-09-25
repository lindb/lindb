package spi

import "github.com/lindb/lindb/spi/types"

type Split interface{}

type BinarySplit struct {
	Page *types.Page
}
