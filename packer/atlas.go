package packer

import (
	"image"
	"image/draw"
	"image/png"
	"io"
	"text/template"

	"github.com/RaniSputnik/lovepac/packing"
)

type atlas struct {
	Name    string
	Sprites []packing.Block

	DescFilename  string
	ImageFilename string

	Width   int
	Height  int
	Padding int
}

func (a *atlas) CreateImage() (image.Image, error) {
	img := image.NewNRGBA(image.Rect(0, 0, a.Width, a.Height))

	// TODO run these draw steps in parallel
	for i := range a.Sprites {
		spr := a.Sprites[i].(*sprite)
		rect := image.Rect(spr.x, spr.y, spr.x+spr.w, spr.y+spr.h)

		assetReader, err := spr.Asset.Reader()
		if err != nil {
			return nil, err
		}
		sprImg, _, err := image.Decode(assetReader)
		if err != nil {
			return nil, err
		}

		draw.Draw(img, rect, sprImg, image.ZP, draw.Src)
	}

	return img, nil
}

func (a *atlas) Output(outputter Outputter, descriptorTemplate *template.Template) error {
	errc := make(chan error, 2)
	go func() {
		// Create and write the resulting image
		errc <- withFile(outputter, a.ImageFilename, func(writer io.Writer) error {
			img, err := a.CreateImage()
			if err != nil {
				return err
			}
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
