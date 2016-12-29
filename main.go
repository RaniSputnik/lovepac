package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"

	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"text/template"

	"github.com/RaniSputnik/packing"
)

// Command line arguments
var (
	pName      *string
	pOutputDir *string
	pVerbose   *bool
	pFormat    *string
	pWidth     *int
	pHeight    *int
)

func main() {

	// Set the function to call when printing command line usage
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage : %s -flags <inputdir>\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Set and parse the command line arguments
	pName = flag.String("name", "atlas", "the base name of the output images and data files")
	pOutputDir = flag.String("out", "", "the directory to output the result to")
	pVerbose = flag.Bool("v", false, "use verbose logging")
	pFormat = flag.String("format", "starling", "the export format of the atlas")
	pWidth = flag.Int("width", 2048, "maximum width of an atlas image")
	pHeight = flag.Int("height", 2048, "maximum height of an atlas image")
	flag.Parse()

	// Validate the parameters
	if !FormatIsValid(*pFormat) {
		log.Fatalf("Format '%s' is not valid", *pFormat)
	}

	// Get the input directory
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Too few arguments passed, missing input directory \n\n")
		flag.Usage()
		return
	}
	inputDir := args[0]

	// Read the images from the input directory
	sprites, err := readDirectory(inputDir)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Arrange the images into the atlas space
	packer := &packing.BinPacker{}
	sort.Sort(packing.ByArea(sprites))
	err = packer.Fit(*pWidth, *pHeight, sprites...)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create and write the resulting image
	err = createImage(sprites)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = createDescriptor(FormatLookup[*pFormat], sprites)
	if err != nil {
		log.Fatal(err.Error())
	}
}

// ReadSprites returns all spritees found in the provided directory
func readDirectory(dir string) ([]packing.Block, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// TODO filter out files that are not .png/.jpg

	sprites := make([]packing.Block, len(files))

	for i := range sprites {
		path := path.Join(dir, files[i].Name())

		logVerbose("Reading input file '%s'", path)

		reader, err := os.Open(path)
		if err != nil {
			logVerbose("Failed to create reader '%s'", err.Error())
			return sprites, err
		}
		defer reader.Close()

		img, _, err := image.Decode(reader)
		if err != nil {
			logVerbose("Failed to decode file '%s'", err.Error())
			return sprites, err
		}

		rect := img.Bounds()
		sprites[i] = &sprite{
			path: path,
			img:  img,
			w:    rect.Dx(),
			h:    rect.Dy(),
		}
	}

	return sprites, nil
}

func createImage(sprites []packing.Block) error {
	img := image.NewRGBA(image.Rect(0, 0, *pWidth, *pHeight))

	for i := range sprites {
		spr := sprites[i].(*sprite)
		rect := image.Rect(spr.x, spr.y, spr.x+spr.w, spr.y+spr.h)
		draw.Draw(img, rect, spr.img, image.ZP, draw.Src)
	}

	index := 1
	format := "png"
	filename := fmt.Sprintf("%s-%d.%s", *pName, index, format)
	writer, err := os.Create(path.Join(*pOutputDir, filename))
	if err != nil {
		logVerbose("Failed to create image writer '%s'", err.Error())
		return err
	}
	defer writer.Close()

	err = png.Encode(writer, img)
	if err != nil {
		return err
	}
	return nil
}

type descriptor struct {
	DescPath  string
	ImagePath string
	Files     []packing.Block
}

const templatesDir = "templates"

func createDescriptor(format Format, sprites []packing.Block) error {
	t, err := template.ParseFiles(path.Join(templatesDir, format.Template))
	if err != nil {
		return err
	}

	index := 1
	filename := fmt.Sprintf("%s-%d.%s", *pName, index, format.Ext)
	writer, err := os.Create(path.Join(*pOutputDir, filename))
	if err != nil {
		logVerbose("Failed to create desc writer '%s'", err.Error())
		return err
	}
	defer writer.Close()

	imageformat := "png"
	t.Execute(writer, &descriptor{
		DescPath:  filename,
		ImagePath: fmt.Sprintf("%s-%d.%s", *pName, index, imageformat),
		Files:     sprites,
	})
	return nil
}

func logVerbose(format string, v ...interface{}) {
	if *pVerbose {
		log.Printf(format, v...)
	}
}
