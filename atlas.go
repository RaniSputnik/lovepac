package main

import (
	"github.com/RaniSputnik/packing"
)

type Atlas struct {
	Name    string
	Sprites []packing.Block

	DescPath  string
	ImagePath string

	Width  int
	Height int
}
