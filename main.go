package main

import (
	"flag"
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	"runtime/pprof"
	"time"

	"github.com/RaniSputnik/lovepac/packer"

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
	pName := flag.String("name", "atlas", "the base name of the output images and data files")
	pOutputDir := flag.String("out", "", "the directory to output the result to")
	pVerbose = flag.Bool("v", false, "use verbose logging")
	pFormat := flag.String("format", "starling", "the export format of the atlas")
	pWidth := flag.Int("width", 2048, "maximum width of an atlas image")
	pHeight := flag.Int("height", 2048, "maximum height of an atlas image")
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

	stopTimer := startTimer("Texture packing")
	err := packer.Run(&packer.Params{
		Name:   *pName,
		Input:  packer.NewFileStream(inputDir),
		Output: packer.NewFileOutputter(*pOutputDir),
		Format: *pFormat,
		Width:  *pWidth,
		Height: *pHeight,
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
