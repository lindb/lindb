package spi

type Split interface {
}

type BinarySplit struct {
	Page *Page
}
