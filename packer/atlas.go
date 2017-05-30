package packer

import (
	"image"
	"image/draw"
	"image/png"
	"io"
	"text/template"

	"github.com/RaniSputnik/packing"
)

type Atlas struct {
	Name    string
	Sprites []packing.Block

	DescFilename  string
	ImageFilename string

	Width  int
	Height int
}

func (a *Atlas) CreateImage() image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, a.Width, a.Height))

	for i := range a.Sprites {
		spr := a.Sprites[i].(*sprite)
		rect := image.Rect(spr.x, spr.y, spr.x+spr.w, spr.y+spr.h)
		draw.Draw(img, rect, spr.img, image.ZP, draw.Src)
	}

	return img
}

func (a *Atlas) Output(outputter Outputter, descriptorTemplate *template.Template) error {
	errc := make(chan error, 2)
	go func() {
		// Create and write the resulting image
		errc <- withFile(outputter, a.ImageFilename, func(writer io.Writer) error {
			img := a.CreateImage()
			return png.Encode(writer, img)
		})
	}()
	go func() {
		// Create and write the file that describes the image
		errc <- withFile(outputter, a.DescFilename, func(writer io.Writer) error {
			return descriptorTemplate.Execute(writer, a)
		})
	}()
	// Drain error channel
	for i := 0; i < 2; i++ {
		if err := <-errc; err != nil {
			return err
		}
	}
	return nil
}
