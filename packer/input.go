package packer

import (
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// Asset represents a single input source into the texture packer.
// Many assets are supplied and packed together to create a single atlas.
//
// Assets commonly represent files in a filesystem, but could also
// be blobs in a blobstore or images on a remote server.
type Asset interface {
	io.ReadCloser

	// Asset returns the name of the given asset
	Asset() string
}

// AssetStreamer is a factory responsible for creating readers that
// can read the contents of a given asset
type AssetStreamer interface {
	AssetStream(ctx context.Context) (<-chan Asset, <-chan error)
}

// AssetStreamerFunc is a function that conforms to the AssetStreamer interface
type AssetStreamerFunc func(ctx context.Context) (<-chan Asset, <-chan error)

// AssetStream implements the AssetStreamer interface
func (f AssetStreamerFunc) AssetStream(ctx context.Context) (<-chan Asset, <-chan error) {
	return f(ctx)
}

type FileAsset struct {
	*os.File
	Name string
}

func (a *FileAsset) Asset() string {
	return a.Name
}

var errContextNil = errors.New("Context must not be nil")

// NewFileStream creates an asset streamer that streams files from a given
// input directory. The input directory will be walked and readers will be
// created using the standard os package.
func NewFileStream(inputDirectory string) AssetStreamer {
	return AssetStreamerFunc(func(ctx context.Context) (<-chan Asset, <-chan error) {
		stream := make(chan Asset)
		errc := make(chan error, 1)
		go func() {
			defer close(stream)
			defer close(errc)

			if ctx == nil {
				errc <- errContextNil
				return
			}

			// No select needed for this send, since errc is buffered.
			errc <- filepath.Walk(inputDirectory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					// TODO walk nested directories...
					return nil
				}
				if !info.Mode().IsRegular() {
					return nil
				}

				relPath, err := filepath.Rel(inputDirectory, path)
				if err != nil {
					return err
				}

				file, err := os.Open(path)
				if err != nil {
					return err
				}

				select {
				case stream <- &FileAsset{File: file, Name: relPath}:
				case <-ctx.Done():
					return ctx.Err()
				}
				return nil
			})
		}()
		return stream, errc
	})
}

// NewFilenameStream creates an asset streamer that streams the specified files.
// The files will be read relative to the given directory.
// Readers will be created using the standard os package.
func NewFilenameStream(directory string, files ...string) AssetStreamer {
	return AssetStreamerFunc(func(ctx context.Context) (<-chan Asset, <-chan error) {
		stream := make(chan Asset)
		errc := make(chan error, 1)

		go func() {
			defer close(stream)
			defer close(errc)

			if ctx == nil {
				errc <- errContextNil
				return
			}

			for _, filename := range files {
				path := filepath.Join(directory, filename)
				reader, err := os.Open(path)
				if err != nil {
					errc <- err
				}
				select {
				case stream <- &FileAsset{File: reader, Name: filename}:
				case <-ctx.Done():
					errc <- ctx.Err()
					return
				}
			}
		}()

		return stream, errc
	})
}

// TODO we could match globs too, that'd be cool
