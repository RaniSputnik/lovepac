/*
Package packer provides a texture packing implementation that will
take a list on input images and build them into a larger image and
a descriptor file.

Texture packing is usually used to optimise runtime performance
where a number of textures must be drawn to the screen at high
frequency but can also be used to improve load times for pages
loaded with HTTP 1.x.

A simple example of usage is;

	params := packer.Params{
		Name:   "myatlas",
		Format: target.Love,
		Input:  packer.NewFileStream("./assets"),
		Output: packer.NewFileOutputter("./build"),
	}
	log.Fatal(packer.Run(context.Background(), &params))

You can specify maximum width and height of an atlas to conform to
platform limitations and you can build multiple atlases with a single
command.

The Input and Output parameters have been designed to be highly flexible.
Simple file system readers and writers are provided out of the box but
consumers could implement the AssetStreamer and Outputter interfaces to
support more complex texture packing environments eg. Reading and writing
to Google Cloud Storage.

For example, here is a complete AssetStreamer that reads assets from memory.

	// Define the Asset that we will return
	// The Asset interface simply augments the
	// io.ReadCloser interface with another simple
	// method Asset, that returns the asset name
	type BytesAsset struct {
		Name string
		*bytes.Reader
	}

	func (b *BytesAsset) Asset() string { return b.Name }
	func (b *BytesAsset) Close() error  { return nil } // noop, implement ReadCloser

	// Create the Asset and reader
	func NewBytesAsset(name string, data []byte) *BytesAsset {
		return &BytesAsset{
			Name:   name,
			Reader: bytes.NewReader(data),
		}
	}

	// Here is the bulk of the implementation, we input a number of assets
	// and output them on an asset stream for decoding. There is no chance
	// of error here so we leave errc unbuffered.
	func NewBytesAssetStream(images ...*BytesAsset) packer.AssetStreamer {
		return packer.AssetStreamerFunc(func(ctx context.Context) (<-chan packer.Asset, <-chan error) {
			stream := make(chan packer.Asset)
			errc := make(chan error)
			go func() {
				defer close(stream)
				defer close(errc)
				for _, img := range images {
					select {
					case stream <- img:
					case <-ctx.Done():
						break
					}
				}
			}()
			return stream, errc
		})
	}

Though that seems like a lot, for the majority of use cases, a custom
asset streamer implementation will not be necessary.
*/
package packer
