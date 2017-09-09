package packer

import (
	"context"
	"errors"
	"fmt"
	"image"
	"sort"

	"sync"

	"github.com/RaniSputnik/lovepac/packing"
	"github.com/RaniSputnik/lovepac/target"
)

var (
	// DefaultAtlasName is the default base name for
	// outputted files when no name is provided
	DefaultAtlasName = "atlas"
	// DefaultAtlasWidth is the width used if no width is specified
	DefaultAtlasWidth = 2048
	// DefaultAtlasHeight is the height used if no height is specified
	DefaultAtlasHeight = 2048
)

// Params are provided to the Run method to configure
// the texture packing output. Input, Ouput and Format parameters are
// required all other parameters are optional. You can use the public
// 'Default' properties to configure the defaults used when parameters
// are missing.
//
// Name is the name that will be prepended to the atlas files
// outputted. Eg. a value of "myatlas" would result in "myatlas-1.png"
//
// Input is used to provide readers for the assets that will be packed.
// In most cases packer.NewFileStream can be used to read from the local
// filesystem, but you could write an input that reads from a server, network
// etc. Input is a required parameter.
//
// Output is used to provide writers for the atlas files to be written.
// In most cases packer.NewFileOutputter will suffice. Output is a required
// parameter.
//
// Format should be a target format, used to define the descriptor format
// of the atlas. The descriptor acompanies the image to indicate where
// subimages can be found within the atlas. A target format should include
// a valid template and file extension format, all other settings are optional.
//
// Width and Height configure the maximum size of the atlases outputted.
// TODO 0 should be interpreted as no maxumum size.
//
// MaxAtlases can be used to limit the number of atlases outputted. A value
// of 0 is interpreted as no limit.
type Params struct {
	Name          string
	Input         AssetStreamer
	Output        Outputter
	Format        target.Format
	Width, Height int
	Padding       int
	MaxAtlases    int
}

// applySensibleDefaults will fill in nil values with values
// from the list of defaults.
func (p *Params) applySensibleDefaults() {
	if p.Name == "" {
		p.Name = DefaultAtlasName
	}
	if p.Width == 0 {
		p.Width = DefaultAtlasWidth
	}
	if p.Height == 0 {
		p.Height = DefaultAtlasHeight
	}
}

// validateRequiredParameters tests the parameters for
// a non-nil input method and a non-nil output method.
func (p *Params) validateRequiredParameters() error {
	if p.Input == nil {
		return errors.New("Input must not be nil")
	}
	if p.Output == nil {
		return errors.New("Output must not be nil")
	}
	return nil
}

// Run performs the texture packing. It reads files from the given
// AssetStreamer and outputs the results to the given Outputter
// returning an error if any critical failures are encountered.
func Run(ctx context.Context, params *Params) error {
	if ctx == nil {
		return errors.New("Context must not be nil")
	}
	if params == nil {
		return errors.New("Params must not be nil")
	}
	if !params.Format.IsValid() {
		return errors.New("Invalid 'Format' parameter")
	}

	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	// Validate the parameters
	if err := params.validateRequiredParameters(); err != nil {
		return err
	}
	params.applySensibleDefaults()

	// Read the images from the input directory
	sprites, err := readAssetStream(ctx, params.Input, params.Padding)
	if err != nil {
		return err
	}
	// TODO allow sorting algorithm to be specified
	sort.Sort(packing.ByArea(sprites))

	totalNumberOfSprites := len(sprites)
	totalNumberOfAtlases := 0
	completedSprites := make([]packing.Block, 0, totalNumberOfSprites)
	incompleteSprites := make([]packing.Block, 0, totalNumberOfSprites)
	wg := &sync.WaitGroup{}
	errc := make(chan error)
	for {
		// Return error if maxAtlases param exceeded
		if params.MaxAtlases > 0 && totalNumberOfAtlases == params.MaxAtlases {
			return fmt.Errorf("Maximum number of atlases (%d) exceeded", params.MaxAtlases)
		}

		// Arrange the images into the atlas space
		completedSprites = completedSprites[:0]
		incompleteSprites = incompleteSprites[:0]
		packer := packing.NewBinPacker(params.Width, params.Height)
		for _, sprite := range sprites {
			switch packer.Pack(sprite) {
			case packing.ErrInputTooLarge:
				return packing.ErrInputTooLarge
			case packing.ErrOutOfRoom:
				incompleteSprites = append(incompleteSprites, sprite)
			default:
				completedSprites = append(completedSprites, sprite)
			}
		}

		totalNumberOfAtlases++
		atlasName := fmt.Sprintf("%s-%d", params.Name, totalNumberOfAtlases)
		atlas := &Atlas{
			Name:         atlasName,
			Sprites:      completedSprites,
			DescFilename: fmt.Sprintf("%s.%s", atlasName, params.Format.Ext),
			// TODO add image type parameter
			ImageFilename: fmt.Sprintf("%s.%s", atlasName, "png"),
			Width:         params.Width,
			Height:        params.Height,
		}
		wg.Add(1)

		go func(ctx context.Context, errc chan<- error, wg *sync.WaitGroup) {
			select {
			case errc <- atlas.Output(params.Output, params.Format.Template):
			case <-ctx.Done():
			}
			wg.Done()
		}(ctx, errc, wg)

		totalNumberOfIncompletedSprites := len(incompleteSprites)
		// If there are no more sprites that are incomplete, we are done!
		if totalNumberOfIncompletedSprites == 0 {
			break
		}
		// If we don't make any progress, then we've failed
		if totalNumberOfIncompletedSprites == totalNumberOfSprites {
			return packing.ErrOutOfRoom
		}
		// Otherwise continue
		sprites = incompleteSprites
	}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
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

func readAssetStream(ctx context.Context, assetStream AssetStreamer, padding int) ([]packing.Block, error) {
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
			decode(ctx, padding, assets, out)
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
func decode(ctx context.Context, padding int, in <-chan Asset, out chan<- *assetDecodeResult) {
	publishResult := func(spr *sprite, err error) {
		select {
		case out <- &assetDecodeResult{spr, err}:
		case <-ctx.Done():
		}
	}

	for asset := range in {
		defer asset.Close()

		img, _, err := image.Decode(asset)
		if err != nil {
			publishResult(nil, err)
			continue
		}

		rect := img.Bounds()
		spr := &sprite{
			path:    asset.Asset(),
			img:     img,
			w:       rect.Dx(),
			h:       rect.Dy(),
			padding: padding,
		}

		publishResult(spr, nil)
	}
}
