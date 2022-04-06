package packer

import (
	"fmt"
	"image"
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
	Scale   float64
}

func (a *atlas) CreateImage() (image.Image, error) {
	img := image.NewNRGBA(image.Rect(0, 0, a.Width, a.Height))

	// TODO run these draw steps in parallel
	for i := range a.Sprites {
		spr := a.Sprites[i].(*sprite)
		rect := image.Rect(spr.x, spr.y, spr.x+spr.w, spr.y+spr.h)

		assetReader, err := spr.Asset.Reader()
		if err != nil {
			return nil, fmt.Errorf("Failed to read asset '%s': %s", spr.path, err)
		}
		sprImg, _, err := image.Decode(assetReader)
		if err != nil {
			return nil, fmt.Errorf("Failed to decode asset '%s': %s", spr.path, err)
		}

		fastDraw(img, rect, sprImg)
	}

	return img, nil
}

func (a *atlas) Output(outputter Outputter, descriptorTemplate *template.Template) error {
	errc := make(chan error, 2)
	go func() {
		// Create and write the resulting image
		errc <- a.OutputImage(outputter, descriptorTemplate)
	}()
	go func() {
		// Create and write the file that describes the image
		errc <- a.OutputDesc(outputter, false, descriptorTemplate)
	}()
	// Drain error channel
	for i := 0; i < 2; i++ {
		if err := <-errc; err != nil {
			return err
		}
	}
	return nil
}

func (a *atlas) OutputImage(imageOutputter Outputter, descriptorTemplate *template.Template) error {
	// Create and write the resulting image
	return withFile(imageOutputter, a.ImageFilename, false, func(writer io.Writer) error {
		img, err := a.CreateImage()
		if err != nil {
			return err
		}
		return png.Encode(writer, img)
	})
}

func (a *atlas) OutputDesc(descOutputter Outputter, append bool, descriptorTemplate *template.Template) error {
	// Create and write the file that describes the image
	return withFile(descOutputter, a.DescFilename, append, func(writer io.Writer) error {
		return descriptorTemplate.Execute(writer, a)
	})
}
