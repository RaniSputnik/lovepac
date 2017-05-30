package packer

import (
	"bytes"
	"io"
	"os"
	"path"
	"sync"
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

type OutputRecorder struct {
	writers map[string]*bufferWithClose
	*sync.Mutex
}

type bufferWithClose struct {
	*bytes.Buffer
}

func (b *bufferWithClose) Close() error { return nil }

func (r *OutputRecorder) GetWriter(filename string) (io.WriteCloser, error) {
	buffer := &bufferWithClose{bytes.NewBufferString("")}
	r.Lock()
	r.writers[filename] = buffer
	r.Unlock()
	return buffer, nil
}

func (r *OutputRecorder) Got() map[string]string {
	r.Lock()
	results := map[string]string{}
	for key, val := range r.writers {
		results[key] = val.String()
	}
	r.Unlock()
	return results
}

func NewOutputRecorder() *OutputRecorder {
	return &OutputRecorder{map[string]*bufferWithClose{}, &sync.Mutex{}}
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
