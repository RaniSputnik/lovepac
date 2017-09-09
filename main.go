package main

import (
	"context"
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	"runtime/pprof"
	"time"

	"github.com/RaniSputnik/lovepac/packer"
	"github.com/RaniSputnik/lovepac/target"

	"log"
	"os"
)

// Command line arguments
var pVerbose *bool

func main() {

	// Set the function to call when printing command line usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage : %s -flags <inputdir>\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Set and parse the command line arguments
	pName := flag.String("name", packer.DefaultAtlasName, "the base name of the output images and data files")
	pOutputDir := flag.String("out", "", "the directory to output the result to")
	pVerbose = flag.Bool("v", false, "use verbose logging")
	pFormat := flag.String("format", "love", "the export format of the atlas")
	pWidth := flag.Int("width", packer.DefaultAtlasWidth, "maximum width of an atlas image")
	pHeight := flag.Int("height", packer.DefaultAtlasHeight, "maximum height of an atlas image")
	pPadding := flag.Int("padding", 0, "the space between images in the atlas")
	pMaxAtlases := flag.Int("maxatlases", 0, "the maximum number of atlases to write, 0 indicates no maximum")
	pCPUProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	pMemprofile := flag.String("memprofile", "", "write memory profile to file")

	flag.Parse()

	if *pCPUProfile != "" {
		f, err := os.Create(*pCPUProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// Get the input directory
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Too few arguments passed, missing input directory \n\n")
		flag.Usage()
		return
	}
	inputDir := args[0]

	format := target.FormatNamed(*pFormat)
	if format == target.Unknown {
		log.Fatalf("Unknown format '%s'", *pFormat)
	}

	stopTimer := startTimer("Texture packing")
	err := packer.Run(context.Background(), &packer.Params{
		Name:       *pName,
		Input:      packer.NewFileStream(inputDir),
		Output:     packer.NewFileOutputter(*pOutputDir),
		Format:     format,
		Width:      *pWidth,
		Height:     *pHeight,
		Padding:    *pPadding,
		MaxAtlases: *pMaxAtlases,
	})
	stopTimer()

	if *pMemprofile != "" {
		f, err := os.Create(*pMemprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	if err != nil {
		log.Fatal(err)
	}
}

func startTimer(name string) func() {
	start := time.Now()
	return func() {
		log.Printf("%s took %s", name, time.Since(start))
	}
}
