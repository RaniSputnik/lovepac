package packer_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/RaniSputnik/lovepac/packer"
)

func TestRunOutputsAtlasAndDescriptor(t *testing.T) {
	files := []string{
		"button_active.png",
		"button_hover.png",
		"button.png",
		"character_evil.png",
		"character_hero.png",
	}
	expected := map[string]string{
		"myatlas-1.png": "",
		"myatlas-1.lua": "",
	}

	outputRecorder := packer.NewOutputRecorder()
	params := &packer.Params{
		Name:   "myatlas",
		Format: packer.FormatLove,
		Input:  packer.NewFilenameStream("./fixtures", files...),
		Output: outputRecorder,
		Width:  1024,
		Height: 1024,
	}

	err := packer.Run(context.Background(), params)
	got := outputRecorder.Got()

	if err != nil {
		t.Errorf("Expected run to succeed without error but got '%s'", err)
	}

	for gotFile := range got {
		if _, ok := expected[gotFile]; !ok {
			t.Errorf("Got unexpected file '%s'", gotFile)
		}
	}

	for expect := range expected {
		if _, ok := got[expect]; !ok {
			t.Errorf("Expected file '%s' to be outputted", expect)
		}
	}
}

func TestRunWithoutParamsSpecifiedUsesSensibleDefaults(t *testing.T) {
	files := []string{"button.png"}
	expected := map[string]string{
		fmt.Sprintf("%s-1.png", packer.DefaultAtlasName): "",
		fmt.Sprintf("%s-1.lua", packer.DefaultAtlasName): "",
	}

	outputRecorder := packer.NewOutputRecorder()
	params := &packer.Params{
		Input:  packer.NewFilenameStream("./fixtures", files...),
		Output: outputRecorder,
	}

	err := packer.Run(context.Background(), params)
	got := outputRecorder.Got()

	if err != nil {
		t.Errorf("Expected run to succeed without error but got '%s'", err)
	}

	for gotFile := range got {
		if _, ok := expected[gotFile]; !ok {
			t.Errorf("Got unexpected file '%s'", gotFile)
		}
	}

	for expect := range expected {
		if _, ok := got[expect]; !ok {
			t.Errorf("Expected file '%s' to be outputted", expect)
		}
	}
}

func TestRunWithNilParamsResultsInError(t *testing.T) {
	emptyParams := &packer.Params{}
	var err error

	err = packer.Run(nil, nil)
	if err == nil {
		t.Errorf("Expected run with nil context and params to fail with error but did not get an error")
	}

	err = packer.Run(nil, emptyParams)
	if err == nil {
		t.Errorf("Expected run with nil context to fail with error but but did not get an error")
	}

	err = packer.Run(context.Background(), emptyParams)
	if err == nil {
		t.Errorf("Expected run with nil input and output to fail with error but but did not get an error")
	}

	err = packer.Run(context.Background(), &packer.Params{
		Input: packer.NewFilenameStream("./fixtures", "button.png"),
	})
	if err == nil {
		t.Errorf("Expected run with nil output to fail with error but but did not get an error")
	}

	err = packer.Run(context.Background(), &packer.Params{
		Output: packer.NewFileOutputter("./doesntmatter"),
	})
	if err == nil {
		t.Errorf("Expected run with nil input to fail with error but but did not get an error")
	}
}

func TestRunWithTooManyFilesForOneAtlasResultsInMultipleAtlases(t *testing.T) {
	files := []string{
		"button_active.png",
		"button_hover.png",
		"button.png",
		"character_evil.png",
		"character_hero.png",
	}
	expected := map[string]string{
		fmt.Sprintf("%s-1.png", packer.DefaultAtlasName): "",
		fmt.Sprintf("%s-1.lua", packer.DefaultAtlasName): "",
		fmt.Sprintf("%s-2.png", packer.DefaultAtlasName): "",
		fmt.Sprintf("%s-2.lua", packer.DefaultAtlasName): "",
	}

	outputRecorder := packer.NewOutputRecorder()
	params := &packer.Params{
		Input:  packer.NewFilenameStream("./fixtures", files...),
		Output: outputRecorder,
		// Here's the crutial part - constrain the width
		// to a size too small for all of the files to fit
		Width:  400,
		Height: 400,
	}

	err := packer.Run(context.Background(), params)
	got := outputRecorder.Got()

	if err != nil {
		t.Errorf("Expected run to succeed without error but got '%s'", err)
	}

	for gotFile := range got {
		if _, ok := expected[gotFile]; !ok {
			t.Errorf("Got unexpected file '%s'", gotFile)
		}
	}

	for expect := range expected {
		if _, ok := got[expect]; !ok {
			t.Errorf("Expected file '%s' to be outputted", expect)
		}
	}
}

func TestRunWithTooManyFilesAndMaxAtlasesResultsInError(t *testing.T) {
	files := []string{
		"button_active.png",
		"button_hover.png",
		"button.png",
		"character_evil.png",
		"character_hero.png",
	}

	outputRecorder := packer.NewOutputRecorder()
	params := &packer.Params{
		Input:  packer.NewFilenameStream("./fixtures", files...),
		Output: outputRecorder,
		// Here's the crutial part - constrain the width
		// to a size too small for all of the files to fit
		// AND limit the number of atlases
		Width:      400,
		Height:     400,
		MaxAtlases: 1,
	}

	err := packer.Run(context.Background(), params)

	if err == nil {
		t.Errorf("Expected run to fail but error was nil")
	}
}
