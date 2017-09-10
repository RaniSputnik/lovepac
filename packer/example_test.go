package packer_test

import (
	"context"
	"log"

	"github.com/RaniSputnik/lovepac/packer"
	"github.com/RaniSputnik/lovepac/target"
)

func ExampleBasic() {
	params := packer.Params{
		Name:   "myatlas",
		Format: target.Love,
		Input:  packer.NewFileStream("./assets"),
		Output: packer.NewFileOutputter("./build"),
		Width:  512,
		Height: 512,
	}
	log.Fatal(packer.Run(context.Background(), &params))
}
