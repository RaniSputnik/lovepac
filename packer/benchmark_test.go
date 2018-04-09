package packer_test

import (
	"context"
	"fmt"
	"testing"

	_ "image/gif"
	_ "image/jpeg"

	"github.com/RaniSputnik/lovepac/packer"
	"github.com/RaniSputnik/lovepac/target"
)

func BenchmarkPack512x512(b *testing.B) {
	benchmarkPackX(b, 512, 512)
}

func BenchmarkPack1024x1024(b *testing.B) {
	benchmarkPackX(b, 1024, 1024)
}

func BenchmarkPack2048x2048(b *testing.B) {
	benchmarkPackX(b, 2048, 2048)
}

func BenchmarkPack4096x4096(b *testing.B) {
	benchmarkPackX(b, 2048, 2048)
}

func benchmarkPackX(b *testing.B, w, h int) {
	params := &packer.Params{
		Name:   "myatlas",
		Format: target.Love,
		Input:  packer.NewFileStream("../assets"),
		Output: packer.NewFileOutputter(fmt.Sprintf("../build")),
		Width:  w,
		Height: h,
	}

	for n := 0; n < b.N; n++ {
		if err := packer.Run(context.Background(), params); err != nil {
			b.Fatalf("%s", err)
		}
	}
}
