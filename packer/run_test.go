package packer_test

import (
	"context"
	"fmt"
	"testing"

	"strings"

	"github.com/RaniSputnik/lovepac/packer"
	"github.com/RaniSputnik/lovepac/target"
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

	outputRecorder := NewOutputRecorder()
	params := &packer.Params{
		Name:   "myatlas",
		Format: target.Love,
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

	outputRecorder := NewOutputRecorder()
	params := &packer.Params{
		Format: target.Love,
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

	outputRecorder := NewOutputRecorder()
	params := &packer.Params{
		Format: target.Love,
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

	outputRecorder := NewOutputRecorder()
	params := &packer.Params{
		Format: target.Love,
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

func TestPaddingIsAppliedCorrectly(t *testing.T) {
	button := "button.png"
	buttonWidth, buttonHeight := 124, 50
	padding := 2

	outputRecorder := NewOutputRecorder()
	params := &packer.Params{
		Input:   packer.NewFilenameStream("./fixtures", button),
		Output:  outputRecorder,
		Name:    "atlas",
		Format:  target.Love,
		Padding: padding,
	}

	err := packer.Run(context.Background(), params)
	got := outputRecorder.Got()

	if err != nil {
		t.Errorf("Expected run to succeed without error but got '%s'", err)
	}

	expectedString := fmt.Sprintf("quads['button'] = love.graphics.newQuad(%d,%d,%d,%d,%d,%d)",
		padding, padding, buttonWidth, buttonHeight, packer.DefaultAtlasWidth, packer.DefaultAtlasHeight)
	seperator := createUnderlineString(expectedString)
	gotStr := got["atlas-1.lua"].String()
	if !strings.Contains(gotStr, expectedString) {
		t.Errorf("Expected descriptor to contain the following sub-string\n\n%s\n%s\n\n%s", expectedString, seperator, gotStr)
	}

	// TODO do we want to ensure the image was placed correctly too?
}

func TestAssetsDoNotFitIfPaddingCannotBeApplied(t *testing.T) {
	button := "button.png"
	buttonWidth, buttonHeight := 124, 50

	outputRecorder := NewOutputRecorder()
	params := &packer.Params{
		Format:  target.Love,
		Input:   packer.NewFilenameStream("./fixtures", button),
		Output:  outputRecorder,
		Padding: 2,
		Width:   buttonWidth,
		Height:  buttonHeight,
	}

	err := packer.Run(context.Background(), params)
	if err == nil {
		t.Errorf("Expected run to fail but unstead got nil error")
	}
}

func createUnderlineString(input string) string {
	inputLength := len(input)
	chars := make([]rune, inputLength)
	for i := 0; i < inputLength; i++ {
		chars[i] = '~'
	}
	return string(chars)
}
