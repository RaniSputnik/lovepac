# lovepac
A texture packer for [love2d](https://love2d.org) written in Go lang.

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
go build
./lovepac -format love -out build ./assets/
```

### Adding Output Targets

Targets are generated from templates using the `/target/gen.go` function. This is run by go generate.

To add a new target output format;

1. Add a new `*.template` file in the `/target` directory. Must be in the go template format.
2. Run `go generate ./target` to regenerate the templates in the `/target/target_generated.go` file.
3. Finally export the target by adding a `template.Fomat` to the `target/target.go` file.

Use should now be able to reference your target by name from the `target` package.