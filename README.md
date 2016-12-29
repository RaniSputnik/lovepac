# lovepac
A texture packer for [love2d](https://love2d.org) written in Go lang.

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
