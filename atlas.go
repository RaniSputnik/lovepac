package main

import (
	"image"
)

type Atlas struct {
	Name  string
	Image image.Image
	Data  []byte
}
