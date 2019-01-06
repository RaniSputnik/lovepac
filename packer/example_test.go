package packer_test

import (
	"context"
	"fmt"

	"github.com/RaniSputnik/lovepac/packer"
	"github.com/RaniSputnik/lovepac/target"
)

func ExampleRun() {
	params := packer.Params{
		Name:   "myatlas",
		Format: target.Love,
		Input:  packer.NewFileStream("./assets"),
		Output: packer.NewFileOutputter("./build"),
		Width:  512,
		Height: 512,
	}
	if err := packer.Run(context.Background(), &params); err != nil {
		fmt.Print("Texture packing complete")
	}
	// Output: Texture packing complete
}
