# lovepac
A texture packer for [love2d](https://love2d.org) written in Go lang.

### Installation

With [go](https://golang.org/) installed;

```
go get github.com/RaniSputnik/lovepac
```

### Features

- Optimise your texture memory usage
- Specify maximum width and height to conform to platform limitations
- FAST
- Generate as many atlases as you need with a single command
- Flexible input and output interfaces to read and write atlases to disk/network/wherever
- No-fuss installation, 100% go code

### Usage

```
Usage : lovepac -flags <inputdir>
  -format string
    	the export format of the atlas (default "starling")
  -height int
    	maximum height of an atlas image (default 2048)
  -name string
    	the base name of the output images and data files (default "atlas")
  -out string
    	the directory to output the result to
  -v	use verbose logging
  -width int
    	maximum width of an atlas image (default 2048)
```

Eg. Pack all files in ./assets directory and output to ./build in love format;

```
lovepac -format love -out build ./assets/
```

### Package

This texture packer can also be used as a library by consuming the packer and target
packages.

```
import (
  ...
  "github.com/RaniSputnik/lovepac/packer"
  "github.com/RaniSputnik/lovepac/target"
)
...
params := packer.Params{
  Name:   "myatlas",
  Format: target.Love,
  Input:  packer.NewFileStream("./assets"),
  Output: packer.NewFileOutputter("./build"),
  Width:  512,
  Height: 512,
}
log.Fatal(packer.Run(context.Background(), &params))
```

See the [godoc](https://godoc.org/github.com/RaniSputnik/lovepac/packer) for
more information and examples.

### Adding Output Targets

Targets are generated from templates using the `/target/gen.go` function. This is run by go generate.

To add a new target output format;

1. Add a new `*.template` file in the `/target` directory. Must be in the go template format.
2. Run `go generate ./target` to regenerate the templates in the `/target/target_generated.go` file.
3. Finally export the target by adding a `template.Fomat` to the `target/target.go` file.

Use should now be able to reference your target by name from the `target` package.

### Benchmarks

These are the results I get on my machine when packing 55 sample assets.

```
$ go test ./packer -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: github.com/RaniSputnik/lovepac/packer
BenchmarkPack512x512-8                        50          27129564 ns/op        13675566 B/op       8004 allocs/op
BenchmarkPack1024x1024-8                      10         176355710 ns/op        17693480 B/op     160332 allocs/op
BenchmarkPack2048x2048-8                       5         287248694 ns/op        25469216 B/op     160370 allocs/op
BenchmarkPack4096x4096-8                       5         290687305 ns/op        25469579 B/op     160372 allocs/op
BenchmarkPackOneAsset512x512-8               100          15301139 ns/op         2005043 B/op        186 allocs/op
BenchmarkPackOneAsset1024x1024-8              30          40203031 ns/op         5166106 B/op        185 allocs/op
BenchmarkPackOneAsset2048x2048-8              10         135694018 ns/op        17776404 B/op        185 allocs/op
BenchmarkPackOneAsset4096x4096-8              10         140962040 ns/op        17776521 B/op        186 allocs/op
```
