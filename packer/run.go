package packer

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"sort"
	"text/template"

	"sync"

	"github.com/RaniSputnik/packing"
)

type Params struct {
	Name          string
	Input         AssetStreamer
	Output        Outputter
	Format        string
	Width, Height int
}

func Run(params *Params) error {
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	// Get the output format
	descFormat, err := GetFormatNamed(params.Format)
	if err != nil {
		return err
	}

	// Read the images from the input directory
	sprites, err := readAssetStream(ctx, params.Input)
	if err != nil {
		return err
	}

	// Arrange the images into the atlas space
	packer := &packing.BinPacker{}
	sort.Sort(packing.ByArea(sprites))
	err = packer.Fit(params.Width, params.Height, sprites...)
	if err != nil {
		return err
	}

	// TODO dynamically adjust number of atlases
	atlases := make([]*Atlas, 1)
	atlasName := fmt.Sprintf("%s-%d", params.Name, 1)
	atlases[0] = &Atlas{
		Name:         atlasName,
		Sprites:      sprites,
		DescFilename: fmt.Sprintf("%s.%s", atlasName, descFormat.Ext),
		// TODO add image type parameter
		ImageFilename: fmt.Sprintf("%s.%s", atlasName, "png"),
		Width:         params.Width,
		Height:        params.Height,
	}

	// TODO should be able to execute all atlases concurrently
	// TODO should write descriptor and image concurrently
	for _, atlas := range atlases {
		// Create and write the resulting image
		err = createImage(atlas, params.Output)
		if err != nil {
			return err
		}
		// Create and write the file that describes the image
		err = createDescriptor(descFormat.Template, atlas, params.Output)
		if err != nil {
			return err
		}
	}

	return nil
}

type assetDecodeResult struct {
	Sprite *sprite
	Err    error
}

func readAssetStream(ctx context.Context, assetStream AssetStreamer) ([]packing.Block, error) {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()
	// Stream the input
	assets, errc := assetStream.AssetStream(ctx)
	// Create decoder pool
	out := make(chan *assetDecodeResult)
	const numDecoders = 5
	var wg sync.WaitGroup
	wg.Add(numDecoders)
	for i := 0; i < numDecoders; i++ {
		go func() {
			decode(ctx, assets, out)
			wg.Done()
		}()
	}
	// Once all decoders complete, close the out channel
	go func() {
		wg.Wait()
		close(out)
	}()
	// Copy results from the out channel to the sprites slice
	var sprites []packing.Block
	for res := range out {
		if res.Err != nil {
			return nil, res.Err
		}
		sprites = append(sprites, res.Sprite)
	}
	// Check if the asset stream failed
	if err := <-errc; err != nil {
		return nil, err
	}

	return sprites, nil
}

// Decodes assets from the in channel and publishes the results to
// the out channel. Will continue even after errors have been discovered
// cancel the context to interrupt early.
func decode(ctx context.Context, in <-chan Asset, out chan<- *assetDecodeResult) {
	publishResult := func(spr *sprite, err error) {
		select {
		case out <- &assetDecodeResult{spr, err}:
		case <-ctx.Done():
		}
	}

	for asset := range in {
		reader, err := asset.CreateReader()
		if err != nil {
			publishResult(nil, err)
			continue
		}
		defer reader.Close()

		img, _, err := image.Decode(reader)
		if err != nil {
			publishResult(nil, err)
			continue
		}

		rect := img.Bounds()
		spr := &sprite{
			path: asset.Asset(),
			img:  img,
			w:    rect.Dx(),
			h:    rect.Dy(),
		}

		publishResult(spr, nil)
	}
}

func createImage(atlas *Atlas, outputter Outputter) error {
	img := atlas.CreateImage()

	writer, err := outputter.GetWriter(atlas.ImageFilename)
	if err != nil {
		return err
	}
	defer writer.Close()

	err = png.Encode(writer, img)
	if err != nil {
		return err
	}
	return nil
}

func createDescriptor(t *template.Template, atlas *Atlas, outputter Outputter) error {
	writer, err := outputter.GetWriter(atlas.DescFilename)
	if err != nil {
		return err
	}
	defer writer.Close()
	t.Execute(writer, atlas)
	return nil
}
