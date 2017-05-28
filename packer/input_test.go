package packer_test

import (
	"context"
	"testing"

	"github.com/RaniSputnik/lovepac/packer"
)

var fixtures = map[string]struct{}{
	"button_active.png":  {},
	"button_hover.png":   {},
	"button.png":         {},
	"character_evil.png": {},
	"character_hero.png": {},
}

func TestFileStreamSendsAllFiles(t *testing.T) {
	assetStreamer := packer.NewFileStream("./fixtures")
	assets, errc := assetStreamer.AssetStream(context.Background())
	results := map[string]int{}
	resultsChan := make(chan string, len(fixtures))

	go func(results chan<- string) {
		defer close(results)
		for asset := range assets {
			assetName := asset.Asset()
			if _, ok := fixtures[assetName]; !ok {
				t.Errorf("Found unexpected asset named '%s'", assetName)
			} else {
				// No select needed because we know this
				// channel is buffered to receive the correct
				// number of results
				results <- assetName
			}
		}
	}(resultsChan)

	if err := <-errc; err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	for result := range resultsChan {
		results[result]++
	}

	for fixture := range fixtures {
		n, ok := results[fixture]
		if !ok {
			t.Errorf("Expected '%s' but was never received", fixture)
		}
		if n > 1 {
			t.Errorf("Expected to recieve '%s' once but was received '%d' times", fixture, n)
		}
	}
}

func TestFileStreamReportsErrorWhenDirectoryDoesNotExist(t *testing.T) {
	assetStreamer := packer.NewFileStream("./doesnotexist")
	assets, errc := assetStreamer.AssetStream(context.Background())

	go func() {
		for asset := range assets {
			assetName := asset.Asset()
			// There should be no assets streamed
			t.Errorf("Found unexpected asset named '%s'", assetName)
		}
	}()

	if err := <-errc; err == nil {
		t.Errorf("Expected 'directory does not exist' error but got nil")
	}
}

func TestFileStreamIsCancellable(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	assetStreamer := packer.NewFileStream("./fixtures")
	assets, errc := assetStreamer.AssetStream(ctx)

	cancelFunc()

	go func() {
		for asset := range assets {
			assetName := asset.Asset()
			// There should be no assets streamed
			t.Errorf("Found unexpected asset named '%s'", assetName)
		}
	}()

	expectedErr := ctx.Err()
	if gotErr := <-errc; gotErr != expectedErr {
		t.Errorf("Expected '%s' but got '%s'", expectedErr, gotErr)
	}
}
