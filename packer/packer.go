package packer

import (
	"fmt"
	"html/template"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"

	"github.com/RaniSputnik/packing"
)

const templatesDir = "templates"

type Params struct {
	Name          string
	Input         string
	Output        string
	Format        string
	Width, Height int
}

func Run(params *Params) error {
	// Validate the parameters
	if !FormatIsValid(params.Format) {
		return fmt.Errorf("Format '%s' is not valid", params.Format)
	}
	descFormat := FormatLookup[params.Format]

	// Create the template that we will use to render descriptor files
	template, err := template.ParseFiles(path.Join(templatesDir, descFormat.Template))
	if err != nil {
		return err
	}

	// Read the images from the input directory
	sprites, err := readDirectory(params.Input)
	if err != nil {
		return err
	}

	// Arrange the images into the atlas space
	packer := &packing.BinPacker{}
	sort.Sort(packing.ByArea(sprites))
	err = packer.Fit(params.Width, params.Height, sprites...)
	if err != nil {
		return err
	}

	// TODO dynamically adjust number of atlases
	atlases := make([]*Atlas, 1)
	atlasName := fmt.Sprintf("%s-%d", params.Name, 1)
	atlases[0] = &Atlas{
		Name:     atlasName,
		Sprites:  sprites,
		DescPath: fmt.Sprintf("%s.%s", atlasName, descFormat.Ext),
		// TODO add image type parameter
		ImagePath: fmt.Sprintf("%s.%s", atlasName, "png"),
		Width:     params.Width,
		Height:    params.Height,
	}

	// TODO should be able to execute all atlases concurrently
	// TODO should write descriptor and image concurrently
	for _, atlas := range atlases {
		// Create and write the resulting image
		err = createImage(atlas, params.Output)
		if err != nil {
			return err
		}
		// Create and write the file that describes the image
		err = createDescriptor(template, atlas, params.Output)
		if err != nil {
			return err
		}
	}

	return nil
}

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

func createImage(atlas *Atlas, outputDir string) error {
	img := atlas.CreateImage()

	writer, err := os.Create(path.Join(outputDir, atlas.ImagePath))
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

func createDescriptor(t *template.Template, atlas *Atlas, outputDir string) error {
	writer, err := os.Create(path.Join(outputDir, atlas.DescPath))
	if err != nil {
		logVerbose("Failed to create desc writer '%s'", err.Error())
		return err
	}
	defer writer.Close()
	t.Execute(writer, atlas)
	return nil
}

func logVerbose(format string, v ...interface{}) {
	//if *pVerbose {
	log.Printf(format, v...)
	//}
}
