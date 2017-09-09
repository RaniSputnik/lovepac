package packer

import (
	"io"
	"os"
	"path"
)

// Outputter is a factory responsible for creating writers that
// atlas files can be written to. In most cases FileOutputter will
// be used to allow you to easilly write to a system directory
// but outputters can be used to write to any destination.
type Outputter interface {
	GetWriter(filename string) (io.WriteCloser, error)
}

// OutputterFunc is a function that conforms to the Outputter interface
type OutputterFunc func(filename string) (io.WriteCloser, error)

func (f OutputterFunc) GetWriter(filename string) (io.WriteCloser, error) {
	return f(filename)
}

// NewFileOutputter is most common form of atlas outputter. Specify an empty
// output directory and it will write all atlas contents to this new directory
// using the os standard library.
func NewFileOutputter(outputDirectory string) Outputter {
	return OutputterFunc(func(filename string) (io.WriteCloser, error) {
		return os.Create(path.Join(outputDirectory, filename))
	})
}

// Helper method that takes care of opening / closing a file with the given outputter
func withFile(outputter Outputter, filename string, do func(writer io.Writer) error) error {
	writer, err := outputter.GetWriter(filename)
	if err != nil {
		return err
	}
	defer writer.Close()
	return do(writer)
}
