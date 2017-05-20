package packer

import (
	"image"
	"image/draw"

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

func (a *Atlas) CreateImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, a.Width, a.Height))

	for i := range a.Sprites {
		spr := a.Sprites[i].(*sprite)
		rect := image.Rect(spr.x, spr.y, spr.x+spr.w, spr.y+spr.h)
		draw.Draw(img, rect, spr.img, image.ZP, draw.Src)
	}

	return img
}
